package feedspider

import (
	"bytes"
	"context"
	"movieSpider/internal/config"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"sync"

	"github.com/youcd/toolkit/log"
)

type Yts struct {
	BaseFeeder
}

func NewYts() *Yts {
	return &Yts{BaseFeeder{
		web: "Yts",
		BaseFeed: types.BaseFeed{
			Scheduling: config.Config.Feed.Yts.Scheduling,
			Url:        config.Config.Feed.Yts.Url,
			UseIPProxy: config.Config.Feed.Yts.UseIPProxy,
		},
	}}
}

func (f *Yts) Crawler(ctx context.Context) (videos []*types.FeedVideoBase, err error) {
	fd, err := f.FeedParser(ctx).ParseURL(f.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.WithCtx(ctx).Debugf("%s Data: %s", f.web, fd.String())
	var tempVideos []*types.FeedVideoBase
	for _, v := range fd.Items {
		fVideo := new(types.FeedVideoBase)
		fVideo.Web = f.web
		fVideo.TorrentName = v.Title
		fVideo.TorrentURL = v.Link
		fVideo.Type = types.VideoTypeMovie.String()
		if len(v.Enclosures) < 1 {
			continue
		}
		fVideo.Magnet = v.Enclosures[0].URL
		tempVideos = append(tempVideos, fVideo)
	}

	var wg sync.WaitGroup
	for _, video := range tempVideos {
		wg.Add(1)
		go func(video *types.FeedVideoBase) {
			defer wg.Done()
			do, err := f.HTTPRequest(ctx, video.Magnet)
			if err != nil {
				log.WithCtx(ctx).Panicw("Yts Magnet Request Error", "err", err)
				return
			}
			magnet, err := magnetconvert.IO2Magnet(bytes.NewReader(do))
			if err != nil {
				log.WithCtx(ctx).Panicw("Yts Magnet Convert Error", "err", err)
				return
			}
			video.Magnet = magnet
			videos = append(videos, video)
		}(video)
	}
	wg.Wait()
	return
}
