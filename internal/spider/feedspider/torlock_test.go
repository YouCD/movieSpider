package feedspider

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func Test_torlock_Crawler(t *testing.T) {
	//facTORLOCK.CreateFeeder("*/1 * * * *", types.VideoTypeTV).Run()
	//select {}
	var err error
	var videos []*types.FeedVideoBase
	for _, r := range config.Config.Feed.TORLOCK {
		if r != nil {
			//if r.ResourceType == types.VideoTypeTV {
			//	feedTorlockTV := NewTorlock(r.Scheduling, r.ResourceType, r.URL, r.UseIPProxy)
			//	videos, err = feedTorlockTV.Crawler()
			//	if err != nil {
			//		log.Errorf("err: %s", err)
			//		return
			//	}
			//}
			if r.ResourceType == types.VideoTypeMovie {
				videos, err = NewTorlock(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy).Crawler()
				if err != nil {
					log.Errorf("err: %s", err)
					return
				}
			}
			//log.Debug(r)
		}
	}

	for _, video := range videos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			log.Errorf("err: %s    %#v", err, video)
			continue
		}
		log.Infof("%#v", filterVideo)
	}
}
