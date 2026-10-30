package mcpserver

import (
	"context"
	"testing"

	"movieSpider/internal/config"
)

func init() {
	config.InitConfig("../../config.local.yaml")
}

func TestMovieService_PlayableTodayMovieTV(t *testing.T) {
	cfg := &MovieServiceConfig{
		AccountID:   config.Config().TMDB.AccountID,
		BearerToken: config.Config().TMDB.BearerToken,
	}
	svc, err := NewMovieService(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	tv, err := svc.PlayableTodayMovieTV(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log(tv)
}
