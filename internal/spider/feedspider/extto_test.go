package feedspider

import (
	"errors"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestNewExtto(t *testing.T) {
	for _, r := range config.Config.Feed.Extto {
		feeder := NewExtto(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy, r.UseCloudflareBypass)
		videos, err := feeder.Crawler()
		if err != nil {
			t.Error(err)
		}
		for _, video := range videos {
			filterVideo, err := model.FilterVideo(video)
			if err != nil {
				if errors.Is(err, model.ErrFeedVideoExclude) {
					log.Warnf("err: %s    %#v", err, video)
					continue
				}
				continue
			}
			log.Infof("%#v", filterVideo)
		}
		//break
	}
}
