package feedspider

import (
	"context"
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/log"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const urlBtbt = "https://www.btbtt12.com/forum-index-fid-951.htm"

type Btbt struct {
	BaseFeeder
}

func NewBtbt(scheduling string) *Btbt {
	//nolint:forcetypeassert
	return &Btbt{
		BaseFeeder{
			web:        "btbt",
			url:        urlBtbt,
			scheduling: scheduling,
		},
	}
}
func (b *Btbt) Crawler() (videos []*types.FeedVideo, err error) {
	c := httpclient.NewHTTPClient()
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, b.url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "btbt new request")
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
					fVideo.TorrentURL = urlParser(u, b.url)
					fVideo.Name = name
				}
			})

			// 年份
			selection.Find("a:nth-child(2)").Each(func(i int, selection *goquery.Selection) {
				fVideo.Year = trim(selection.Text())
			})
			//nolint:goconst
			fVideo.Type = "movie"
			Videos1 = append(Videos1, fVideo)
		})
	})

	// 洗出  下载电影种子下载的页面
	var wg sync.WaitGroup
	//nolint:prealloc
	var Videos2 []*types.FeedVideo
	for _, v := range Videos1 {
		parse, err := url.Parse(v.TorrentURL)
		if err != nil {
			log.Error(err)
			continue
		}
		if parse.String() == "" {
			continue
		}
		wg.Add(1)
		downloadURL, err := moviePageURL(v.TorrentURL)
		if err != nil {
			log.Error(err)
		}
		if downloadURL == "" {
			wg.Done()
			continue
		}
		v.TorrentName = v.Name
		//nolint:exhaustruct
		v.RowData = sql.NullString{String: downloadURL}
		Videos2 = append(Videos2, v)
		wg.Done()
	}
	wg.Wait()

	// 洗出  磁力连接
	var wg1 sync.WaitGroup
	for _, v := range Videos2 {
		wg1.Add(1)
		torrentDownloadURLStr, err := torrentDownloadURL(v.RowData.String)
		if err != nil {
			return nil, err
		}
		magnet, err := getMagnet(torrentDownloadURLStr)
		if err != nil {
			log.Error(err)
		}
		v.Magnet = magnet
		v.Web = "btbt"
		videos = append(videos, v)
		wg1.Done()
	}
	wg1.Wait()
	//nolint:nakedret
	return
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
	}
	return parse.String()
}

func moviePageURL(pageURL string) (url string, err error) {
	resp, err := client(pageURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.WithMessagef(err, "BTBT: 电影页没有数据: %s", pageURL)
	}

	doc.Find("#body > div > table:nth-child(2) > tbody > tr:nth-child(1) > td.post_td > div.attachlist > table > tbody > tr:nth-child(3) > td:nth-child(1) > a").Each(func(i int, selection *goquery.Selection) {
		u, ok := selection.Attr("href")
		if ok {
			url = urlParser(u, pageURL)
		}
	})
	return
}

func torrentDownloadURL(url string) (newURL string, err error) {
	resp, err := client(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.WithMessage(err, "torrentDownloadURL goquery")
	}

	doc.Find("#body > div > dl > dd:nth-child(14) > a").Each(func(i int, selection *goquery.Selection) {
		u, ok := selection.Attr("href")
		if ok {
			newURL = urlParser(u, url)
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
	//nolint:wrapcheck
	return magnetconvert.IO2Magnet(resp.Body)
}

func client(url string) (resp *http.Response, err error) {
	c := httpclient.NewHTTPClient()
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "client")
	}
	resp, err = c.Do(req)
	return
}
