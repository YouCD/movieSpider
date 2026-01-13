package feedspider

import (
	"context"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestYts_Crawler(t *testing.T) {
	//web1337x := NewWeb1337x("*/3 * * * *", types.VideoTypeTV, "https://www.1337xx.to/popular-tv")
	//log.SetLogLevel("DEBUG")
	//web1337x := NewWeb1337x("*/3 * * * *", types.VideoTypeMovie, "https://www.1337xx.to/popular-movies", true)
	yts := NewYts()
	gotVideos, err := yts.Crawler(context.Background())
	if err != nil {
		t.Error(err)
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
