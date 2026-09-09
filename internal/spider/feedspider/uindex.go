package feedspider

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"movieSpider/internal/types"

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

func (u *Uindex) Crawler(ctx context.Context) ([]*types.FeedVideoBase, error) {
	log.WithCtx(ctx).Debugf("%s type: %v url: %s", u.web, u.typ, u.Url)
	resp, err := u.HTTPRequest(ctx, u.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", u.web, u.Url, err)
	}
	var videosArr []*types.FeedVideoBase
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}
	selector := "#content > div > div.top-table-container > table > tbody > tr"
	doc.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		var video types.FeedVideoBase
		// #content > div > div.top-table-container > table > tbody > tr:nth-child(1) > td.sr-col-name > a.sr-torrent-link
		video.TorrentName = selection.Find("td.sr-col-name > a.sr-torrent-link").Text()

		video.TorrentURL = fmt.Sprintf("%s/%s", u.webHost, selection.Find("td>a:nth-child(2)").AttrOr("href", ""))

		video.Type = u.typ.String()
		video.Web = u.web
		videosArr = append(videosArr, &video)
	})
	return u.fetchMagnetDownLoad(ctx, videosArr), nil
}

func (u *Uindex) fetchMagnetDownLoad(ctx context.Context, videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	return fetchMagnetDownLoad(ctx, u.BaseFeeder, "#content > div.dt-download-box > a", videos)
}
