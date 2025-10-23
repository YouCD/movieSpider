package job

import (
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}

func Test_reportAria2TaskStatistics(t *testing.T) {
	reportAria2TaskStatistics()
}
