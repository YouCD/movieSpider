package feedspider

import (
	"bytes"
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"strings"
	"sync"

	"github.com/youcd/toolkit/log"

	"github.com/PuerkitoBio/goquery"
)

type Rarbg2 struct {
	BaseFeeder
	typ     types.VideoType
	webHost string
}

func NewRarbg2(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	parse, err := url.Parse(siteURL)
	if err != nil {
		panic(err)
	}
	return &Rarbg2{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Scheduling: scheduling,
				Url:        siteURL,
				UseIPProxy: useIPProxy,
			},
			web: "rarbg2",
		},
		typ:     resourceType,
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}

func (r *Rarbg2) Crawler() ([]*types.FeedVideoBase, error) {
	resp, err := r.HTTPRequest(r.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", r.web, r.Url, err)
	}
	var videosArr []*types.FeedVideoBase
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}
	doc.Find("#content > table > tbody > tr").Each(func(_ int, selection *goquery.Selection) {
		name, exists := selection.Find("td.td.tdnormal.hide > div > a > button").Attr("title")
		if !exists {
			return
		}
		downloadURI, exists := selection.Find("td.td.tdnormal.hide > div > a").Attr("href")
		if !exists {
			return
		}

		URLStr := fmt.Sprintf("%s%s", r.webHost, downloadURI)

		videosArr = append(videosArr, &types.FeedVideoBase{
			TorrentName: name,
			TorrentURL:  URLStr,
			Type:        r.typ.String(),
			Web:         r.web,
		})
	})
	var videos []*types.FeedVideoBase
	var wg sync.WaitGroup
	for _, videoBase := range videosArr {
		wg.Add(1)
		go func(videoBase *types.FeedVideoBase) {
			defer wg.Done()
			magnet, err := r.moviePageURL(videoBase.TorrentURL)
			if err != nil {
				log.Warnf("rarbg2: %s", err)
				return
			}
			videos = append(videos, &types.FeedVideoBase{
				TorrentName: videoBase.TorrentName,
				TorrentURL:  videoBase.TorrentURL,
				Magnet:      magnet,
				Type:        videoBase.Type,
				Web:         videoBase.Web,
			})
		}(videoBase)
	}
	wg.Wait()
	return videos, nil
}

func (r *Rarbg2) moviePageURL(pageURL string) (string, error) {
	resp, err := r.HTTPRequest(pageURL)
	if err != nil {
		return "", fmt.Errorf("连接请求: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", fmt.Errorf("moviePageURL: %w", err)
	}
	magnet := ""
	doc.Find("#content > table > tbody > tr > td > div > div > a ").Each(func(_ int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if !exists {
			return
		}
		if strings.Contains(val, "magnet") {
			magnet = val
			return
		}
	})

	return magnet, nil
}
