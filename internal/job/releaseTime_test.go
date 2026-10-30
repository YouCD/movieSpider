package job

import (
	"context"
	"testing"

	"movieSpider/internal/config"
	"movieSpider/internal/model"
)

func init() {
	config.InitConfig("../../config.local.yaml")
	model.NewMovieDB()
}

func TestDatePublished_Run(t *testing.T) {
	r := &ReleaseTimeJob{
		scheduling: "tt.fields.scheduling",
	}
	r.Run(context.Background())
	select {}
}
