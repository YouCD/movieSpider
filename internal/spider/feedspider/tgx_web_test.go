package feedspider

import (
	"context"
	"errors"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestTgxWeb_Crawler(t *testing.T) {
	//feeder := NewTgxWeb("*/3 * * * *", "https://torrentgalaxy.info/get-posts/category:TV:order:-se/", types.VideoTypeTV)
	feeder := NewTgxWeb("*/3 * * * *", types.VideoTypeTV, "https://torrentgalaxy.info/get-posts/category:Movies:order:-se/", false)
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
			// log.Errorf("err: %s    %#v", err, video)
			continue
		}
		log.WithCtx(context.Background()).Infof("%#v", filterVideo)
	}
}
