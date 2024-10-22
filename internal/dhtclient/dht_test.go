package dhtclient

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	//model.NewMovieDB().SaveFeedVideoFromChan()
}
func TestBoot(t *testing.T) {
	model.NewMovieDB().SaveFeedVideoFromChan()
	Boot()
}
