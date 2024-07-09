package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestNewTgxWeb(t *testing.T) {
	web := NewTgxWeb("*/3 * * * *", "https://tgx.rs/torrents.php?c3=1&c42=1&c41=1&c11=1&search=&lang=0&nox=2#resultss")
	//webHost := NewTgxWeb("*/3 * * * *", "https://www.btbtt12.com/forum-index-fid-951.htm")
	videos, err := web.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		fmt.Println(video)
	}
}
