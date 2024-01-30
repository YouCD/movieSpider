package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func Test_eztv_Crawler(t *testing.T) {

	var e = Eztv{
		BaseFeeder{
			web: "eztv",
			url: fmt.Sprintf("%s/%s", urlBaseEztv, urlRssURIEztv),
		},
	}
	videos, err := e.Crawler()

	if err != nil {
		t.Error(err)
	}
	model.ProxySaveVideo2DB(videos...)
	select {}
}
