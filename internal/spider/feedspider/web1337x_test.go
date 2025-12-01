package feedspider

import (
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestWeb1337x_Crawler(t *testing.T) {
	//web1337x := NewWeb1337x("*/3 * * * *", types.VideoTypeTV, "https://1337x.to/popular-tv")
	web1337x := NewWeb1337x("*/3 * * * *", types.VideoTypeMovie, "https://1337x.to/popular-movies", true)
	gotVideos, err := web1337x.Crawler()
	if err != nil {
		t.Errorf("Crawler() error = %v", err)
		return
	}
	for _, video := range gotVideos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			log.WithCtx(context.Background()).Errorf("err: %s    %#v", err, video)
			continue
		}
		log.WithCtx(context.Background()).Infof("%#v", filterVideo)
	}
}
