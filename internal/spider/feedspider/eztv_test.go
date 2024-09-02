package feedspider

import (
	"movieSpider/internal/config"
	"testing"
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
		t.Log(video)
	}
}
