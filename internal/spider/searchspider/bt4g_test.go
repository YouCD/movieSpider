package searchspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
)

func Test_bt4g_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/core/bin/core/config.yaml")
	model.NewMovieDB()
	//https: //BT4g.org/search/%E9%BE%99%E4%B9%8B%E5%AE%B6%E6%97%8F?page=rss
	//https://bt4g.org/search/%E9%BE%99%E4%B9%8B%E5%AE%B6%E6%97%8F?page=rss
	b := NewFeedBt4g("杀手疾风号", types.ResolutionOther)

	gotVideos, err := b.Search()
	if err != nil {
		t.Error(err)
	}
	for _, v := range gotVideos {
		fmt.Println(v.Magnet)
	}
}
