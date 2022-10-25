package feed

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
)

func Test_torlock_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()
	var facTORLOCK FactoryTORLOCK

	videos, err := facTORLOCK.CreateFeeder("*/1 * * * *", types.ResourceTV).Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		fmt.Println(video)
	}

	//facTORLOCK.CreateFeeder("*/1 * * * *", types.ResourceTV).Run()
	//select {}

	//videos, err := facTORLOCK.CreateFeeder("*/1 * * * *", types.ResourceMovie).Crawler()
	//if err != nil {
	//	t.Error(err)
	//}
	//for _, video := range videos {
	//	fmt.Println(video)
	//}
}
