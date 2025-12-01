package model

import (
	"context"
	"fmt"
	"movieSpider/internal/nameparser"
	"movieSpider/internal/types"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestMovieDB_GetFeedVideoMovieByNames(t *testing.T) {
	movieDB := NewMovieDB()

	got, err := movieDB.findMovie([]string{"The.Fogsa"}, `name = ? and magnet!="" and type="movie"`)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(got)
}
func TestNameParserModelHandler(t *testing.T) {
	dbA := NewMovieDB().GetDB()
	var tvs []*types.FeedVideo
	log.WithCtx(context.Background()).Info(tvs)
	query := dbA.Table("feed_video").Model(&types.FeedVideo{}).Where("type = ?", "tv").Limit(1)
	log.WithCtx(context.Background()).Info("准备执行查询")
	err := query.Find(&tvs).Error
	log.WithCtx(context.Background()).Infof("查询结果: %+v, 错误: %v", tvs, err)

	for _, tv := range tvs {
		fmt.Println(tv)
	}

	result, err := nameparser.ModelHandler(context.Background(), "The Morning Show S04E02 The Revolution Will Be Televised XviD-AFG")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}
