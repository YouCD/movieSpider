package httpclient

import (
	"crypto/tls"
	"fmt"
	"movieSpider/internal/ipproxy"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/youcd/toolkit/log"
)

//nolint:gochecknoglobals
var (
	HTTPClient = &http.Client{}
)

func NewProxyHTTPClient(proxy string) *http.Client {
	proxyURL, _ := url.Parse(proxy)
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	return &http.Client{Transport: transport, Timeout: time.Minute}
}

func NewIPProxyPoolHTTPClient(exampleURL string) (*http.Client, string) {
	var proxyObj *ipproxy.PoolDataIP
	switch {
	case strings.HasPrefix(exampleURL, "https"):
		proxyObj = ipproxy.FetchProxy("https")
	case strings.HasPrefix(exampleURL, "http"):
		proxyObj = ipproxy.FetchProxy("http")
	}
	if proxyObj == nil {
		return nil, ""
	}

	proxy := strings.ToLower(fmt.Sprintf("%s://%s:%d", proxyObj.ProxyType, proxyObj.ProxyHost, proxyObj.ProxyPort))
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Error(err)
		return nil, ""
	}
	// 设置网络传输
	//nolint:exhaustruct,gosec
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxyURL),
		DisableKeepAlives:     true,
		MaxConnsPerHost:       20,
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       20 * time.Second,
		ResponseHeaderTimeout: time.Second * time.Duration(30),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	log.Debugf("use proxy: %s", proxy)
	return &http.Client{Transport: netTransport, Timeout: time.Second * 30}, proxy
}
