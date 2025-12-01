package feedspider

import (
	"errors"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestNewGlodls(t *testing.T) {
	feeder := NewGlodls()
	videos, err := feeder.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		_, err := model.FilterVideo(video)
		if err != nil {
			if errors.Is(err, model.ErrFeedVideoExclude) {
				continue
			}
			log.WithCtx(context.Background()).Errorf("err: %s    %#v", err, video)
			continue
		}
		//log.Infof("%#v", filterVideo)
	}

}
