package feedspider

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/youcd/toolkit/log"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"regexp"
	"strings"
	"sync"
)

const (
	urlBaseTgx   = "https://tgx.rs"
	urlRssURITgx = "rss"
)

type Tgx struct {
	BaseFeeder
}

//nolint:forcetypeassert
func NewTgx(args ...interface{}) *Tgx {
	var url string
	if len(args) == 0 {
		url = fmt.Sprintf("%s/%s", urlBaseTgx, urlRssURITgx)
	} else if args[1] != "" {
		url = fmt.Sprintf("%s/%s", args[1], urlRssURITgx)
	}

	return &Tgx{
		BaseFeeder{
			web:        "tgx",
			url:        url,
			scheduling: args[0].(string),
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

func (t *Tgx) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(t.url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(t.web), fd.String())

	videos1 := make([]*types.FeedVideo, 0)
	for _, v := range fd.Items {
		if len(v.Categories) > 0 {
			if inSkipCategories(v.Categories[0]) {
				log.Infof("TGx: 跳过类别: [%s], Title: %s", v.Categories[0], v.Title)
				continue
			}
		}

		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		var name, year, typ string
		compileRegex := regexp.MustCompile(`(.*)\.(\d{4})\.`)
		matchArr := compileRegex.FindStringSubmatch(torrentName)
		if len(matchArr) < 3 {
			continue
		}
		year = matchArr[2]
		if len(matchArr) == 0 {
			tvReg := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9][eE][0-9][0-9])`)
			TVNameArr := tvReg.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(TVNameArr) == 0 {
				continue
			}
			name = TVNameArr[1]
		} else {
			name = matchArr[1]
		}

		// 过滤掉 其他类型的种子
		switch {
		case strings.HasPrefix(strings.ToLower(v.Categories[0]), "tv :"):
			typ = "tv"
		case strings.HasPrefix(strings.ToLower(v.Categories[0]), "movies :"):
			typ = "movie"
		default:
			continue
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = t.web
		fVideo.Year = year

		// 片名
		fVideo.Name = fVideo.FormatName(name)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)
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
			return nil, errors.WithMessage(err, "FetchMagnet")
		}
		video.Magnet = magnet
		videos = append(videos, video)
		wg.Done()
	}
	wg.Wait()
	//nolint:nakedret
	return
}
