package job

import (
	"movieSpider/internal/config"
	"testing"
)

func Test_reportIpProxyStatistics(t *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	reportIPProxyStatistics()
}
