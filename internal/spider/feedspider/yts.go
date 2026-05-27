package feedspider

import (
	"bytes"
	"context"
	"sync"

	"movieSpider/internal/config"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"

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
	var mu sync.Mutex // 保护 videos 切片的并发访问
	for _, video := range tempVideos {
		wg.Add(1)
		go func(video *types.FeedVideoBase) {
			defer wg.Done()
			do, err := f.HTTPRequest(ctx, video.Magnet)
			if err != nil {
				log.WithCtx(ctx).Errorw("Yts Magnet Request Error", "err", err, "torrent_url", video.TorrentURL)
				return
			}
			magnet, err := magnetconvert.IO2Magnet(bytes.NewReader(do))
			if err != nil {
				log.WithCtx(ctx).Errorw("Yts Magnet Convert Error", "err", err, "torrent_url", video.TorrentURL)
				return
			}
			video.Magnet = magnet
			mu.Lock()
			videos = append(videos, video)
			mu.Unlock()
		}(video)
	}
	wg.Wait()
	return
}
