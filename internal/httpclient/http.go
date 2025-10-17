package httpclient

import (
	"crypto/tls"
	"fmt"
	"movieSpider/internal/config"
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

const (
	httpsHTTPBin = "https://www.google.com"
	httpHTTPBin  = "http://httpbin.org/get"
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
		IdleConnTimeout:       time.Second * time.Duration(config.Config.Global.Timeout),
		ResponseHeaderTimeout: time.Second * time.Duration(config.Config.Global.Timeout),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	log.Debugf("use proxy: %s", proxy)
	return &http.Client{Transport: netTransport, Timeout: time.Second * time.Duration(config.Config.Global.Timeout)}, proxy
}

func NewIPProxyPoolHTTPClientDel(exampleURL string) (*http.Client, string) {
RETRY:
	var testURL string
	var proxyObj *ipproxy.PoolDataIP
	switch {
	case strings.HasPrefix(exampleURL, "https"):
		proxyObj = ipproxy.FetchProxy("https")
		testURL = httpsHTTPBin
	case strings.HasPrefix(exampleURL, "http"):
		proxyObj = ipproxy.FetchProxy("http")
		testURL = httpHTTPBin
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
		ResponseHeaderTimeout: time.Second * time.Duration(60),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	log.Debugf("use proxy: %s", proxy)

	httpClient := &http.Client{Transport: netTransport, Timeout: time.Second * 60}
	//nolint:noctx
	res, err := httpClient.Get(testURL)
	if err != nil {
		log.Debugf("proxy: %s, error: %s", proxy, err)
		ipproxy.DelProxy(proxyObj.ProxyHost)
		goto RETRY
	}

	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()
	if res.StatusCode == http.StatusOK {
		return httpClient, proxy
	}
	// 回调删除代理
	ipproxy.DelProxy(proxyObj.ProxyHost)
	goto RETRY
}
