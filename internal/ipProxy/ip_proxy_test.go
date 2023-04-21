package ipProxy

import (
	"fmt"
	"movieSpider/internal/config"
	"testing"
)

func TestFetchProxy(t *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	proxy := FetchProxy("https")
	fmt.Println(proxy)
}
