package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	fmt.Println(DouBanList)
}
