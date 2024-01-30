package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"testing"
)

func Test_tpbpirateproxy_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	g := &Tpbpirateproxy{
		BaseFeeder{
			web: "tpbpirateproxy",
			url: fmt.Sprintf("%s/%s", urlBaseTpbpirateProxy, urlRssURITpbpirateProxy),
		},
	}
	gotVideos, err := g.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range gotVideos {
		fmt.Println(video)
	}
}
