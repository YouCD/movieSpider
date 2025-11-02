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
	log.Debugw(r.web, "type", r.typ, "url", r.Url)
	resp, err := r.HTTPRequest(r.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", r.web, r.Url, err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("moviePageURL: %w", err)
	}
	var urls []string
	doc.Find("tbody > tr > td.cellName > div >  a:nth-child(1)").Each(func(_ int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if !exists {
			return
		}
		if strings.Contains(val, "post-detail") {
			urls = append(urls, fmt.Sprintf("%s%s", r.webHost, val))
			return
		}
	})

	var videos []*types.FeedVideoBase
	var wg sync.WaitGroup
	log.Debugw(r.web, "type", r.typ, "urls", len(urls))
	for _, urlItem := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			name, magnet, err := r.moviePageURL(u)
			if err != nil {
				log.Warnf("the_rarbg: %s", err)
				return
			}

			videos = append(videos, &types.FeedVideoBase{
				TorrentName: name,
				TorrentURL:  urlItem,
				Magnet:      magnet,
				Type:        r.typ.String(),
				Web:         r.web,
			})
		}(urlItem)
	}
	wg.Wait()
	log.Debugw(r.web, "type", r.typ, "videos", len(videos))
	return videos, nil
}

func (r *TheRarbg) moviePageURL(pageURL string) (string, string, error) {
	resp, err := r.HTTPRequest(pageURL)
	if err != nil {
		return "", "", fmt.Errorf("连接请求: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", "", fmt.Errorf("moviePageURL: %w", err)
	}

	selectorA := "body > div.topnav > div:nth-child(5) > div.postContL.col-12.col-md-9.col-lg-11 > div.table-responsive > table > tbody > tr:nth-child(2) > td > div > div.download-primary > a.btn-download.magnet-btn"
	selectorB := "body > div.topnav > div:nth-child(4) > div.postContL.col-12.col-md-9.col-lg-11 > div.table-responsive > table > tbody > tr:nth-child(2) > td > div > div.download-primary > a.btn-download.magnet-btn"
	selectorC := "h4.text-center.m-4"
	magnet := r.getMagnet(doc, selectorA)
	if magnet == "" {
		magnet = r.getMagnet(doc, selectorB)
	}

	var name string

	doc.Find(selectorC).Each(func(_ int, s *goquery.Selection) {
		name = s.Text()
	})
	return name, magnet, nil
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
