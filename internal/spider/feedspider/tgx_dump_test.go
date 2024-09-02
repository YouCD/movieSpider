package feedspider

import (
	"movieSpider/internal/config"
	"strings"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	//model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestNewTgxDump(t *testing.T) {
	for _, tgx := range config.Config.Feed.TGX {
		if strings.HasSuffix(tgx.Url, "tgx24hdump.txt.gz") {
			NewTgxDump(tgx.Scheduling, tgx.Url, true).Crawler()
		}
	}
}
