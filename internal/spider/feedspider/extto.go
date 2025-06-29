package feedspider

import (
	"bytes"
	"database/sql"
	"fmt"
	"movieSpider/internal/types"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Extto struct {
	BaseFeeder
	webHost string
	typ     types.VideoType
}

func NewExtto(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy, useCloudflareBypass bool) Feeder {
	parse, _ := url.Parse(siteURL)
	return &Extto{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Scheduling:          scheduling,
				Url:                 siteURL,
				UseIPProxy:          useIPProxy,
				UseCloudflareBypass: useCloudflareBypass,
			},
			web: "extto",
		},
		typ:     resourceType,
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}

func (b *Extto) Crawler() (videos []*types.FeedVideoBase, err error) {
	log.Debug(b.Url)

	resp, err := b.HTTPRequest(b.Url)
	if err != nil {
		return nil, fmt.Errorf("extto new request,url: %s, err: %w", b.Url, err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("goquery,err: %w", err)
	}

	doc.Find("body > div.container-fluid > div > div > div.col-12.col-md-12.col-xl-10.py-md-3.bd-content.main-block > div > div > div > form > div.col-md-10.result-list-block > div.found-content-block > div > table > tbody > tr").Each(func(_ int, s *goquery.Selection) {
		var torrentName, Magnet, TorrentURL string
		s.Find("td.text-left > div.float-left > a > b").Each(func(_ int, selection *goquery.Selection) {
			torrentName = selection.Text()
		})

		s.Find(".torrent-dwn").Each(func(_ int, selection *goquery.Selection) {
			val, exists := selection.Attr("href")
			if !exists {
				return
			}
			Magnet = val
		})

		var selector string
		if b.typ == types.VideoTypeMovie {
			selector = "td.text-left > div.float-left.has-movie > a"
		} else {
			selector = "td.text-left > div.float-left > a"
		}

		s.Find(selector).Each(func(_ int, selection *goquery.Selection) {
			href, exists := selection.Attr("href")
			if !exists {
				return
			}
			TorrentURL = fmt.Sprintf("%s%s", b.webHost, href)
		})

		videos = append(videos, &types.FeedVideoBase{
			TorrentName: torrentName,
			TorrentURL:  TorrentURL,
			Magnet:      Magnet,
			Type:        b.typ.String(),
			RowData:     sql.NullString{},
			Web:         b.web,
		})

	})

	return videos, nil
}
