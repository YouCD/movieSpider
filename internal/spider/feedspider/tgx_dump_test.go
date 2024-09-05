package feedspider

import (
	"errors"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"strings"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	//model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestNewTgxDump(t *testing.T) {
	for _, tgx := range config.Config.Feed.TGX {
		if strings.HasSuffix(tgx.Url, "tgx24hdump.txt.gz") {
			videos, err := NewTgxDump(tgx.Scheduling, tgx.Url, true).Crawler()
			if err != nil {
				t.Error(err)
				continue
			}
			for _, video := range videos {
				vv, err := model.FilterVideo(video)
				if err != nil {
					if errors.Is(err, model.ErrFeedVideoExclude) {
						continue
					}
					log.Errorf("err: %s    %#v", err, video)
					continue
				}
				log.Info(vv)
			}
		}
	}
}
