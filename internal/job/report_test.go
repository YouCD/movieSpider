package job

import (
	"testing"

	"movieSpider/internal/config"
)

func init() {
	config.InitConfig("../../config.local.yaml")
}

func Test_reportAria2TaskStatistics(t *testing.T) {
	reportAria2TaskStatistics()
}
