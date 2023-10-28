package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"testing"
)

func Test_tpbpirateproxy_Crawler(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/core/bin/core/config.yaml")
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	g := &tpbpirateproxy{
		scheduling: "tpbpirateproxy",
		web:        "tpbpirateproxy",
	}
	gotVideos, err := g.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range gotVideos {
		fmt.Println(video)
	}
}
