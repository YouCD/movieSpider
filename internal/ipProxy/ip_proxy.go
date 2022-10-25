package ipProxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"net/http"
)

type proxyData struct {
	Anonymous  string `json:"anonymous"`
	CheckCount int    `json:"check_count"`
	FailCount  int    `json:"fail_count"`
	Https      bool   `json:"https"`
	LastStatus bool   `json:"last_status"`
	LastTime   string `json:"last_time"`
	Proxy      string `json:"proxy"`
	Region     string `json:"region"`
	Source     string `json:"source"`
}

func FetchProxy(typ string) string {
	if config.ProxyPool != "" {
		var urlStr string
		if typ == "" {
			urlStr = fmt.Sprintf("%s/get", config.ProxyPool)
		} else {
			urlStr = fmt.Sprintf("%s/get/?type=https", config.ProxyPool)
		}
		resp, err := http.Get(urlStr)
		if err != nil {
			log.Errorf("Feed.ProxyPool %s,err: %s", config.ProxyPool, err.Error())
			return ""
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		defer resp.Body.Close()

		var data proxyData
		err = json.Unmarshal(buf.Bytes(), &data)
		if err != nil {
			log.Warnf("FetchProxy: %s.", err.Error())
			return ""
		}
		if data.Proxy != "" {
			return fmt.Sprintf("http://%s", data.Proxy)
		}
		return ""
	}
	log.Warn("FetchProxy: Global.ProxyPool没有配置.")
	return ""
}
