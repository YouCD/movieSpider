package feedspider

import (
	"movieSpider/internal/config"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestUindex_Crawler(t *testing.T) {
	for _, uindex := range config.Config.Feed.Uindex {

		u := NewUindex(uindex.Scheduling, uindex.Url, uindex.UseIPProxy)
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
