package feedspider

import (
	"errors"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestEztv_Crawler(t *testing.T) {
	feeder := NewFeedKnaben()
	videos, err := feeder.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			if errors.Is(err, model.ErrFeedVideoExclude) {
				continue
			}
			//log.Errorf("err: %s    %#v", err, video)
			continue
		}
		log.WithCtx(context.Background()).Error("%#v", filterVideo)
	}

}
