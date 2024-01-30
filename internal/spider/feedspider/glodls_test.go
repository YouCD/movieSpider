package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestGlodls_Crawler(t *testing.T) {
	f := &glodls{
		url: "http://glodls.to/rss.php?cat=1,41",
	}
	log.Info("GLODLS: is working...")
	videos, err := f.Crawler()
	if err != nil {
		log.Error(err)
		return
	}
	for _, video := range videos {
		fmt.Printf("%#v\n", video)
	}

}
