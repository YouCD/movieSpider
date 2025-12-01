package feedspider

import (
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestNewTheRarbg(t *testing.T) {
	for _, r := range config.Config.Feed.TheRarbg {
		if r.ResourceType == types.VideoTypeTV {
			feeder := NewTheRarbg(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
			videos, err := feeder.Crawler()
			if err != nil {
				t.Error(err)
			}
			for _, video := range videos {
				log.WithCtx(context.Background()).Infof("%#v", video)
				//filterVideo, err := model.FilterVideo(video)
				//if err != nil {
				//	if errors.Is(err, model.ErrFeedVideoExclude) {
				//		log.WithCtx(context.Background()).Warnf("err: %s    %#v", err, video)
				//		continue
				//	}
				//	continue
				//}
				//log.WithCtx(context.Background()).Infof("%#v", filterVideo)
			}
		}
	}
}
