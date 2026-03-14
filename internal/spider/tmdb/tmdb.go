package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
)

const (
	url                      = "https://api.themoviedb.org/3/"
	patternFindByImdbID      = url + "find/%s?language=zh-CN&external_source=imdb_id" // 通过imdbID获取
	patternFindMovieByID     = url + "movie/%d?language=en-us"                        // 通过id获取movie
	patternFindMovieByIDzhCN = url + "movie/%d?language=zh-CN"                        // 通过id获取movie
	patternFindTVByID        = url + "tv/%d?language=en-us"                           // 通过id获取tv
	patternFindTVByIDzhCN    = url + "tv/%d?language=zh-CN"                           // 通过id获取tv
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
		client:     httpclient.NewProxyHTTPClient(context.Background()),
	}
}

func (t *TmDB) FindByImdbID(imdbID string) (*types.TmDBFindByImdbIDData, error) {
	res := types.TmDBFindByImdbIDData{}

	urlStr := fmt.Sprintf(patternFindByImdbID, imdbID)
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
		urlStr = fmt.Sprintf(patternFindMovieByIDzhCN, id)
	} else {
		urlStr = fmt.Sprintf(patternFindMovieByID, id)
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
		urlStr = fmt.Sprintf(patternFindTVByIDzhCN, id)
	} else {
		urlStr = fmt.Sprintf(patternFindTVByID, id)
	}
	err := t.request(urlStr, &tv)
	if err != nil {
		return nil, err
	}
	return &tv, err
}
func (t *TmDB) Crawler() {
	// 1. 获取 豆瓣想看列表的所有电视剧
	list, err := model.NewMovieDB().FetchDouBanVideoByType(types.VideoTypeTV)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
		return
	}
	log.WithCtx(context.Background()).Error("获取到电视剧数量：", len(list))
	// 2	遍历所有电视剧 ，获取 tmDB信息 详情
	dataMap := make(map[*types.DouBanVideo]*types.TmDBFindByImdbIDData)

	// 将 list 中的名字提取出来
	videoDoubanNamesMap := make(map[*types.DouBanVideo][]string)

	for v, n := range list {
		videoDoubanNamesMap[v] = append(videoDoubanNamesMap[v], n...)
		got, err := t.FindByImdbID(v.ImdbID)
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			continue
		}
		dataMap[v] = got
	}
	//	3. 遍历所有tmDB，获取详情，以及名字

	for video, got := range dataMap {
		if len(got.TvEpisodeResults) > 0 {
			tv, err := t.GetTVDetailByID(got.TvEpisodeResults[0].ShowID, false)
			if err != nil {
				log.WithCtx(context.Background()).Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], allName(tv.Name, tv.NumberOfSeasons))
			tv, err = t.GetTVDetailByID(got.TvEpisodeResults[0].ShowID, true)
			if err != nil {
				log.WithCtx(context.Background()).Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], joinDot(tv.Name))
		}

		if len(got.TvResults) > 0 {
			tv, err := t.GetTVDetailByID(got.TvResults[0].ID, false)
			if err != nil {
				log.WithCtx(context.Background()).Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], allName(tv.Name, tv.NumberOfSeasons))

			tv, err = t.GetTVDetailByID(got.TvResults[0].ID, true)
			if err != nil {
				log.WithCtx(context.Background()).Error(err)
				continue
			}
			videoDoubanNamesMap[video] = append(videoDoubanNamesMap[video], joinDot(tv.Name))
		}
	}

	// 4. 遍历 videoTMDBNamesMap ，去掉重复的名字
	for video, names := range videoDoubanNamesMap {
		log.WithCtx(context.Background()).Errorf("video1: %#v    names：   %#v", video.ImdbID, names)
		videoDoubanNamesMap[video] = slice.Unique(names)

		marshal, err := json.Marshal(slice.Unique(names))
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
		}

		log.WithCtx(context.Background()).Errorf("video3: %#v    names：   %#v", video.ImdbID, string(marshal))
	}
}

func (t *TmDB) Run() {
	if t.scheduling == "" {
		log.WithCtx(context.Background()).Error("DouBan Scheduling is null")
		os.Exit(1)
	}
	log.WithCtx(context.Background()).Infof("DouBan Scheduling is: [%s]", t.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(t.scheduling, func() { t.Crawler() })
	c.Start()
}

//nolint:noctx
func (t *TmDB) request(urlStr string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败,err:%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apikey)

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败,err:%w", err)
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取body失败,err:%w", err)
	}

	if err = json.Unmarshal(all, result); err != nil {
		log.WithCtx(context.Background()).Debug(err)
		return fmt.Errorf("解析json失败,err:%w", err)
	}

	return err
}

func allName(str string, seasons int) string {
	name := joinDot(str)
	if seasons < 10 {
		return name + ".S0" + strconv.Itoa(seasons)
	}
	return name + ".S" + strconv.Itoa(seasons)
}

func joinDot(str string) string {
	return strings.Join(strings.Split(str, " "), ".")
}
