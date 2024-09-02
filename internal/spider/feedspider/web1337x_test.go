package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestWeb1337x_Crawler(t *testing.T) {
	//web1337x := NewWeb1337x("*/3 * * * *", types.VideoTypeTV, "https://1337x.to/popular-tv")
	web1337x := NewWeb1337x("*/3 * * * *", types.VideoTypeMovie, "https://1337x.to/popular-movies", true)
	gotVideos, err := web1337x.Crawler()
	if err != nil {
		t.Errorf("Crawler() error = %v", err)
		return
	}
	fmt.Println(gotVideos)
	for _, video := range gotVideos {
		log.Errorf("%#v", video)
	}
}
