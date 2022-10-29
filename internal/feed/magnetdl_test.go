package feed

import (
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"testing"
)

func Test_magnetdl_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	tv := &magnetdl{
		typ:        types.ResourceTV,
		web:        "magnetdl",
		scheduling: "*/1 * * * *",
	}
	tv.Run()

	m := &magnetdl{
		typ:        types.ResourceMovie,
		web:        "magnetdl",
		scheduling: "*/2 * * * *",
	}
	m.Run()
	select {}
}
