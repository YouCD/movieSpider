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
//	@Description: 初始化普通http client
//	@return *http.Client
func NewHTTPClient() *http.Client {
	once.Do(func() {
		//nolint:exhaustruct
		httpClient = &http.Client{Timeout: time.Second * 60}
	})

	return httpClient
}

// NewHTTPProxyClient
//
//	@Description: 创建代理http client
//	@return *http.Client
func NewHTTPProxyClient() *http.Client {
	once.Do(func() {
		if config.Config.TG.Proxy.URL != "" {
			proxyURL, _ := url.Parse(config.Config.TG.Proxy.URL)
			proxy := http.ProxyURL(proxyURL)
			//nolint:exhaustruct
			transport := &http.Transport{Proxy: proxy}
			//nolint:exhaustruct
			httpClient = &http.Client{Transport: transport}
		} else {
			//nolint:exhaustruct
			httpClient = &http.Client{Timeout: time.Second * 30}
		}
	})

	return httpClient
}
