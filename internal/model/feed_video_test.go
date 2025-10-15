package model

import (
	"context"
	"fmt"
	"movieSpider/internal/nameParser"
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
	log.Info(tvs)
	query := dbA.Table("feed_video").Model(&types.FeedVideo{}).Where("type = ?", "tv").Limit(1)
	log.Info("准备执行查询")
	err := query.Find(&tvs).Error
	log.Infof("查询结果: %+v, 错误: %v", tvs, err)

	for _, tv := range tvs {
		fmt.Println(tv)
	}

	typeStr, newName, year, resolution, err := nameParser.NameParserModelHandler(context.Background(), "The Morning Show S04E02 The Revolution Will Be Televised XviD-AFG")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(typeStr, newName, year, resolution)
}
