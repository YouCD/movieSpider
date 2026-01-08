package feedspider

import (
	"context"
	"errors"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestNewGlodls(t *testing.T) {
	feeder := NewGlodls()
	videos, err := feeder.Crawler(context.Background())
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			if errors.Is(err, model.ErrFeedVideoExclude) {
				continue
			}
			log.WithCtx(context.Background()).Errorf("err: %s    %#v", err, video)
			continue
		}
		log.WithCtx(context.Background()).Infof("%#v", filterVideo)
	}
}
