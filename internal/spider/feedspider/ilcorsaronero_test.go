package feedspider

import (
	"movieSpider/internal/config"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestIlcorsaronero_Crawler(t *testing.T) {
	for _, item := range config.Config.Feed.Ilcorsaronero {
		u := NewIlcorsaronero(item.Scheduling, item.ResourceType, item.Url, item.UseIPProxy)
		got, err := u.Crawler()
		if err != nil {
			t.Error(err)
			return
		}
		for _, base := range got {
			log.Infof("%#v", base)
		}
	}
}
