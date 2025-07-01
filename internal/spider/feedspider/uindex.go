package feedspider

import (
	"bytes"
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Uindex struct {
	webHost string
	typ     types.VideoType
	BaseFeeder
}

func NewUindex(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	parse, _ := url.Parse(siteURL)
	return &Uindex{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Scheduling: scheduling,
				Url:        siteURL,
				UseIPProxy: useIPProxy,
			},
			web: "Uindex",
		},
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
		typ:     resourceType,
	}
}
func (u *Uindex) Crawler() ([]*types.FeedVideoBase, error) {
	resp, err := u.HTTPRequest(u.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", u.web, u.Url, err)
	}
	var videosArr []*types.FeedVideoBase
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}

	selector := "#content > table > tbody > tr"
	doc.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		var video types.FeedVideoBase
		video.TorrentName = selection.Find("td>a:nth-child(2)").Text()
		video.TorrentURL = fmt.Sprintf("%s/%s", u.webHost, selection.Find("td>a:nth-child(2)").AttrOr("href", ""))
		video.Type = u.typ.String()
		video.Web = u.web
		videosArr = append(videosArr, &video)

	})
	return u.fetchMagnetDownLoad(videosArr), nil
}

func (u *Uindex) fetchMagnetDownLoad(videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	var wg sync.WaitGroup
	var videos2 []*types.FeedVideoBase
	for _, video := range videos {
		wg.Add(1)
		//nolint:noctx
		go func() {
			defer wg.Done()
			resp, err := u.HTTPClientDynamic().Get(video.TorrentURL)
			if err != nil {
				log.Errorf("Uindex.%s %s http request url is %s, error:%s", video.Type, video.TorrentName, video.TorrentURL, err)
				return
			}
			defer resp.Body.Close()
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Errorf("Uindex.%s %s http goquery error:%s", video.Type, video.TorrentName, err)
				return
			}
			val, exists := doc.Find("#downloadbox > h2 > a").Attr("href")
			if exists {
				video.Magnet = strings.ReplaceAll(val, "\n        ", "")
				videos2 = append(videos2, video)
			}
		}()
	}
	wg.Wait()
	return videos2
}
