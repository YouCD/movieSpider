package feedspider

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	httpClient2 "movieSpider/internal/httpclient"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	urlBaseGlodls   = "http://gtso.cc"
	urlRssURIGlodls = "rss.php?cat=1,41"
)

type Glodls struct {
	urlBase string
	BaseFeeder
}

func NewGlodls(scheduling, mirrorSite string) *Glodls {
	urlBase := urlBaseGlodls
	urlStr := fmt.Sprintf("%s/%s", urlBaseGlodls, urlRssURIGlodls)
	if mirrorSite != "" {
		urlStr = fmt.Sprintf("%s/%s", mirrorSite, urlRssURIGlodls)
		urlBase = mirrorSite
	}

	//nolint:exhaustruct
	return &Glodls{
		urlBase,
		BaseFeeder{
			web:        "glodls",
			url:        urlStr,
			scheduling: scheduling,
		},
	}
}

//nolint:gosimple,ineffassign,goconst
func (g *Glodls) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(g.url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(g.web), fd.String())
	//nolint:prealloc
	var videosA []*types.FeedVideo
	for _, v := range fd.Items {
		// 片名
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		// 片名处理
		var name, year string
		if strings.ToLower(v.Categories[0]) == "tv" {
			compileRegex := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(matchArr) == 0 {
				continue
			}
			//nolint:wastedassign
			name = matchArr[1]
		} else if strings.ToLower(v.Categories[0]) == "movies" {
			compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			if len(matchArr) == 0 {
				//nolint:wastedassign
				name = torrentName
			} else {
				//nolint:wastedassign
				name = matchArr[1]
				year = matchArr[2]
			}
		}
		name = torrentName

		fVideo := new(types.FeedVideo)
		fVideo.Name = fVideo.FormatName(name)
		fVideo.Year = year

		fVideo.Web = g.web
		parse, _ := url.Parse(v.Link)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)

		if len(parse.Query()["id"]) == 0 {
			log.Error("没有ID")
		}
		id := parse.Query()["id"][0]
		all := strings.ReplaceAll(v.Title, " ", "-")

		TorrentURL := fmt.Sprintf("%s/%s-f-%s.html", g.urlBase, strings.ToLower(all), id)

		fVideo.TorrentURL = TorrentURL

		// 处理 资源类型 是 电影 还是电视剧
		typ := strings.ToLower(v.Categories[0])
		if typ == "movies" {
			fVideo.Type = "movie"
		} else {
			fVideo.Type = typ
		}
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		//nolint:exhaustruct
		fVideo.RowData = sql.NullString{String: string(bytes)}

		videosA = append(videosA, fVideo)
	}
	var wg sync.WaitGroup
	for _, v := range videosA {
		wg.Add(1)
		go func(video *types.FeedVideo) {
			defer wg.Done()
			magnet, err := g.fetchMagnet(video.TorrentURL)
			if err != nil {
				log.Error(err)
			}
			if magnet == "" {
				return
			}
			video.Magnet = magnet
			videos = append(videos, video)
		}(v)
	}
	wg.Wait()
	//nolint:nakedret
	return
}

func (g *Glodls) fetchMagnet(url string) (magnet string, err error) {
	request, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "GLODLS: 创建新的请求")
	}
	httpClient := httpClient2.NewHTTPClient()
	httpClient.Timeout = 20 * time.Second
	resp, err := httpClient.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "GLODLS: 请求错误")
	}
	if resp == nil {
		return "", errors.New("GLODLS: response is nil")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "GLODLS: Document 查找出错")
	}
	selector := "#downloadbox > table > tbody > tr > td:nth-child(1) > a:nth-child(2)"
	magnet, exists := doc.Find(selector).Attr("href")
	if !exists {
		return "", errors.Wrap(err, "GLODLS: 查找href出错")
	}
	return magnet, nil
}
