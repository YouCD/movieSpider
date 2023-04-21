package feed

import (
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func Test_eztv_Crawler(t *testing.T) {

	var e = eztv{
		scheduling: "",
		url:        urlEztv,
		web:        "eztv",
	}
	videos, err := e.Crawler()

	if err != nil {
		t.Error(err)
	}
	proxySaveVideo2DB(videos...)
	select {}
}
