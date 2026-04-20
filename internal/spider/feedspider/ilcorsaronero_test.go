package feedspider

import (
	"context"
	"testing"

	"movieSpider/internal/config"

	"github.com/youcd/toolkit/log"
)

func TestIlcorsaronero_Crawler(t *testing.T) {
	for _, item := range config.Config.Feed.Ilcorsaronero {
		log.WithCtx(context.Background()).Info(item.Url)
		u := NewIlcorsaronero(item.Scheduling, item.ResourceType, item.Url, item.UseIPProxy)
		got, err := u.Crawler(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		for _, base := range got {
			log.WithCtx(context.Background()).Infof("%#v", base)
		}
	}
}
