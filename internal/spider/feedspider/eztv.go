package feedspider

import (
	"database/sql"
	"encoding/json"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"strings"

	"github.com/youcd/toolkit/log"
)

type Eztv struct {
	BaseFeeder
}

func NewEztv() *Eztv {
	return &Eztv{BaseFeeder{
		web: "eztv",
		BaseFeed: types.BaseFeed{
			Scheduling: config.Config.Feed.EZTV.Scheduling,
			Url:        config.Config.Feed.EZTV.Url,
			UseIPProxy: config.Config.Feed.EZTV.UseIPProxy,
		},
	}}
}

func (f *Eztv) Crawler() (videos []*types.FeedVideoBase, err error) {
	fd, err := f.FeedParser().ParseURL(f.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(f.web), fd.String())

	for _, v := range fd.Items {
		fVideo := new(types.FeedVideoBase)
		fVideo.Web = f.web
		fVideo.TorrentName = v.Title
		fVideo.TorrentURL = v.Link
		fVideo.Magnet = v.Extensions["torrent"]["magnetURI"][0].Value
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		fVideo.Type = strings.ToLower(v.Categories[0])
		fVideo.RowData = sql.NullString{String: string(bytes)}
		videos = append(videos, fVideo)
	}
	return
}
