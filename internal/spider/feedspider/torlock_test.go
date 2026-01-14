package feedspider

import (
	"context"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func Test_torlock_Crawler(t *testing.T) {
	var err error
	var videos []*types.FeedVideoBase
	for _, r := range config.Config.Feed.TORLOCK {
		if r != nil {
			if r.ResourceType == types.VideoTypeTV {
				feedTorlockTV := NewTorlock(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
				videos, err = feedTorlockTV.Crawler(context.Background())
				if err != nil {
					log.WithCtx(context.Background()).Errorf("err: %s", err)
					return
				}
			}
			//if r.ResourceType == types.VideoTypeMovie {
			//	videos, err = NewTorlock(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy).Crawler(context.Background())
			//	if err != nil {
			//		log.WithCtx(context.Background()).Errorf("err: %s", err)
			//		return
			//	}
			//}
			//log.Debug(r)
		}
	}

	for _, video := range videos {
		//filterVideo, err := model.FilterVideo(video)
		//if err != nil {
		//	log.WithCtx(context.Background()).Errorf("err: %s    %#v", err, video)
		//	continue
		//}
		log.WithCtx(context.Background()).Infof("%#v", video)
	}
}
