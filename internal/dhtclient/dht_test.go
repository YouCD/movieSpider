package dhtclient

import (
	"context"
	"testing"

	"movieSpider/internal/config"
	"movieSpider/internal/model"
)

func init() {
	config.InitConfig("../../config.local.yaml")
	// model.NewMovieDB().SaveFeedVideoFromChan()
}

func TestBoot(t *testing.T) {
	ctx := context.Background()
	model.NewMovieDB().SaveFeedVideoFromChan(ctx)
	Boot(2, ctx)
}
