package feedspider

import (
	"bytes"
	"context"
	"fmt"
	"movieSpider/internal/types"
	"net/url"

	"github.com/youcd/toolkit/log"

	"github.com/PuerkitoBio/goquery"
)

type Ilcorsaronero struct {
	BaseFeeder

	webHost string
	typ     types.VideoType
}

func NewIlcorsaronero(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	parse, _ := url.Parse(siteURL)
	return &Ilcorsaronero{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Scheduling: scheduling,
				Url:        siteURL,
				UseIPProxy: useIPProxy,
			},
			web: "Ilcorsaronero",
		},
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
		typ:     resourceType,
	}
}
func (u *Ilcorsaronero) Crawler() ([]*types.FeedVideoBase, error) {
	log.WithCtx(context.Background()).Debugf("%s type: %v url: %s", u.web, u.typ, u.Url)
	// body > main > div.container.md\:rounded-xl.md\:shadow.md\:border.bg-neutral-800.border-neutral-900.text-neutral-400 > div.overflow-x-auto > table > tbody > tr:nth-child(1)
	resp, err := u.HTTPRequest(u.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", u.web, u.Url, err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}
	var videosArr []*types.FeedVideoBase

	selector := `body > main > div.container.md\:rounded-xl.md\:shadow.md\:border.bg-neutral-800.border-neutral-900.text-neutral-400 > div.overflow-x-auto > table > tbody > tr`
	doc.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		// body > main > div.container.md\:rounded-xl.md\:shadow.md\:border.bg-neutral-800.border-neutral-900.text-neutral-400 > div.overflow-x-auto > table > tbody > tr:nth-child(1) > th > div > a
		// body > main > div.container.md\:rounded-xl.md\:shadow.md\:border.bg-neutral-800.border-neutral-900.text-neutral-400 > div.overflow-x-auto > table > tbody > tr:nth-child(2) > th > div > a
		href := selection.Find("th > div > a").AttrOr("href", "")
		video := types.FeedVideoBase{
			TorrentName: selection.Find("th > div > a").Text(),
			TorrentURL:  fmt.Sprintf("%s%s", u.webHost, href),
			Web:         u.web,
		}
		videosArr = append(videosArr, &video)
	})

	return u.fetchMagnetDownLoad(videosArr), nil
}

func (u *Ilcorsaronero) fetchMagnetDownLoad(videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	s := `body > main > div.w-full.max-w-screen-xl.mx-auto.p-8.md\:rounded-xl.md\:shadow.md\:border.bg-neutral-800.border-neutral-900.text-neutral-300 > div.flex.flex-wrap.items-center.gap-4 > a.w-full.sm\:w-auto.px-5.py-2\.5.rounded-xl.text-sm.text-center.text-black.bg-neutral-300.hover\:bg-neutral-200.focus\:ring-4.focus\:ring-neutral-100.focus\:outline-none`
	return fetchMagnetDownLoad(u.BaseFeeder, s, videos)
}
