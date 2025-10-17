package feedspider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"strings"
	"sync"

	"github.com/youcd/toolkit/log"

	"github.com/PuerkitoBio/goquery"
)

//nolint:tagliatelle
type theRarbgItem struct {
	URL        string `json:"url"`
	Thumbnail  string `json:"thumbnail"`
	Rating     string `json:"rating"`
	Imdb       string `json:"imdb"`
	Name       string `json:"name"`
	DetailPage string `json:"detail_page"`
	Poster     string `json:"poster"`
}
type TheRarbg struct {
	BaseFeeder
	typ     types.VideoType
	webHost string
}

func NewTheRarbg(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	parse, err := url.Parse(siteURL)
	if err != nil {
		panic(err)
	}
	return &TheRarbg{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Scheduling: scheduling,
				Url:        siteURL,
				UseIPProxy: useIPProxy,
			},
			web: "the_rarbg",
		},
		typ:     resourceType,
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}

func (r *TheRarbg) Crawler() ([]*types.FeedVideoBase, error) {
	resp, err := r.HTTPRequest(r.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", r.web, r.Url, err)
	}

	var result []*theRarbgItem
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("%s json unmarshal, err: %w", r.web, err)
	}

	videosArr := make([]*types.FeedVideoBase, 0)
	for _, item := range result {
		URLStr := fmt.Sprintf("%s%s", r.webHost, item.DetailPage)
		videosArr = append(videosArr, &types.FeedVideoBase{
			TorrentName: item.Name,
			TorrentURL:  URLStr,
			Type:        r.typ.String(),
			Web:         r.web,
		})
	}
	var videos []*types.FeedVideoBase
	var wg sync.WaitGroup
	for _, videoBase := range videosArr {
		wg.Add(1)
		go func(videoBase *types.FeedVideoBase) {
			defer wg.Done()
			magnet, err := r.moviePageURL(videoBase.TorrentURL)
			if err != nil {
				log.Warnf("the_rarbg: %s", err)
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
	log.Infof("the_rarbg: %v", len(videos))
	return videos, nil
}

func (r *TheRarbg) moviePageURL(pageURL string) (string, error) {
	resp, err := r.HTTPRequest(pageURL)
	if err != nil {
		return "", fmt.Errorf("连接请求: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", fmt.Errorf("moviePageURL: %w", err)
	}
	magnet := ""
	selectorA := "body > div.topnav > div:nth-child(5) > div.postContL.col-12.col-md-9.col-lg-11 > div.table-responsive > table > tbody > tr:nth-child(2) > td > div > div.download-primary > a.btn-download.magnet-btn"
	selectorB := "body > div.topnav > div:nth-child(4) > div.postContL.col-12.col-md-9.col-lg-11 > div.table-responsive > table > tbody > tr:nth-child(2) > td > div > div.download-primary > a.btn-download.magnet-btn"
	magnet = r.getMagnet(doc, selectorA)
	if magnet == "" {
		magnet = r.getMagnet(doc, selectorB)
	}
	return magnet, nil
}

func (r *TheRarbg) getMagnet(doc *goquery.Document, selector string) string {
	var magnet string
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if !exists {
			return
		}
		if strings.Contains(val, "magnet") {
			magnet = val
			return
		}
	})
	return magnet
}
