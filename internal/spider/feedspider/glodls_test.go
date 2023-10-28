package feedspider

import (
	"movieSpider/internal/config"
	httpClient2 "movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"net/http"
	"testing"
	"time"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestGlodls_Crawler(t *testing.T) {
	f := &glodls{
		url:        "http://glodls.to/rss.php?cat=1,41",
		httpClient: &http.Client{Timeout: time.Second * 3},
	}
	for {
	Start:
		log.Info("GLODLS: is working...")
		videos, err := f.Crawler()
		if err != nil {
			log.Error(err)
		}
		if len(videos) == 0 || videos == nil {
			log.Info("GLODLS: 切换代理")
			f.httpClient = httpClient2.NewProxyHTTPClient("http")
			f.httpClient.Timeout = 3 * time.Second
			goto Start
		}

	}

}
