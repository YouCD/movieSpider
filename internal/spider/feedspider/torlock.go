package feedspider

import (
	"database/sql"
	"encoding/json"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Torlock struct {
	typ types.VideoType
	BaseFeeder
}

func NewTorlock(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) *Torlock {
	return &Torlock{
		typ: resourceType,
		BaseFeeder: BaseFeeder{
			web:      "torlock",
			BaseFeed: types.BaseFeed{Url: siteURL, Scheduling: scheduling, UseIPProxy: useIPProxy},
		},
	}
}

func (t *Torlock) Crawler() ([]*types.FeedVideoBase, error) {
	var Videos []*types.FeedVideoBase
	fp := t.FeedParserUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	if t.typ == types.VideoTypeMovie {
		fd, err := fp.ParseURL(t.Url)
		if err != nil {
			return nil, ErrFeedParseURL
		}
		log.Debugf("%s Data: %#v", strings.ToUpper(t.web), fd.String())

		var videos1 []*types.FeedVideoBase
		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")
			var fVideo types.FeedVideoBase
			fVideo.Web = t.web
			fVideo.TorrentName = name
			fVideo.TorrentURL = v.Link
			fVideo.Type = "movie"

			// 原始数据
			//nolint:errchkjson
			bytes, _ := json.Marshal(v)
			fVideo.RowData = sql.NullString{String: string(bytes)}
			videos1 = append(videos1, &fVideo)
		}

		videos2 := t.fetchMagnetDownLoad(videos1)

		Videos = t.fetchMagnet(videos2)
		return Videos, nil
	}
	if t.typ == types.VideoTypeTV {
		fd, err := fp.ParseURL(t.Url)
		if err != nil {
			return nil, ErrFeedParseURL
		}
		log.Debugf("TORLOCK.tv Data: %#v", fd.String())
		var videos1 []*types.FeedVideoBase
		for _, v := range fd.Items {
			var fVideo types.FeedVideoBase
			fVideo.TorrentName = v.Title
			fVideo.TorrentURL = v.Link
			fVideo.Type = "tv"
			//nolint:errchkjson
			bytes, _ := json.Marshal(v)

			fVideo.RowData = sql.NullString{String: string(bytes)}
			fVideo.Web = t.web
			videos1 = append(videos1, &fVideo)
		}

		videos2 := t.fetchMagnetDownLoad(videos1)

		Videos = t.fetchMagnet(videos2)
		return Videos, nil
	}
	return nil, nil
}

func (t *Torlock) fetchMagnet(videos []*types.FeedVideoBase) (feedVideos []*types.FeedVideoBase) {
	var wg sync.WaitGroup
	for _, video := range videos {
		wg.Add(1)
		magnet, err := magnetconvert.FetchMagnet(video.Magnet)
		if err != nil {
			log.Errorf("TORLOCK: get %s magnet download url is %s", video.TorrentName, video.Magnet)
			wg.Done()
			continue
		}
		video.Magnet = magnet
		feedVideos = append(feedVideos, video)
		wg.Done()
	}
	wg.Wait()
	return feedVideos
}

func (t *Torlock) fetchMagnetDownLoad(videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	var wg sync.WaitGroup
	var videos2 []*types.FeedVideoBase
	for _, video := range videos {
		wg.Add(1)
		//nolint:noctx
		go func() {
			defer wg.Done()
			resp, err := t.HTTPClientDynamic().Get(video.TorrentURL)
			if err != nil {
				log.Errorf("TORLOCK.%s %s http request url is %s, error:%s", video.Type, video.TorrentName, video.TorrentURL, err)
				return
			}
			defer resp.Body.Close()
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Errorf("TORLOCK.%s %s http goquery error:%s", video.Type, video.TorrentName, err)
				return
			}
			val, exists := doc.Find("body > article > div:nth-child(6) > div > div:nth-child(2) > a").Attr("href")
			if exists {
				video.Magnet = val
				videos2 = append(videos2, video)
			}
		}()
	}
	wg.Wait()
	return videos2
}
