package feedspider

import (
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestNewBtbt(t *testing.T) {
	feeder := NewBtbt()
	videos, err := feeder.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			log.Infof("%#v", video)
			continue
		}
		log.Warnf("%#v", filterVideo)
	}

}
