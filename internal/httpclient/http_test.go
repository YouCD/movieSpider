package httpclient

import (
	"movieSpider/internal/config"
	"testing"

	"github.com/youcd/toolkit/log"
)

func init() {
	log.Init(true)
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func TestNewIpProxyPoolHTTPClient(t *testing.T) {
	NewProxyHTTPClient("https://thepiratebay.org/search.php?q=top100:200")
}
