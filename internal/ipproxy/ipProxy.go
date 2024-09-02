package ipproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"movieSpider/internal/config"
	"net/http"

	"github.com/pkg/errors"
	"github.com/youcd/toolkit/log"
)

var (
	ErrProxyIsEmpty = errors.New("proxy is empty")
)

type PoolDataIP struct {
	ProxyHost     string `json:"proxyHost"`
	ProxyPort     int    `json:"proxyPort"`
	ProxyType     string `json:"proxyType"`
	ProxyLocation string `json:"proxyLocation"`
	ProxySpeed    int    `json:"proxySpeed"`
	ProxySource   string `json:"proxySource"`
	UpdateTime    string `json:"updateTime"`
}

// FetchProxy
//
//	@Description: 获取代理
//	@param typ
//	@return string
func FetchProxy(typ string) *PoolDataIP {
	if config.Config.Global.IPProxyPool == "" {
		log.Warn("FetchProxy: Global.IpProxyPool没有配置.")
		return nil
	}
	urlStr := fmt.Sprintf("%s/%s", config.Config.Global.IPProxyPool, typ)
	//nolint:noctx
	resp, err := http.DefaultClient.Get(urlStr)
	if err != nil {
		log.Errorf("Feed.ProxyPool %s,err: %s", config.Config.Global.IPProxyPool, err.Error())
		return nil
	}
	defer resp.Body.Close()

	data, err := parseIPProxyPoolData(resp.Body)
	if err != nil {
		log.Warnf("FetchProxy: %s.", err.Error())
		return nil
	}

	return data
}

func parseIPProxyPoolData(body io.Reader) (*PoolDataIP, error) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(body)
	var data PoolDataIP
	err := json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return nil, errors.WithMessage(err, "parseIPProxyPoolData")
	}
	if data.ProxyHost == "" {
		return nil, ErrProxyIsEmpty
	}
	return &data, nil
}

type ProxyTypeCount struct {
	HTTP  int64 `json:"http"`
	HTTPS int64 `json:"https"`
	TCP   int64 `json:"tcp"`
	Other int64 `json:"other"`
}

func FetchProxyTypeCount() (c *ProxyTypeCount) {
	//nolint:exhaustruct
	c = &ProxyTypeCount{}
	if config.Config.Global.IPProxyPool == "" {
		log.Warn("FetchProxy: Global.ProxyPool没有配置.")
		return nil
	}

	urlStr := config.Config.Global.IPProxyPool + "/count"
	//nolint:noctx
	resp, err := http.DefaultClient.Get(urlStr)
	if err != nil {
		log.Errorf("Feed.ProxyPool %s,err: %s", config.Config.Global.IPProxyPool, err.Error())
		return nil
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("FetchProxyTypeCount: %s.", err.Error())
		return nil
	}
	err = json.Unmarshal(all, c)
	if err != nil {
		log.Warnf("FetchProxyTypeCount: %s.", err.Error())
		return nil
	}

	return
}
