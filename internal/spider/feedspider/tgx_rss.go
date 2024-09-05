package feedspider

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"strings"
	"sync"

	"github.com/youcd/toolkit/log"
)

type Tgx struct {
	BaseFeeder
}

func NewTgx(scheduling, siteURL string, useIPProxy bool) *Tgx {
	return &Tgx{
		BaseFeeder{
			web: "tgx",
			BaseFeed: types.BaseFeed{
				Url:        siteURL,
				Scheduling: scheduling,
				UseIPProxy: useIPProxy,
			},
		},
	}
}

//nolint:gochecknoglobals
var (
	//  跳过的类别
	skipCategories = []string{"games", "xxx", "apps", "music", "books"}
)

func inSkipCategories(categories string) bool {
	for _, category := range skipCategories {
		return strings.Contains(strings.ToLower(categories), category)
	}
	return false
}

func (t *Tgx) Crawler() (videos []*types.FeedVideoBase, err error) {
	fd, err := t.FeedParser().ParseURL(t.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(t.web), fd.String())

	videos1 := make([]*types.FeedVideoBase, 0)
	for _, v := range fd.Items {
		if len(v.Categories) > 0 {
			if inSkipCategories(v.Categories[0]) {
				log.Debugf("TGx: 跳过类别: [%s], Title: %s", v.Categories[0], v.Title)
				continue
			}
		}

		var typ string
		switch {
		case strings.HasPrefix(strings.ToLower(v.Categories[0]), "tv :"):
			typ = "tv"
		case strings.HasPrefix(strings.ToLower(v.Categories[0]), "movies :"):
			typ = "movie"
		default:
			// 过滤掉 其他类型的种子
			continue
		}

		fVideo := new(types.FeedVideoBase)
		fVideo.Web = t.web
		// 种子名
		fVideo.TorrentName = v.Title
		fVideo.Type = typ
		fVideo.TorrentURL = v.Link
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)

		fVideo.RowData = sql.NullString{String: string(bytes)}
		videos1 = append(videos1, fVideo)
	}
	var wg sync.WaitGroup
	for _, video := range videos1 {
		wg.Add(1)
		magnet, err := magnetconvert.FetchMagnet(video.TorrentURL)
		if err != nil {
			wg.Done()
			return nil, fmt.Errorf("FetchMagnet: %w", err)
		}
		video.Magnet = magnet
		videos = append(videos, video)
		wg.Done()
	}
	wg.Wait()
	//nolint:nakedret
	return
}
