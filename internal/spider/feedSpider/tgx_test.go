package feedSpider

import (
	"movieSpider/internal/config"
	"testing"
)

func Test_tgx_Run(t1 *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	t := tgx{
		scheduling: "tt.fields.scheduling",
		url:        urlTgx,
		web:        "tgx",
	}
	t.Run()
}
