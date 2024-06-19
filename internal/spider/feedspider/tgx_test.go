package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"testing"
)

func Test_tgx_Run(t1 *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	t := Tgx{
		BaseFeeder{
			scheduling: "tt.fields.scheduling",
			url:        fmt.Sprintf("%s/%s", urlBaseTgx, urlRssURITgx),
			web:        "tgx",
		},
	}
	videos, err := t.Crawler()
	if err != nil {
		t1.Errorf("tgx.Run() error = %v", err)
	}
	for _, video := range videos {
		fmt.Printf("%#v\n", video)
	}
}
