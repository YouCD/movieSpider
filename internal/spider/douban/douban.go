package douban

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/config"
	httpClient2 "movieSpider/internal/httpclient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/spider"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const movieURLPrefix = "https://movie.douban.com/subject/"

type DouBan struct {
	url        string
	scheduling string
}

// NewSpiderDouBan
//
//	@Description: 新建DouBan
//	@param doubanUrl
//	@param scheduling
//	@return *DouBan
func NewSpiderDouBan(cfg *config.DouBan) (douBanList []spider.Spider) {
	for _, db := range cfg.DouBanList {
		if db.Scheduling == "" {
			db.Scheduling = cfg.Scheduling
		}
		douBanList = append(douBanList, &DouBan{
			url:        db.URL,
			scheduling: db.Scheduling,
		})
	}

	return
}

// Crawler
//
//	@Description: 爬取
//	@receiver d
//	@return videos
//
//nolint:gosimple
func (d *DouBan) Crawler() (videos []*types.DouBanVideo) {
	doc, err := d.newRequest(d.url)
	if err != nil {
		log.Error(err)
		//nolint:nakedret
		return
	}
	var summaryVideo []*types.DouBanVideo

	doc.Find("#content > div.grid-16-8.clearfix > div.article > div.grid-view> div").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, selection *goquery.Selection) {
			doubanVideo := new(types.DouBanVideo)
			// 片名
			nameStr := selection.Find("div.info> ul > li.title > a > em ").Text()
			doubanVideo.Names = nameStr
			//#content > div.grid-16-8.clearfix > div.article > div.grid-view > div:nth-child(1) > div.info > ul > li.title > a
			val, _ := selection.Find("div.info>ul > li.title > a").Attr("href")

			compileRegex := regexp.MustCompile("[0-9]+")
			matchArr := compileRegex.FindStringSubmatch(val)
			doubanVideo.DoubanID = matchArr[0]
			//#content > div.grid-16-8.clearfix > div.article > div.grid-view > div:nth-child(2) > div.info > ul > li.title > span
			Playable := selection.Find(" div.info > ul > li.title > span").Text()
			Playable = strings.ReplaceAll(Playable, "[", "")
			Playable = strings.ReplaceAll(Playable, "]", "")
			doubanVideo.Playable = Playable

			summaryVideo = append(summaryVideo, doubanVideo)
		})
	})
	//nolint:prealloc
	var videos2 []*types.DouBanVideo
	var wg sync.WaitGroup
	for _, video := range summaryVideo {
		wg.Add(1)
		// 访问 豆瓣 具体的电影首页
		doc, err := d.newRequest(fmt.Sprintf("%s%s", movieURLPrefix, video.DoubanID))
		if err != nil {
			wg.Done()
			log.Error(err)
			//nolint:nakedret
			return
		}
		// 获取电影原始数据
		content := doc.Find("script[type='application/ld+json']").Text()
		content = strings.ReplaceAll(content, "\n", "")
		content = strings.ReplaceAll(content, "@type", "type")
		var data types.RowData
		err = json.Unmarshal([]byte(content), &data)
		if err != nil {
			wg.Done()
			log.Error(err)
			//nolint:nakedret
			return
		}

		// 处理 Genre
		var genre []string
		for _, g := range data.Genre {
			unicode, err := d.zhToUnicode([]byte(g))
			if err != nil {
				return nil
			}
			genre = append(genre, string(unicode))
		}
		data.Genre = genre
		marshal, err := json.Marshal(data)
		if err != nil {
			return nil
		}
		// 赋值 原始数据
		video.RowData = string(marshal)

		// 上映时间
		video.DatePublished = data.DatePublished

		// 处理类型
		video.Type = video.FormatType(data.Type)
		// 处理 名称
		video.Names = video.FormatName(video.Names)

		html, err := doc.Html()
		if err != nil {
			return nil
		}

		compileRegex := regexp.MustCompile("tt\\d+")
		matchArr := compileRegex.FindStringSubmatch(html)
		if len(matchArr) > 0 {
			video.ImdbID = matchArr[0]
		}

		videos2 = append(videos2, video)
		wg.Done()
	}
	wg.Wait()

	for _, video := range videos2 {
		if err = model.NewMovieDB().CreatDouBanVideo(video); err != nil && errors.Is(err, model.ErrDataExist) {
			log.Debugf("DouBan %s 已更新", video.Names)
		} else {
			log.Infof("DouBan %s 已添加", video.Names)
		}
	}

	return videos2
}

// newRequest
//
//	@Description: 新建请求
//	@receiver d
//	@param url
//	@return document
//	@return err
func (d *DouBan) newRequest(url string) (document *goquery.Document, err error) {
	request, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "newRequest")
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	request.Header.Set("Cookie", `bid=TfBEKp4BWMc; dbcl2="251312920:rkR4Ujzy3W8"; ck=5MvF; ll="118267"; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1708572426%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_id.100001.4cf6=aecc8277cdf140cf.1708572426.; _pk_ses.100001.4cf6=1; __utma=30149280.1796252786.1708572426.1708572426.1708572426.1; __utmb=30149280.0.10.1708572426; __utmc=30149280; __utmz=30149280.1708572426.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utma=223695111.1879872192.1708572426.1708572426.1708572426.1; __utmb=223695111.0.10.1708572426; __utmc=223695111; __utmz=223695111.1708572426.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; push_noty_num=0; push_doumail_num=0; frodotk_db="202dcb43b7ae3817608961cb270aa845"; _vwo_uuid_v2=D8A0FB067F659FDB23642E61CBF1F722B|11e34d62d2eb48d088514b92792fc6b5`)
	client := httpClient2.NewHTTPClient()
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.WithMessage(err, "client.Do")
	}
	if resp == nil {
		return nil, errors.New("未能正常获取豆瓣数据")
	}
	defer resp.Body.Close()

	document, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "goquery.NewDocumentFromReader")
	}
	return
}

// zhToUnicode
//
//	@Description: 中文转 unicode
//	@receiver d
//	@param raw
//	@return []byte
//	@return error
//
//nolint:gocritic
func (d *DouBan) zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, errors.WithMessage(err, "zhToUnicode")
	}
	return []byte(str), nil
}
func (d *DouBan) Run() {
	if d.scheduling == "" {
		log.Error("DouBan Scheduling is null")
		os.Exit(1)
	}
	log.Infof("DouBan Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(d.scheduling, func() { d.Crawler() })
	c.Start()
}

// todo 还需要搞一个定时任务，定时更新 DatePublished
