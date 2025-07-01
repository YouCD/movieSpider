package feedspider

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Glodls struct {
	urlBase string
	BaseFeeder
}

func NewGlodls() *Glodls {
	parse, err := url.Parse(config.Config.Feed.GLODLS.Url)
	if err != nil {
		log.Errorf("url.Parse err: %v", err)
		return nil
	}

	urlBase := parse.Scheme + "://" + parse.Host
	return &Glodls{
		urlBase,
		BaseFeeder{
			web: "glodls",
			BaseFeed: types.BaseFeed{
				Url:        config.Config.Feed.GLODLS.Url,
				Scheduling: config.Config.Feed.GLODLS.Scheduling,
				UseIPProxy: config.Config.Feed.GLODLS.UseIPProxy,
			},
		},
	}
}

//nolint:goconst
func (g *Glodls) Crawler() (videos []*types.FeedVideoBase, err error) {
	fd, err := g.FeedParser().ParseURL(g.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(g.web), fd.String())
	//nolint:prealloc
	var videosA []*types.FeedVideoBase
	for _, v := range fd.Items {
		fVideo := new(types.FeedVideoBase)
		fVideo.Web = g.web
		parse, _ := url.Parse(v.Link)
		// 种子名
		fVideo.TorrentName = v.Title

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

		fVideo.RowData = sql.NullString{String: string(bytes)}
		fVideo.Web = g.web
		videosA = append(videosA, fVideo)
	}
	var wg sync.WaitGroup
	for _, v := range videosA {
		wg.Add(1)
		go func(video *types.FeedVideoBase) {
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
		return "", fmt.Errorf("GLODLS: 请求错误，err:%w", err)
	}

	resp, err := g.HTTPClientDynamic().Do(request)
	if err != nil {
		return "", fmt.Errorf("GLODLS: 请求错误，err:%w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("GLODLS: response is nil，err:%w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("GLODLS: Document 查找出错，err:%w", err)
	}
	selector := "#downloadbox > table > tbody > tr > td:nth-child(1) > a:nth-child(2)"
	magnet, exists := doc.Find(selector).Attr("href")
	if !exists {
		return "", fmt.Errorf("GLODLS: 查找href出错，url:%s", url)
	}
	return magnet, nil
}
