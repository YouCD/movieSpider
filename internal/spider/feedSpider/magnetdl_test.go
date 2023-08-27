package feedSpider

import (
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"testing"
)

func Test_magnetdl_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	tv := &magnetdl{
		typ:        types.VideoTypeTV,
		web:        "magnetdl",
		scheduling: "*/1 * * * *",
	}
	tv.Run(bus.FeedVideoChan)

	m := &magnetdl{
		typ:        types.VideoTypeMovie,
		web:        "magnetdl",
		scheduling: "*/2 * * * *",
	}
	m.Run(bus.FeedVideoChan)
	select {}
}
