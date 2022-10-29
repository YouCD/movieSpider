package model

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	types2 "movieSpider/internal/types"
	"regexp"
	"strings"
	"testing"
)

func Test_movieDB_FetchMagnetByName(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")

	NewMovieDB()

	var names = []string{"Im.Westen.nichts.Neues"}
	videos, err := movieDatabase.FetchMovieMagnetByName(names)
	if err != nil {
		t.Error(err)
	}
	for _, v := range videos {
		fmt.Println(v.Magnet)
	}

}

func Test_movieDB_UpdateFeedVideoDownloadByID(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	sql := "select  id,name,torrent_name  from feed_video where web=? and torrent_name like ?"
	query, err := NewMovieDB().db.Query(sql, "btbt", fmt.Sprintf("https%%"))
	if err != nil {
		t.Error(err)
	}

	for query.Next() {
		var v types2.FeedVideo

		err := query.Scan(&v.ID, &v.Name, &v.TorrentName)
		if err != nil {
			t.Error(err)
		}
		_, err = NewMovieDB().db.Exec("update feed_video set torrent_name=? where id=?", v.Name, v.ID)
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				_, err := NewMovieDB().db.Exec("delete  from  feed_video  where id=?", v.ID)
				if err != nil {
					t.Error(err)
				}
			}
			t.Error(err)
		}
		fmt.Println(v)
	}
	//
	//err := movieDatabase.UpdateFeedVideoDownloadByID(55)
	//if err != nil {
	//	t.Error(err)
	//}
	isURL := govalidator.IsURL("https://api.wmdb.tv/movie/api?id=")
	//parse, err := url.ParseRequestURI("http://dsadsa.com")
	//if err != nil {
	//	t.Error(err)
	//}
	fmt.Println(isURL)
}

func Test_movieDB_RandomOneDouBanVideo(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	NewMovieDB()
	movieDatabase.RandomOneDouBanVideo()
}

func Test_movieDB_FetchTVMagnetByName(t *testing.T) {

	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	NewMovieDB()

	var names = []string{"Black.Adam"}
	videos, err := movieDatabase.FetchTVMagnetByName(names)
	if err != nil {
		t.Error(err)
	}
	is1, is3 := sotByResolution(videos)
	for _, video := range is1 {
		fmt.Println(video)
	}
	fmt.Println("-----------------")
	for _, video := range is3 {
		fmt.Println(video)

	}

}
func sotByResolution(videos []*types2.FeedVideo) (downloadIs1 []*types2.FeedVideo, downloadIs3 []*types2.FeedVideo) {
	var Videos2160P []*types2.FeedVideo
	var Videos1080P []*types2.FeedVideo
	for _, v := range videos {
		switch {
		case strings.Contains(v.TorrentName, "2160"):
			Videos2160P = append(Videos2160P, v)
		case strings.Contains(v.TorrentName, "1080"):
			Videos1080P = append(Videos1080P, v)
		}
	}
	if len(Videos2160P) >= 0 {
		if len(Videos2160P) >= 2 {
			downloadIs1 = append(downloadIs1, Videos2160P[0:2]...)
			downloadIs3 = append(downloadIs3, Videos2160P[2:]...)
			downloadIs3 = append(downloadIs3, Videos1080P...)
		} else {
			downloadIs1 = append(downloadIs1, Videos2160P...)
			downloadIs3 = append(downloadIs3, Videos1080P...)
		}

	} else {
		if len(Videos2160P) >= 2 {
			downloadIs1 = append(downloadIs1, Videos1080P[0:2]...)
			downloadIs3 = append(downloadIs3, Videos1080P[2:]...)
		} else {
			downloadIs1 = append(downloadIs1, Videos1080P...)
		}

	}
	return
}

func Test_movieDB_FetchDouBanVideoByType(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	NewMovieDB()

	Videos, err := movieDatabase.FetchDouBanVideoByType(types2.ResourceMovie)
	if err != nil {
		t.Error(err)
	}

	for _, v := range Videos {
		fmt.Println(v)
	}
}

// 数据统计
func Test_movieDB_CountFeedVideo(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	NewMovieDB()
	count, err := movieDatabase.CountFeedVideo()
	if err != nil {
		t.Error(err)
	}
	var s string
	var Total int

	for _, reportCount := range count {
		Total += reportCount.Count
		s += fmt.Sprintf("%s: %d ", reportCount.Web, reportCount.Count)
	}
	fmt.Println(s)
	log.Infof("Report: feed_video 数据统计: Total: %d  %s", Total, s)
}

// 批量修改名字
func Test_movieDB_FindLikeTVFromFeedVideo(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	NewMovieDB()
	videos, err := movieDatabase.FindLikeTVFromFeedVideo("S19")
	if err != nil {
		t.Error(err)
	}
	compileRegex := regexp.MustCompile("(.*)\\.[sS][0-9][0-9]|[Ee][0-9][0-9]?\\.")
	for _, video := range videos {

		matchArr := compileRegex.FindStringSubmatch(video.Name)
		if len(matchArr) > 0 {
			err := movieDatabase.UpdateFeedVideoNameByID(video.ID, matchArr[1], types2.ResourceTV)
			if err != nil {
				t.Error(err)
			}
			log.Infof("ID: %d name: %s 修改为 %s", video.ID, video.Name, matchArr[1])

		} else {
			log.Warnf("ID:%d name: %s 无需修改", video.ID, video.Name)
		}

	}
}
