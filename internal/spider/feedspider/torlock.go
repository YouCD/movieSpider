package feedspider

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
	"go.uber.org/zap/buffer"
)

type Torlock struct {
	BaseFeeder

	typ types.VideoType
}

func NewTorlock(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) Feeder {
	return &Torlock{
		typ: resourceType,
		BaseFeeder: BaseFeeder{
			web:      "torlock",
			BaseFeed: types.BaseFeed{Url: siteURL, Scheduling: scheduling, UseIPProxy: useIPProxy},
		},
	}
}

func (t *Torlock) Crawler(ctx context.Context) ([]*types.FeedVideoBase, error) {
	var Videos []*types.FeedVideoBase
	fp := t.FeedParserUserAgent(ctx, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	if t.typ == types.VideoTypeMovie {
		fd, err := fp.ParseURL(t.Url)
		if err != nil {
			return nil, fmt.Errorf("err:%s, err:%w", err, ErrFeedParseURL)
		}
		log.WithCtx(ctx).Debugf("%s type: %v Data: %s", t.web, t.typ, fd.String())
		var videos1 []*types.FeedVideoBase
		for _, v := range fd.Items {
			// 片名
			var fVideo types.FeedVideoBase
			fVideo.Web = t.web
			fVideo.TorrentName = v.Title
			fVideo.TorrentURL = v.Link
			fVideo.Type = "movie"

			// 原始数据
			//nolint:errchkjson
			bytes, _ := json.Marshal(v)
			fVideo.RowData = sql.NullString{String: string(bytes)}
			videos1 = append(videos1, &fVideo)
		}

		videos2 := t.fetchMagnetDownLoad(ctx, videos1)
		log.WithCtx(ctx).Infof("Movie Magnet Url Count: %d", len(videos2))

		Videos = t.fetchMagnet(ctx, videos2)
		log.WithCtx(ctx).Infof("Movie Magnet Url Count: %d", len(Videos))

		return Videos, nil
	}
	if t.typ == types.VideoTypeTV {
		fd, err := fp.ParseURL(t.Url)
		if err != nil {
			return nil, ErrFeedParseURL
		}
		log.WithCtx(ctx).Debugf("%s type: %v Data: %s", t.web, t.typ, fd.String())
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

		videos2 := t.fetchMagnetDownLoad(ctx, videos1)
		log.WithCtx(ctx).Infof("TV Magnet Url Count: %d", len(videos2))
		Videos = t.fetchMagnet(ctx, videos2)
		log.WithCtx(ctx).Infof("TV Magnet Url Count: %d", len(Videos))
		return Videos, nil
	}
	return nil, nil
}

func (t *Torlock) fetchMagnet(ctx context.Context, videos []*types.FeedVideoBase) (feedVideos []*types.FeedVideoBase) {
	var wg sync.WaitGroup
	for _, video := range videos {
		wg.Add(1)
		go func(video *types.FeedVideoBase) {
			defer wg.Done()
			var count int
		RETRY:
			if strings.HasPrefix(video.Magnet, "magnet:?x") {
				feedVideos = append(feedVideos, video)
				return
			}
			magnet, err := magnetconvert.FetchMagnetWithHTTPClient(ctx, video.Magnet, t.HTTPClientDynamic(ctx))
			if err != nil {
				if count < 3 {
					count++
					if count == 2 {
						log.WithCtx(ctx).Warnf("RETRY: torlock.%s %s http request url is %s retry count:%d, error:%s", video.Type, video.TorrentName, video.TorrentURL, count, err)
					}
					goto RETRY
				}
			} else {
				video.Magnet = magnet
				feedVideos = append(feedVideos, video)
			}
		}(video)
	}
	wg.Wait()
	return feedVideos
}

func (t *Torlock) fetchMagnetDownLoad(ctx context.Context, videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	var wg sync.WaitGroup
	var videos2 []*types.FeedVideoBase
	for _, video := range videos {
		wg.Add(1)
		//nolint:noctx
		go func(v *types.FeedVideoBase) {
			defer wg.Done()
			var count int
		RETRY:
			resp, err := t.HTTPClientDynamic(ctx).Get(v.TorrentURL)
			if err == nil {
				defer resp.Body.Close()
				var buf buffer.Buffer
				_, _ = io.Copy(&buf, resp.Body)
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf.Bytes()))
				if err != nil {
					log.WithCtx(ctx).Errorf("torlock.%s %#v url: %s  content：%s  error:%s", v.Type, v.TorrentName, v.TorrentURL, buf.String(), err)
					return
				}
				val, exists := doc.Find("body > article > div:nth-child(6) > div > div:nth-child(2) > a").Attr("href")
				if exists {
					v.Magnet = val
					videos2 = append(videos2, v)
					log.WithCtx(ctx).Debugf("Add: torlock.%s   %#v", v.Type, v)
				}
			}

			if count < 3 {
				count++
				if count == 2 {
					log.WithCtx(ctx).Debugf("RETRY: torlock.%s %s http request url is %s , retry count: %d", v.Type, v.TorrentName, v.TorrentURL, count)
				}
				goto RETRY
			}
		}(video)
	}
	wg.Wait()
	return videos2
}
