package model

import (
	"database/sql"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"testing"
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

func Test_movieDB_CreatDouBanVideo(t *testing.T) {
	err := NewMovieDB().CreatDouBanVideo(&types.DouBanVideo{
		ID:        99119,
		Names:     `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
		DoubanID:  "878",
		ImdbID:    "444444444",
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
		ID:          888888,
		Name:        "888888",
		TorrentName: "888888",
		TorrentURL:  "888888",
		Magnet:      "888888",
		Year:        "888888",
		Type:        "888888",
		RowData:     sql.NullString{},
		Web:         "888888",
		Download:    0,
		Timestamp:   99999,
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_movieDB_FetchDouBanVideoByType(t *testing.T) {
	var tt = types.VideoTypeMovie
	list, err := NewMovieDB().FetchDouBanVideoByType(tt)

	if err != nil {
		t.Error(err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

func Test_movieDB_FetchMovieMagnetByName(t *testing.T) {
	videos, err := NewMovieDB().GetFeedVideoMovieByName("满江红")
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		fmt.Println(video)
	}
}

func Test_movieDB_FetchOneDouBanVideoByDouBanID(t *testing.T) {
	video, err := NewMovieDB().FetchOneDouBanVideoByDouBanID("30222734")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(video)
}

func Test_movieDB_FetchTVMagnetByName(t *testing.T) {
	videos, err := NewMovieDB().GetFeedVideoTVByName("29384657", []string{"Ahsoka"}...)
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
	for _, video := range videos {
		fmt.Println(video)
	}
}

func Test_movieDB_RandomOneDouBanVideo(t *testing.T) {
	video, err := NewMovieDB().RandomOneDouBanVideo()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(video)
}

func Test_movieDB_UpdateDouBanVideo(t *testing.T) {
	NewMovieDB().UpdateDouBanVideo(&types.DouBanVideo{
		ID:        99119,
		Names:     `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
		DoubanID:  "878",
		ImdbID:    "5555555",
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

func Test_movieDB_UpdateFeedVideoNameByID(t *testing.T) {
	NewMovieDB().UpdateFeedVideoNameByID(56, "万千瓦请", types.VideoTypeMovie)
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
	obj := types.DouBanVideo{
		ID:            99119,
		Names:         `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
		DoubanID:      "878",
		ImdbID:        "444444444",
		RowData:       "444444444",
		Timestamp:     0,
		Type:          "444444444",
		Playable:      "三21问问是岁",
		DatePublished: "2023-06-13",
	}
	fmt.Println(obj.IsDatePublished())
}
