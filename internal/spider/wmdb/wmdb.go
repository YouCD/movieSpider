package wmdb

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
	"io"
	"movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"strings"
)

var (
	client *http.Client
)

const urlStr = "https://api.wmdb.tv/movie/api?id="

type WMDB struct {
	url        string
	scheduling string
}

//
// NewSpiderWmdb
//  @Description: 初始化
//  @param scheduling
//  @return *WMDB
//
func NewSpiderWmdb(scheduling string) *WMDB {
	return &WMDB{
		urlStr,
		scheduling,
	}
}

//
// crawlerImdb 30s 内只允许一个请求
//  @Description:
//  @receiver d
//  @param doubanID
//  @return video
//  @return err
//
// crawlerImdb 30s 内只允许一个请求
func (d *WMDB) crawler(doubanID string) (video *types.DouBanVideo, err error) {
	log.Infof("WMDB: crawler Douban ID: %s", doubanID)
	video = new(types.DouBanVideo)
	urlStr := d.url + doubanID
	log.Debugf("WMDB: url: %s", urlStr)

	request, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, nil
	}
	client = httpClient.NewHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp == nil {
		log.Warn("未能正常获取wmdb数据")
		return nil, errors.New("未能正常获取wmdb数据")
	}

	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("WMDB Config: %#v", string(all))
	rowData := string(all)

	if strings.Contains(rowData, "Too Many Requests") {
		return nil, errors.New(fmt.Sprintf("WMDB: 没有Spider数据,url:%s", urlStr))
	}
	if strings.Contains(rowData, "your requests today are full") {
		return nil, errors.New(fmt.Sprintf("WMDB: 没有Spider数据,url:%s", urlStr))
	}
	if strings.Contains(rowData, "Bad Request") {
		return nil, errors.New(fmt.Sprintf("WMDB: 没有Spider数据,url:%s", urlStr))
	}
	video.ImdbID = gjson.Get(rowData, "imdbId").String()
	video.Type = strings.ToLower(gjson.Get(rowData, "type").String())
	video.RowData = rowData

	var ns []string
	array := gjson.Get(rowData, "data").Array()
	for _, v := range array {
		name := gjson.Get(v.String(), "name").String()
		replace := strings.ReplaceAll(name, " ", ".")
		ns = append(ns, replace)
	}
	marshal, _ := json.Marshal(ns)
	video.Names = string(marshal)

	return
}

func (d *WMDB) Run() {
	if d.scheduling == "" {
		log.Error("WMDB: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("WMDB Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	c.AddFunc(d.scheduling, func() {
		video, err := model.NewMovieDB().RandomOneDouBanVideo()
		if err != nil {
			log.Warn("WMDB: 没有可爬取的豆瓣数据")
			return
		}

		v, err := d.crawler(video.DoubanID)
		if err != nil {
			log.Error(err)
			client = httpClient.NewProxyHttpClient("https")
			v, err := d.crawler(video.DoubanID)
			if err != nil {
				log.Error(err)
				return
			}
			d.updateVideo(video, v)
			return
		}
		d.updateVideo(video, v)

	})
	c.Start()
}

func (d *WMDB) updateVideo(video *types.DouBanVideo, crawlerVideo *types.DouBanVideo) {
	video.ImdbID = crawlerVideo.ImdbID
	video.RowData = crawlerVideo.RowData
	if strings.ToLower(crawlerVideo.Type) == "tvseries" {
		video.Type = "tv"
	} else {
		video.Type = crawlerVideo.Type
	}

	video.Names = crawlerVideo.Names

	err := model.NewMovieDB().UpdateDouBanVideo(video)
	if err != nil {
		log.Error(err)
		return
	}
	log.Warnf("WMDB: %s 更新完毕", video.Names)
}
