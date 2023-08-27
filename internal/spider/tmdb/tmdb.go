package tmdb

import (
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/robfig/cron/v3"
	"io"
	"movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	url                      = "https://api.themoviedb.org/3/"
	patternFindByImdbID      = url + "find/%s?language=zh-CN&external_source=imdb_id&api_key=%s" // 通过imdbID获取
	patternFindMovieByID     = url + "movie/%d?language=en-us&api_key=%s"                        // 通过id获取movie
	patternFindMovieByIDzhCN = url + "movie/%d?language=zh-CN&api_key=%s"                        // 通过id获取movie
	patternFindTVByID        = url + "tv/%d?language=en-us&api_key=%s"                           // 通过id获取tv
	patternFindTVByIDzhCN    = url + "tv/%d?language=zh-CN&api_key=%s"                           // 通过id获取tv
)

type TmDB struct {
	url        string
	apikey     string
	scheduling string
	client     *http.Client
}

func NewSpiderTmDB(scheduling, apikey string) *TmDB {
	return &TmDB{
		url:        url,
		apikey:     apikey,
		scheduling: scheduling,
		client:     httpClient.NewHttpClient(),
	}
}

func (t *TmDB) FindByImdbID(imdbID string) (*types.TmDBFindByImdbIdData, error) {
	res := types.TmDBFindByImdbIdData{}

	urlStr := fmt.Sprintf(patternFindByImdbID, imdbID, t.apikey)
	err := t.request(urlStr, &res)
	if err != nil {
		return nil, err
	}

	return &res, err
}

func (t *TmDB) GetMovieDetailByID(id int, zhCN bool) (*types.TmDBMovieDetailData, error) {
	var tmDBResult types.TmDBMovieDetailData
	var urlStr string
	if zhCN {
		urlStr = fmt.Sprintf(patternFindMovieByIDzhCN, id, t.apikey)
	} else {
		urlStr = fmt.Sprintf(patternFindMovieByID, id, t.apikey)
	}

	err := t.request(urlStr, &tmDBResult)
	if err != nil {
		return nil, err
	}
	return &tmDBResult, err
}

func (t *TmDB) GetTVDetailByID(id int, zhCN bool) (*types.TmDBTVDetailData, error) {

	var tv types.TmDBTVDetailData
	var urlStr string
	if zhCN {
		urlStr = fmt.Sprintf(patternFindTVByIDzhCN, id, t.apikey)
	} else {
		urlStr = fmt.Sprintf(patternFindTVByID, id, t.apikey)
	}
	err := t.request(urlStr, &tv)
	if err != nil {
		return nil, err
	}
	return &tv, err
}

func (t *TmDB) request(urlStr string, result interface{}) error {
	resp, err := t.client.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(all, result)
	if err != nil {
		log.Debug(err)
		return err
	}

	return nil

}

func (t *TmDB) Crawler() {
	// 1. 获取 豆瓣想看列表的所有电视剧
	list, err := model.NewMovieDB().FetchDouBanVideoByType(types.VideoTypeTV)
	if err != nil {
		log.Error(err)
		return
	}
	log.Error("获取到电视剧数量：", len(list))
	//2	遍历所有电视剧 ，获取 tmDB信息 详情
	dataMap := make(map[*types.DouBanVideo]*types.TmDBFindByImdbIdData)

	// 将 list 中的名字提取出来
	videoDoubanNamesMap := make(map[*types.DouBanVideo][]string)

	for v, n := range list {
		videoDoubanNamesMap[v] = append(videoDoubanNamesMap[v], n...)
		got, err := t.FindByImdbID(v.ImdbID)
		if err != nil {
			log.Error(err)
			continue
		}
		dataMap[v] = got
	}
	//	3. 遍历所有tmDB，获取详情，以及名字

	for video, got := range dataMap {
		if len(got.TvEpisodeResults) > 0 {
			tv, err := t.GetTVDetailByID(got.TvEpisodeResults[0].ShowId, false)
			if err != nil {
				log.Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], allName(tv.Name, tv.NumberOfSeasons))
			tv, err = t.GetTVDetailByID(got.TvEpisodeResults[0].ShowId, true)
			if err != nil {
				log.Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], joinDot(tv.Name))
		}

		if len(got.TvResults) > 0 {
			tv, err := t.GetTVDetailByID(got.TvResults[0].Id, false)
			if err != nil {
				log.Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], allName(tv.Name, tv.NumberOfSeasons))

			tv, err = t.GetTVDetailByID(got.TvResults[0].Id, true)
			if err != nil {
				log.Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], joinDot(tv.Name))
		}

	}

	// 4. 遍历 videoTMDBNamesMap ，去掉重复的名字
	for video, names := range videoDoubanNamesMap {
		log.Errorf("video1: %#v    names：   %#v", video.ImdbID, names)
		videoDoubanNamesMap[video] = slice.Unique(names)

		marshal, err := json.Marshal(slice.Unique(names))
		if err != nil {
			log.Error(err)
		}

		log.Errorf("video3: %#v    names：   %#v", video.ImdbID, string(marshal))

	}
}

func (t *TmDB) Run() {
	//TODO implement me
	panic("implement me")

	if t.scheduling == "" {
		log.Error("DouBan Scheduling is null")
		os.Exit(1)
	}
	log.Infof("DouBan Scheduling is: [%s]", t.scheduling)
	c := cron.New()
	c.AddFunc(t.scheduling, func() { t.Crawler() })
	c.Start()

}

func allName(str string, seasons int) string {
	name := joinDot(str)
	if seasons < 10 {
		return name + ".S0" + strconv.Itoa(seasons)
	} else {
		return name + ".S" + strconv.Itoa(seasons)
	}
}

func joinDot(str string) string {
	return strings.Join(strings.Split(str, " "), ".")
}
