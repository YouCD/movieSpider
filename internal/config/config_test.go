package config

import (
	"fmt"
	"movieSpider/internal/tools"
	"testing"
)

func init() {
	InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}

func TestExcludeVideo(t *testing.T) {
	video := tools.ExcludeVideo("Foundation.S02E07.2160p.Dolby.Vision.Multi.Sub.DDP5.1.Atmos.DV.x265.MKV-BEN.THE.MEN", ExcludeWords)
	fmt.Println(video)
}
