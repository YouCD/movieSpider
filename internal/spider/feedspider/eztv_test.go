package feedspider

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}

func TestNewEztv(t *testing.T) {
	eztv := NewEztv()
	videos, err := eztv.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			log.WithCtx(context.Background()).Errorf("err:%s,%#v", err, video)
			continue
		}
		log.WithCtx(context.Background()).Infof("%#v", filterVideo)
	}

}
