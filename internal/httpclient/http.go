package httpclient

import (
	"context"
	"crypto/tls"
	"movieSpider/internal/config"
	"net/http"
	"net/url"
	"time"

	"github.com/youcd/toolkit/log"
)

//nolint:gochecknoglobals
var (
	HTTPClient = &http.Client{}
)

func NewProxyHTTPClient(ctx context.Context) *http.Client {
	proxyURL, err := url.Parse(config.Config.Global.ProxyURL)
	if err != nil {
		log.WithCtx(ctx).Error(err)
		return nil
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

		ResponseHeaderTimeout: time.Second * time.Duration(config.Config.Global.Timeout),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	log.WithCtx(ctx).Debugf("use proxy: %s", config.Config.Global.ProxyURL)
	return &http.Client{Transport: netTransport, Timeout: time.Second * time.Duration(config.Config.Global.Timeout)}
}
