package feedspider

import (
	"database/sql"
	"encoding/json"
	"movieSpider/internal/types"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/youcd/toolkit/log"
)

type Eztv struct {
	BaseFeeder
}

func NewEztv(scheduling, siteURL string) *Eztv {
	return &Eztv{BaseFeeder{
		web:        "eztv",
		url:        siteURL,
		scheduling: scheduling,
	}}
}

//nolint:gosimple
func (f *Eztv) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(f.url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(f.web), fd.String())

	for _, v := range fd.Items {
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		var name, year string
		compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
		matchArr := compileRegex.FindStringSubmatch(torrentName)
		if len(matchArr) == 0 {
			tvReg := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			TVNameArr := tvReg.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(TVNameArr) == 0 {
				continue
			}
			name = TVNameArr[1]
		} else {
			year = matchArr[2]
			name = matchArr[1]
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = f.web
		fVideo.Year = year

		// 片名
		fVideo.Name = fVideo.FormatName(name)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)
		fVideo.TorrentURL = v.Link
		fVideo.Magnet = v.Extensions["torrent"]["magnetURI"][0].Value
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		fVideo.Type = strings.ToLower(v.Categories[0])

		fVideo.RowData = sql.NullString{String: string(bytes)}

		videos = append(videos, fVideo)
	}
	//nolint:nakedret
	return
}
