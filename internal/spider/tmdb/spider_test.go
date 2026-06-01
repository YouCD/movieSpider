package tmdb

import (
	"context"
	"testing"

	"movieSpider/internal/config"
)

func TestNewTMDBSpider(t *testing.T) {
	spider, err := NewTMDBSpider(config.Config.TMDB.AccountID, config.Config.TMDB.BearerToken)
	if err != nil {
		t.Errorf("创建TMDB爬虫失败: %s", err)
		return
	}
	//spider.Run(context.Background())
	//spider.checkWatchProviders(context.Background())
	spider.crawlTV(context.Background())
}
