package httpclient

import (
	"context"
	"testing"

	"movieSpider/internal/config"

	"github.com/youcd/toolkit/log"
)

func init() {
	log.Init(&log.Config{Stdout: true})
	config.InitConfig("../../config.local.yaml")
}

func TestNewIpProxyPoolHTTPClient(t *testing.T) {
	NewProxyHTTPClient(context.Background())
}
