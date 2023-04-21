package httpClient

import (
	"movieSpider/internal/config"
	"net/http"
	"net/url"
	"sync"
)

var (
	httpClient *http.Client
	once       = &sync.Once{}
)

//
// NewHttpClient
//  @Description: 初始化http client
//  @return *http.Client
//
func NewHttpClient() *http.Client {
	once.Do(func() {
		if config.TG.Proxy.Enable {
			if config.TG.Proxy.Url != "" {
				proxyUrl, _ := url.Parse(config.TG.Proxy.Url)
				proxy := http.ProxyURL(proxyUrl)
				transport := &http.Transport{Proxy: proxy}
				httpClient = &http.Client{Transport: transport}
			}
		} else {
			httpClient = &http.Client{}
		}
	})

	return httpClient
}
