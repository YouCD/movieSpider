package feedspider

import (
	"context"
	"strings"
	"sync"

	"movieSpider/internal/types"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

func fetchMagnetDownLoad(ctx context.Context, baseFeeder BaseFeeder, selector string, videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	var wg sync.WaitGroup
	var videos2 []*types.FeedVideoBase
	for _, video := range videos {
		wg.Add(1)
		//nolint:noctx
		go func() {
			defer wg.Done()
			resp, err := baseFeeder.HTTPClientDynamic(ctx).Get(video.TorrentURL)
			if err != nil {
				log.WithCtx(ctx).Errorf("%s.%s %s http request url is %s, error:%s", baseFeeder.web, video.Type, video.TorrentName, video.TorrentURL, err)
				return
			}
			defer resp.Body.Close()
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.WithCtx(ctx).Errorf("%s.%s %s http goquery error:%s", baseFeeder.web, video.Type, video.TorrentName, err)
				return
			}
			val, exists := doc.Find(selector).Attr("href")
			if exists {
				video.Magnet = strings.ReplaceAll(val, "\n        ", "")
				videos2 = append(videos2, video)
			}
		}()
	}
	wg.Wait()
	return videos2
}
