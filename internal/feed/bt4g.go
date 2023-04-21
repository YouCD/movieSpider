package feed

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"movieSpider/internal/log"
	types2 "movieSpider/internal/types"
	"net/url"
	"os"
	"strings"
)

const (
	urlBt4g = "https://bt4g.org"
)

type bt4g struct {
	url        string
	resolution types2.Resolution
	web        string
}

func NewFeedBt4g(name string, resolution types2.Resolution) *bt4g {
	parse, err := url.Parse(urlBt4g)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	bUrl := fmt.Sprintf("%s://%s/search/%s/bysize/1?page=rss", parse.Scheme, parse.Host, name)
	return &bt4g{url: bUrl, resolution: resolution, web: "bt4g"}
}

func (b *bt4g) Crawler() (videos []*types2.FeedVideo, err error) {
	f := gofeed.NewParser()
	fd, err := f.ParseURL(b.url)
	if fd == nil {
		return nil, errors.New("BT4G: 没有feed数据")
	}
	log.Debugf("BT4G Config: %#v", b)
	log.Debugf("BT4G Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("BT4G: 没有feed数据")
	}
	for _, v := range fd.Items {
		// 片名
		name := strings.ReplaceAll(v.Title, " ", ".")
		ok := excludeVideo(name)
		if ok {
			continue
		}
		if v.Link == "" {
			continue
		}

		fVideo := new(types2.FeedVideo)
		fVideo.Web = b.web
		fVideo.Name = fVideo.FormatName(name)
		fVideo.Magnet = v.Link
		// 种子名
		fVideo.TorrentName = fVideo.Name

		fVideo.TorrentUrl = v.GUID
		fVideo.Type = "other"
		bytes, _ := json.Marshal(v)
		fVideo.RowData = sql.NullString{String: string(bytes)}
		videos = append(videos, fVideo)
	}

	proxySaveVideo2DB(videos...)
	// 指定清晰度
	if b.resolution != types2.ResolutionOther {
		var resolutionVideos []*types2.FeedVideo
		for _, v := range videos {
			if strings.Contains(v.Name, b.resolution.Res()) {
				resolutionVideos = append(resolutionVideos, v)
			}
		}
		return resolutionVideos, nil
	}
	return
}
func (b *bt4g) Run() {

}
