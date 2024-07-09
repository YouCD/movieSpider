package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestThePirateBay_Crawler(t1 *testing.T) {
	thePirateBay := NewThePirateBay("*/3 * * * *", "https://thepiratebay.org/search.php?q=top100:200")
	gotVideos, err := thePirateBay.Crawler()
	if err != nil {
		t1.Errorf("Crawler() error = %v", err)
		return
	}
	fmt.Println(gotVideos)
	for _, video := range gotVideos {
		log.Errorf("%#v", video)
	}
}
