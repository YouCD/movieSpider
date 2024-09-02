package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	marshal, err := json.Marshal(Config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(marshal))
}
