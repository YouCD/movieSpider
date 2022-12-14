package douban

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/config"
	httpClient2 "movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const movieUrlPrefix = "https://movie.douban.com/subject/"

type DouBan struct {
	doubanUrl  string
	scheduling string
}
type rowData struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Image    string `json:"image"`
	Director []struct {
		Type string `json:"type"`
		Url  string `json:"url"`
		Name string `json:"name"`
	} `json:"director"`
	Author []struct {
		Type string `json:"type"`
		Url  string `json:"url"`
		Name string `json:"name"`
	} `json:"author"`
	Actor []struct {
		Type string `json:"type"`
		Url  string `json:"url"`
		Name string `json:"name"`
	} `json:"actor"`
	DatePublished   string   `json:"datePublished"`
	Genre           []string `json:"genre"`
	Duration        string   `json:"duration"`
	Description     string   `json:"description"`
	Type            string   `json:"type"`
	AggregateRating struct {
		Type        string `json:"type"`
		RatingCount string `json:"ratingCount"`
		BestRating  string `json:"bestRating"`
		WorstRating string `json:"worstRating"`
		RatingValue string `json:"ratingValue"`
	} `json:"aggregateRating"`
}

func NewSpiderDouBan(doubanUrl, scheduling string) *DouBan {
	return &DouBan{
		doubanUrl,
		scheduling,
	}
}

func (d *DouBan) Crawler() (videos []*types.DouBanVideo) {

	doc, err := d.newRequest(d.doubanUrl)
	if err != nil {
		log.Error(err)
		return
	}
	var videos1 []*types.DouBanVideo
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

			videos1 = append(videos1, doubanVideo)
		})

	})

	var videos2 []*types.DouBanVideo
	var wg sync.WaitGroup
	for _, video := range videos1 {
		wg.Add(1)
		// 访问 豆瓣 具体的电影首页
		doc, err := d.newRequest(fmt.Sprintf("%s%s", movieUrlPrefix, video.DoubanID))
		if err != nil {
			wg.Done()
			log.Error(err)
			return
		}
		// 获取电影原始数据
		content := doc.Find("script[type='application/ld+json']").Text()
		content = strings.ReplaceAll(content, "\n", "")
		content = strings.ReplaceAll(content, "@type", "type")
		var data rowData
		err = json.Unmarshal([]byte(content), &data)
		if err != nil {
			wg.Done()
			log.Error(err)
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
		video.ImdbID = matchArr[0]

		videos2 = append(videos2, video)
		wg.Done()
	}
	wg.Wait()

	for _, video := range videos2 {
		err = model.NewMovieDB().CreatDouBanVideo(video)
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("DouBan %s 已保存", video.Names)
	}

	return
}

func (d *DouBan) newRequest(url string) (document *goquery.Document, err error) {

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	if config.DouBan.Cookie != "" {
		request.Header.Set("Cookie", config.DouBan.Cookie)
	}
	client := httpClient2.NewHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("未能正常获取豆瓣数据")
	}
	defer resp.Body.Close()

	document, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return
}
func (d *DouBan) zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
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
	c.AddFunc(d.scheduling, func() { d.Crawler() })
	c.Start()
}
