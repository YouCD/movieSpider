package httpClient

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"movieSpider/internal/config"
	"movieSpider/internal/ipProxy"
	"movieSpider/internal/log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
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
			httpClient = &http.Client{Timeout: time.Second * 30}
		}
	})

	return httpClient
}

func NewProxyHttpClient(proxyType string) *http.Client {

	proxyObj := ipProxy.FetchProxy(proxyType)
	if proxyObj == nil {
		return nil
	}

	proxy := strings.ToLower(fmt.Sprintf("%s://%s:%d", proxyObj.ProxyType, proxyObj.ProxyHost, proxyObj.ProxyPort))
	log.Infof("proxy change to  %s", proxy)
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		log.Error(err)
		return nil
	}
	if proxyUrl == nil {
		return nil
	}

	switch {
	case strings.Contains(proxy, "tcp"):
		dialer := &net.Dialer{
			// 限制创建一个TCP连接使用的时间（如果需要一个新的链接）
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		conn, err := dialer.Dial("tcp", proxyObj.ProxyHost+":"+strconv.Itoa(proxyObj.ProxyPort))
		if err != nil {
			log.Error(err)
		}
		// 设置网络传输
		netTransport := &http.Transport{
			DialContext:           dialer.DialContext,
			Proxy:                 nil,
			DisableKeepAlives:     true,
			MaxConnsPerHost:       20,
			MaxIdleConns:          20,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       20 * time.Second,
			ResponseHeaderTimeout: time.Second * time.Duration(30),
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
		_ = http2.ConfigureTransport(netTransport)

		// 创建连接客户端
		defer conn.Close()
		return &http.Client{Transport: netTransport}
	case strings.Contains(proxy, "https"), strings.Contains(proxy, "http"):
		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		return &http.Client{Transport: transport, Timeout: time.Second * 30}
	}
	return nil
}
