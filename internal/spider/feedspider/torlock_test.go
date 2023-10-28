package feedspider

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
)

func Test_torlock_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/core/bin/core/config.yaml")
	model.NewMovieDB()
	var facTORLOCK FactoryTORLOCK

	videos, err := facTORLOCK.CreateFeeder("*/1 * * * *", types.VideoTypeTV).Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		fmt.Println(video)
	}

	//facTORLOCK.CreateFeeder("*/1 * * * *", types.VideoTypeTV).Run()
	//select {}

	//videos, err := facTORLOCK.CreateFeeder("*/1 * * * *",  types.VideoTypeMovie).Search()
	//if err != nil {
	//	t.Error(err)
	//}
	//for _, video := range videos {
	//	fmt.Println(video)
	//}
}
