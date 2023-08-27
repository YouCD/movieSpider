package tmdb

import (
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"strings"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")

}

func TestTmDB_FindByImdbID(t1 *testing.T) {
	videos, err := model.NewMovieDB().FetchThisYearVideo()
	if err != nil {
		t1.Logf("%+v", err)
	}
	t := NewSpiderTmDB(config.TmDB.Scheduling, config.TmDB.ApiKey)
	//log.Infof("len(videos)  %d  %v   ", len(videos), videos)
	for _, video := range videos {
		//log.Infof("video   %s   %s    %s ", video.Type, video.Names, video.ImdbID)
		//
		got, err := t.FindByImdbID(video.ImdbID)
		if err != nil {
			log.Errorf("FindByImdbID() error = %v", err)
			continue
		}

		switch types.Convert2VideoType(video.Type) {
		case types.VideoTypeTV:

			if len(got.TvEpisodeResults) > 0 {
				tv, err := t.GetTVDetailByID(got.TvEpisodeResults[0].ShowId, false)
				if err != nil {
					t1.Errorf("TvEpisodeResults() error = %v", err)
					continue
				}
				log.Errorf("TvEpisodeResults  %s    imdb: %s    seasons： %v   TvEpisodeResults: %v  全部是英文名: %v  ", tv.Name, video.ImdbID, tv.NumberOfSeasons, got.TvEpisodeResults[0].ShowId, isEnglishString(strings.ReplaceAll(tv.Name, " ", "")))

				tv, err = t.GetTVDetailByID(got.TvEpisodeResults[0].ShowId, true)
				if err != nil {
					t1.Errorf("TvEpisodeResults() error = %v", err)
					continue
				}
				log.Errorf("TvEpisodeResults  %s    imdb: %s    seasons： %v   TvEpisodeResults: %v 全部是英文名: %v  ", tv.Name, video.ImdbID, tv.NumberOfSeasons, got.TvEpisodeResults[0].ShowId, isEnglishString(strings.ReplaceAll(tv.Name, " ", "")))

				continue
			}

			if len(got.TvResults) > 0 {
				//log.Warn("TvResults   ", got.TvResults)
				tv, err := t.GetTVDetailByID(got.TvResults[0].Id, false)
				if err != nil {
					t1.Errorf("TvResults() error = %v", err)
					continue
				}
				log.Errorf("TvResults  %s  seasons： %v   全部是英文名: %v ", tv.Name, tv.NumberOfSeasons, isEnglishString(strings.ReplaceAll(tv.Name, " ", "")))

				tv, err = t.GetTVDetailByID(got.TvResults[0].Id, true)
				if err != nil {
					t1.Errorf("TvResults() error = %v", err)
					continue
				}

				log.Errorf("TvResults  %s  seasons： %v  全部是英文名: %v  ", tv.Name, tv.NumberOfSeasons, isEnglishString(strings.ReplaceAll(tv.Name, " ", "")))
				continue
			}

			if len(got.TvSeasonResults) > 0 {
				//log.Warn("TvSeasonResults    ", got.TvSeasonResults)
				//id, err := t.GetTVDetailByID(got.TvSeasonResults[0].TvId)
				//if err != nil {
				//	t1.Errorf("TvSeasonResults() error = %v", err)
				//	continue
				//}
				log.Warnf("TvSeasonResults  %s ", video.Names)
				continue
			}

			//log.Warn("VideoTypeTV    ", got)
			//case types.VideoTypeMovie:
			//	if len(got.MovieResults) > 0 {
			//		Movie, err := t.GetMovieDetailByID(got.MovieResults[0].Id)
			//		if err != nil {
			//			t1.Errorf("GetMovieDetailByID() error = %v", err)
			//			continue
			//		}
			//		log.Error("Movie  ", Movie)
			//	}
			//default:
			//	log.Warn("default    ", got)
		}
	}

	//got, err := t.FindByImdbID("tt13622776")
	//if err != nil {
	//	t1.Errorf("FindByImdbID() error = %v", err)
	//	return
	//}
	//t1.Logf("%+v", got)

}
func isEnglishChar(char rune) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}
func isEnglishString(str string) bool {
	for _, char := range str {
		if !isEnglishChar(char) {
			return false
		}
	}
	return true
}

func TestTmDB_Crawler(t1 *testing.T) {

	t := NewSpiderTmDB(config.TmDB.Scheduling, config.TmDB.ApiKey)

	t.Crawler()

}
