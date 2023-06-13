package ipProxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"net/http"
)

type IpProxyPoolDataIP struct {
	ProxyHost     string `json:"proxyHost"`
	ProxyPort     int    `json:"proxyPort"`
	ProxyType     string `json:"proxyType"`
	ProxyLocation string `json:"proxyLocation"`
	ProxySpeed    int    `json:"proxySpeed"`
	ProxySource   string `json:"proxySource"`
	UpdateTime    string `json:"updateTime"`
}

// FetchProxy
//  @Description: 获取代理
//  @param typ
//  @return string
//
func FetchProxy(typ string) *IpProxyPoolDataIP {
	if config.ProxyPool == "" {
		log.Warn("FetchProxy: Global.ProxyPool没有配置.")
		return nil
	}

	urlStr := fmt.Sprintf("%s/%s", config.ProxyPool, typ)
	resp, err := http.DefaultClient.Get(urlStr)
	if err != nil {
		log.Errorf("Feed.ProxyPool %s,err: %s", config.ProxyPool, err.Error())
		return nil
	}
	defer resp.Body.Close()

	data, err := parseIpProxyPoolData(resp.Body)
	if err != nil {
		log.Warnf("FetchProxy: %s.", err.Error())
		return nil
	}

	return data
}

func parseIpProxyPoolData(body io.Reader) (*IpProxyPoolDataIP, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	var data IpProxyPoolDataIP
	err := json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

type proxyTypeCount struct {
	Http  int64 `json:"http"`
	Https int64 `json:"https"`
	Tcp   int64 `json:"tcp"`
	Other int64 `json:"other"`
}

func FetchProxyTypeCount() (c *proxyTypeCount) {
	c = &proxyTypeCount{}
	if config.ProxyPool == "" {
		log.Warn("FetchProxy: Global.ProxyPool没有配置.")
		return nil
	}

	urlStr := fmt.Sprintf("%s/count", config.ProxyPool)
	resp, err := http.DefaultClient.Get(urlStr)
	if err != nil {
		log.Errorf("Feed.ProxyPool %s,err: %s", config.ProxyPool, err.Error())
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
