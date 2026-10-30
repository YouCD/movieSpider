package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig("../../config.local.yaml")
	marshal, err := json.Marshal(Config())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(marshal))
}
