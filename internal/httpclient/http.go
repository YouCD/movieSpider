package httpclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"movieSpider/internal/config"

	"github.com/youcd/toolkit/log"
)

//nolint:gochecknoglobals
var (
	HTTPClient = &http.Client{}
)

// NewProxyHTTPClient 创建带代理的 HTTP 客户端
// 代理配置无效时降级为普通客户端, 保证永远不返回 nil
func NewProxyHTTPClient(ctx context.Context) *http.Client {
	if config.Config().Global.ProxyURL == "" {
		return HTTPClient
	}
	proxyURL, err := url.Parse(config.Config().Global.ProxyURL)
	if err != nil {
		log.WithCtx(ctx).Errorf("解析代理地址失败, 使用直连: %v", err)
		return HTTPClient
	}
	// 设置网络传输
	//nolint:exhaustruct,gosec
	netTransport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),

		// 连接池
		MaxConnsPerHost:     20,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     time.Second * 90,

		ResponseHeaderTimeout: time.Second * time.Duration(config.Config().Global.Timeout),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	log.WithCtx(ctx).Debugf("use proxy: %s", config.Config().Global.ProxyURL)
	return &http.Client{Transport: netTransport, Timeout: time.Second * time.Duration(config.Config().Global.Timeout)}
}
