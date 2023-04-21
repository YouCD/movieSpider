package feed

import (
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/magnetConvert"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const urlBtbt = "https://www.btbtt12.com/forum-index-fid-951.htm"

type btbt struct {
	url        string
	scheduling string
}

func (b *btbt) Crawler() (Videos []*types.FeedVideo, err error) {
	c := httpClient.NewHttpClient()
	req, err := http.NewRequest("GET", b.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "getMovies resp")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "getMovies goquery")
	}

	// 洗出  下载电影的页面
	var Videos1 []*types.FeedVideo
	doc.Find("[valign='middle'].subject").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, selection *goquery.Selection) {
			fVideo := new(types.FeedVideo)
			selection.Find("a:nth-child(6)").Each(func(i int, selection *goquery.Selection) {
				name := splitTitle(selection.Text())
				u, ok := selection.Attr("href")
				if ok {
					if name == "" {
						return
					}
					fVideo.TorrentUrl = urlParser(u, b.url)
					fVideo.Name = name
				}
			})

			// 年份
			selection.Find("a:nth-child(2)").Each(func(i int, selection *goquery.Selection) {
				fVideo.Year = trim(selection.Text())

			})

			fVideo.Type = "movie"
			Videos1 = append(Videos1, fVideo)
		})
	})

	// 洗出  下载电影种子下载的页面
	var wg sync.WaitGroup
	var Videos2 []*types.FeedVideo
	for _, v := range Videos1 {
		parse, err := url.Parse(v.TorrentUrl)
		if err != nil {
			log.Error(err)
			continue
		}
		if parse.String() == "" {
			continue
		}
		wg.Add(1)
		downloadUrl, err := moviePageUrl(v.TorrentUrl)
		if err != nil {
			log.Error(err)
		}
		if downloadUrl == "" {
			wg.Done()
			continue
		}
		v.TorrentName = v.Name
		v.RowData = sql.NullString{String: downloadUrl}
		Videos2 = append(Videos2, v)
		wg.Done()
	}
	wg.Wait()

	// 洗出  磁力连接
	var wg1 sync.WaitGroup
	for _, v := range Videos2 {
		wg1.Add(1)
		torrentDownloadUrlStr, err := torrentDownloadUrl(v.RowData.String)
		if err != nil {
			return nil, err
		}
		magnet, err := getMagnet(torrentDownloadUrlStr)
		if err != nil {
			log.Error(err)
		}
		v.Magnet = magnet
		v.Web = "btbt"
		Videos = append(Videos, v)
		wg1.Done()
	}
	wg1.Wait()

	return

}

func (b *btbt) Run() {
	if b.scheduling == "" {
		log.Error("BTBT: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("BTBT: Scheduling is: [%s]", b.scheduling)
	c := cron.New()
	c.AddFunc(b.scheduling, func() {
		videos, err := b.Crawler()
		if err != nil {
			log.Error(err)
			return
		}

		proxySaveVideo2DB(videos...)
	})
	c.Start()

}

func splitTitle(str string) string {
	if strings.Contains(str, "720") {
		return ""
	}
	str = strings.TrimSpace(str)
	if len(strings.Split(str, "][")) > 1 {
		return strings.TrimSpace(strings.Split(str, "][")[1])
	}
	return ""
}
func trim(str string) string {
	str = strings.Trim(str, "[")
	str = strings.Trim(str, "]")
	return str
}
func urlParser(uri, urlStr string) string {
	parse, _ := url.Parse(urlStr)
	if uri != "" {
		parse.Path = uri
		return parse.String()
	} else {
		return parse.String()
	}
}

func moviePageUrl(pageUrl string) (url string, err error) {
	resp, err := client(pageUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.WithMessagef(err, "BTBT: 电影页没有数据: %s", pageUrl)
	}

	doc.Find("#body > div > table:nth-child(2) > tbody > tr:nth-child(1) > td.post_td > div.attachlist > table > tbody > tr:nth-child(3) > td:nth-child(1) > a").Each(func(i int, selection *goquery.Selection) {
		u, ok := selection.Attr("href")
		if ok {
			url = urlParser(u, pageUrl)
		}
	})
	return
}

func torrentDownloadUrl(url string) (Url string, err error) {
	resp, err := client(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.WithMessage(err, "torrentDownloadUrl goquery")
	}

	doc.Find("#body > div > dl > dd:nth-child(14) > a").Each(func(i int, selection *goquery.Selection) {

		u, ok := selection.Attr("href")
		if ok {
			Url = urlParser(u, url)
		}
	})
	return
}

func getMagnet(url string) (magnet string, err error) {

	resp, err := client(url)
	if err != nil {
		return "", errors.WithMessagef(err, "BTBT: 获取磁链错误 %s", url)
	}
	defer resp.Body.Close()
	return magnetConvert.IO2Magnet(resp.Body)

}

func client(Url string) (resp *http.Response, err error) {
	c := httpClient.NewHttpClient()
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}
	resp, err = c.Do(req)
	return
}
