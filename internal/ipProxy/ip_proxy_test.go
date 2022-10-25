package ipProxy

import (
	"fmt"
	"movieSpider/internal/config"
	"testing"
)

func TestFetchProxy(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	proxy := FetchProxy("https")
	fmt.Println(proxy)
}
