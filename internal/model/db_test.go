package model

import (
	"context"
	"fmt"
	"testing"

	"movieSpider/internal/config"
	"movieSpider/internal/types"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}

func TestNewMovieDB(t *testing.T) {
	NewMovieDB()
}

func Test_movieDB_AddDownloadHistory(t *testing.T) {
	err := NewMovieDB().AddDownloadHistory(&types.DownloadHistory{
		ID:          0,
		Name:        "Raven.of.the.Inner.Palace",
		Type:        "tv",
		TorrentName: "Raven.of.the.Inner.Palace.S01E03.1080p.WEB.H264-SENPAI",
		Timestamp:   0,
		Resolution:  0,
		Season:      "01",
		Episode:     "03",
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_movieDB_CountFeedVideo(t *testing.T) {
	counts, err := NewMovieDB().CountFeedVideo()
	if err != nil {
		t.Error(err)
	}
	for _, count := range counts {
		fmt.Println(count)
	}
}

func Test_movieDB_CreatTMDBVideo(t *testing.T) {
	err := NewMovieDB().CreatTMDBVideo(context.Background(), &types.TMDBVideo{
		ID:        99119,
		Names:     `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
		ImdbID:    "tt444444444",
		RowData:   "444444444",
		Timestamp: 0,
		Type:      "444444444",
		Playable:  "三21问问是岁",
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_movieDB_CreatFeedVideo(t *testing.T) {
	err := NewMovieDB().CreatFeedVideo(&types.FeedVideo{
		ID:        888888,
		Name:      "888888",
		Download:  0,
		Timestamp: 99999,
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_movieDB_FetchTMDBVideoByType(t *testing.T) {
	tt := types.VideoTypeTV
	list, err := NewMovieDB().FetchTMDBVideoByType(tt)
	if err != nil {
		t.Error(err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

func Test_movieDB_FetchOneTMDBVideoByImdbID(t *testing.T) {
	video, err := NewMovieDB().FetchOneTMDBVideoByImdbID(context.Background(), "tt30222734")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(video)
}

func Test_movieDB_FetchTVMagnetByName(t *testing.T) {
	videos, err := NewMovieDB().GetFeedVideoTVByNames([]string{"Ahsoka"}...)
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		fmt.Println(video)
	}
}

func Test_movieDB_FindLikeTVFromFeedVideo(t *testing.T) {
	videos, err := NewMovieDB().FindLikeTVFromFeedVideo("Raven.of.the.Inner.Palace")
	if err != nil {
		t.Error(err)
	}
	//for _, video := range videos {
	//	fmt.Println(video)
	//}
	fmt.Println(len(videos))
}

func Test_movieDB_RandomOneTMDBVideo(t *testing.T) {
	video, err := NewMovieDB().RandomOneTMDBVideo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(video)
}

func Test_movieDB_UpdateTMDBVideo(t *testing.T) {
	NewMovieDB().UpdateTMDBVideo(context.Background(), &types.TMDBVideo{
		ID:        99119,
		Names:     `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
		ImdbID:    "tt5555555",
		RowData:   "5555555",
		Timestamp: 0,
		Type:      "444444444",
		Playable:  "三21问问是岁",
	})
}

func Test_movieDB_UpdateFeedVideoDownloadByID(t *testing.T) {
	err := NewMovieDB().UpdateFeedVideoDownloadByID(56, 543543)
	if err != nil {
		t.Error(err)
	}
}

func Test_movieDB_checkDownloadHistory(t *testing.T) {
	_, flag := NewMovieDB().checkDownloadHistory(&types.DownloadHistory{
		Name:        "Raven.of.the.Inner.Palace",
		Type:        "tv",
		TorrentName: "Raven.of.the.Inner.Palace.S01E03.1080p.WEB.H264-SENPAI",
		Timestamp:   0,
		Resolution:  1080,
		Season:      "01",
		Episode:     "03",
	})
	fmt.Println(flag)
}

func Test_movieDB_IsDatePublished(t *testing.T) {
	obj := types.TMDBVideo{
		ID:            99119,
		Names:         `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
		ImdbID:        "tt444444444",
		RowData:       "444444444",
		Timestamp:     0,
		Type:          "444444444",
		Playable:      "三21问问是岁",
		DatePublished: "2023-06-13",
	}
	fmt.Println(obj.IsDatePublished())
}
