package searchspider

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/url"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/youcd/toolkit/log"
)

const (
	urlBt4g = "https://bt4g.org"
)

type BT4g struct {
	url        string
	resolution types.Resolution
	web        string
}

func NewFeedBt4g(name string, resolution types.Resolution) *BT4g {
	parse, err := url.Parse(urlBt4g)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	bURL := fmt.Sprintf("%s://%s/search/%s/bysize/1?page=rss", parse.Scheme, parse.Host, name)
	return &BT4g{url: bURL, resolution: resolution, web: "BT4g"}
}

//nolint:nakedret
func (b *BT4g) Search() (videos []*types.FeedVideo, err error) {
	f := gofeed.NewParser()
	fd, err := f.ParseURL(b.url)
	if fd == nil {
		//nolint:goerr113
		return nil, errors.New("BT4G: 没有feed数据")
	}
	log.Debugf("BT4G Config: %#v", b)
	log.Debugf("BT4G Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		//nolint:goerr113
		return nil, errors.New("BT4G: 没有feed数据")
	}
	for _, v := range fd.Items {
		// 片名
		if v.Link == "" {
			continue
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = b.web
		fVideo.Magnet = v.Link
		// 种子名
		fVideo.TorrentName = v.Title

		fVideo.TorrentURL = v.GUID
		fVideo.Type = "other"
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)

		fVideo.RowData = sql.NullString{String: string(bytes)}
		videos = append(videos, fVideo)
	}

	model.ProxySaveVideo2DB(videos...)
	// 指定清晰度
	if b.resolution != types.ResolutionOther {
		var resolutionVideos []*types.FeedVideo
		for _, v := range videos {
			if strings.Contains(v.Name, b.resolution.Res()) {
				resolutionVideos = append(resolutionVideos, v)
			}
		}
		return resolutionVideos, nil
	}
	return
}
