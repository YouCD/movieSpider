package feedspider

import (
	"bytes"
	"context"
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type TgxWeb struct {
	webHost string
	BaseFeeder
	typ types.VideoType
}

func NewTgxWeb(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	parse, err := url.Parse(siteURL)
	if err != nil {
		panic(err)
	}
	return &TgxWeb{
		BaseFeed: types.BaseFeed{
			Scheduling: scheduling,
			Url:        siteURL,
			UseIPProxy: useIPProxy,
		},
		web:     "tgx",
		typ:     resourceType,
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}

func (t *TgxWeb) Crawler(ctx context.Context) ([]*types.FeedVideoBase, error) {
	log.WithCtx(ctx).Debugf("%s type: %v url: %s", t.web, t.typ, t.Url)
	resp, err := t.HTTPRequest(ctx, t.Url)
	if err != nil {
		return nil, fmt.Errorf("%s new request,url: %s, err: %w", t.web, t.Url, err)
	}
	var videosArr []*types.FeedVideoBase
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}

	selector := "#panelmain > div.panel-body > div.container-fluid > div > div"

	doc.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		var video types.FeedVideoBase
		video.TorrentName = selection.Find("#click > div > a.txlight").Text()
		video.TorrentURL = fmt.Sprintf("%s/%s", t.webHost, selection.Find("#click > div > a.txlight").AttrOr("href", ""))
		video.Type = t.typ.String()
		video.Web = t.web
		videosArr = append(videosArr, &video)
	})
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		videosList []*types.FeedVideoBase
	)
	for _, videoBase := range videosArr {
		wg.Add(1)
		go func(v *types.FeedVideoBase) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			resp, err := t.HTTPRequest(ctx, v.TorrentURL)
			if err != nil {
				log.WithCtx(ctx).Warnf("tgx: %s", err)
				return
			}
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
			if err != nil {
				log.WithCtx(ctx).Warnf("tgx: %s", err)
				return
			}
			v.Magnet = doc.Find("#covercell > center > a.btn.btn-danger.lift.txlight").AttrOr("href", "")
			mu.Lock()
			videosList = append(videosList, v)
			mu.Unlock()
		}(videoBase)
	}
	wg.Wait()
	return videosList, nil
}
