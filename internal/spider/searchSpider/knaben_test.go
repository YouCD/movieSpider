package searchSpider

import (
	"encoding/json"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
)

func TestEztv_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/core/bin/core/config.yaml")
	model.NewMovieDB()

	f := NewFeedKnaben("House Of The Dragon", types.ResolutionOther)
	videos, err := f.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, v := range videos {
		bytes, _ := json.Marshal(v)

		fmt.Println(string(bytes))
		//fmt.Println(v.Name)

	}
}
