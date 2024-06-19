package feedspider

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB().SaveFeedVideoFromChan()
}
