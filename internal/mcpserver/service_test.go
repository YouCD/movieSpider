package mcpserver

import (
	"context"
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func TestMovieService_PlayableTodayMovieTV(t *testing.T) {
	tv, err := NewMovieService().PlayableTodayMovieTV(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log(tv)
}
