package feedspider

import (
	"context"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"strings"

	"github.com/youcd/toolkit/log"
)

type ThePirateBay struct {
	BaseFeeder
}

func NewThePirateBay() *ThePirateBay {
	return &ThePirateBay{
		BaseFeeder{
			web: "ThePirateBay",
			BaseFeed: types.BaseFeed{
				Url:        config.Config.Feed.ThePirateBay.Url,
				Scheduling: config.Config.Feed.ThePirateBay.Scheduling,
				UseIPProxy: config.Config.Feed.ThePirateBay.UseIPProxy,
			},
		},
	}
}
func (t *ThePirateBay) Crawler() ([]*types.FeedVideoBase, error) {
	fd, err := t.FeedParser().ParseURL(t.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.WithCtx(context.Background()).Debugf("%s Data: %s", t.web, fd.String())
	videos := make([]*types.FeedVideoBase, 0)
	for _, v := range fd.Items {
		if len(v.Categories) < 1 {
			continue
		}
		video := new(types.FeedVideoBase)
		if strings.Contains(strings.ToLower(v.Categories[0]), "movie") {
			video.Type = types.VideoTypeMovie.String()
		}
		if strings.Contains(strings.ToLower(v.Categories[0]), "tv") {
			video.Type = types.VideoTypeTV.String()
		}

		video.TorrentName = v.Title
		video.Magnet = v.Link
		video.TorrentURL = v.GUID
		video.Web = t.web
		videos = append(videos, video)
	}
	return videos, nil
}
