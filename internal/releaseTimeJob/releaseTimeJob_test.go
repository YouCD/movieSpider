package releaseTimeJob

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB()
}
func TestDatePublished_Run(t *testing.T) {
	r := &ReleaseTimeJob{
		scheduling: "tt.fields.scheduling",
	}
	r.Run()
	select {}
}
