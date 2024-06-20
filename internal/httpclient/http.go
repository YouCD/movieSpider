package httpclient

import (
	"movieSpider/internal/config"
	"net/http"
	"net/url"
	"sync"
	"time"
)

//nolint:gochecknoglobals
var (
	httpClient *http.Client
	once       = &sync.Once{}
)

// NewHTTPClient
//
//	@Description: 初始化http client
//	@return *http.Client
func NewHTTPClient() *http.Client {
	once.Do(func() {
		httpClient = &http.Client{Timeout: time.Second * 60}
		if config.Config.Global.Proxy.URL != "" {
			proxyURL, _ := url.Parse(config.Config.Global.Proxy.URL)
			proxy := http.ProxyURL(proxyURL)
			transport := &http.Transport{Proxy: proxy}
			httpClient = &http.Client{Transport: transport, Timeout: time.Minute}
		}
	})

	return httpClient
}
