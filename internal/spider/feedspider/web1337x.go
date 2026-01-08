package feedspider

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Web1337x struct {
	BaseFeeder
	typ     types.VideoType
	webHost string
}

func NewWeb1337x(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	parse, _ := url.Parse(siteURL)
	return &Web1337x{
		typ: resourceType,
		BaseFeeder: BaseFeeder{
			web:      "1337x",
			BaseFeed: types.BaseFeed{Url: siteURL, Scheduling: scheduling, UseIPProxy: useIPProxy},
		},
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}

func (w *Web1337x) Crawler(ctx context.Context) ([]*types.FeedVideoBase, error) {
	log.WithCtx(ctx).Debugf("%s type: %v url: %s", w.web, w.typ, w.Url)
	videosTemp := make([]*types.FeedVideoBase, 0)
	data, err := w.HTTPRequest(ctx, w.Url)
	if err != nil {
		return nil, fmt.Errorf("fetchHTMLData, err:%w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery, err:%w", err)
	}

	selector := "body > main > div > div > div.featured-list > div > table > tbody > tr"

	compileRegex := regexp.MustCompile(`/torrent/\d+/(.*)/.*`)
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		var torrentURL, magnet, typ string
		// href
		href, exists := s.Find("td.coll-1.name > a:nth-child(2)").Attr("href")
		if !exists {
			return
		}
		torrentURL = fmt.Sprintf("%s%s", w.webHost, href)
		log.WithCtx(ctx).Debugf("1337x.%s %s", w.typ, torrentURL)
		matchArr := compileRegex.FindStringSubmatch(href)
		if len(matchArr) <= 1 {
			return
		}
		// magnet 链接
		var count int
	RETRY:
		data, err = w.HTTPRequest(ctx, torrentURL)
		if err != nil {
			if count < 3 {
				count++
				if count == 3 {
					log.WithCtx(ctx).Warnf("RETRY: 1337x.%s retry count:%d, error:%s", w.typ, count, err)
				}
				goto RETRY
			}
			return
		}
		count = 0
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
		if err != nil {
			log.WithCtx(ctx).Error(err)
			return
		}

		// #openPopup
		magnet, exists = doc.Find(".torrentdown1").Attr("href")
		if !exists {
			return
		}
		switch w.typ {
		case types.VideoTypeMovie:
			typ = "movie"
		case types.VideoTypeTV:
			typ = "tv"
		}

		video := &types.FeedVideoBase{
			TorrentName: matchArr[1],
			TorrentURL:  torrentURL,
			Magnet:      magnet,
			Type:        typ,
			RowData:     sql.NullString{},
			Web:         w.web,
		}
		videosTemp = append(videosTemp, video)
	})
	return videosTemp, nil
}
