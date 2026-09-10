package feedspider

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
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
	fp := t.FeedParserUserAgent(ctx, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

	fd, err := fp.ParseURL(t.Url)
	if err != nil {
		return nil, fmt.Errorf("err:%s, err:%w", err, ErrFeedParseURL)
	}

	videos := t.parseFeedItems(fd.Items)
	log.WithCtx(ctx).Infof("%s parsed feed items: %d", t.typ.String(), len(videos))

	videosWithMagnet := t.fetchMagnetDownLoad(ctx, videos)
	log.WithCtx(ctx).Infof("%s MagnetDownLoad Url Count: %d", t.typ.String(), len(videosWithMagnet))

	result := t.fetchMagnet(ctx, videosWithMagnet)
	log.WithCtx(ctx).Infof("%s Magnet Url Count: %d", t.typ.String(), len(result))

	return result, nil
}

// parseFeedItems 解析 feed items 并转换为 FeedVideoBase 切片
func (t *Torlock) parseFeedItems(items []*gofeed.Item) []*types.FeedVideoBase {
	videos := make([]*types.FeedVideoBase, 0, len(items))
	for _, v := range items {
		//nolint:errchkjson
		jsonBytes, _ := json.Marshal(v)
		videos = append(videos, &types.FeedVideoBase{
			Web:         t.web,
			TorrentName: v.Title,
			TorrentURL:  v.Link,
			Type:        t.typ.String(),
			RowData:     sql.NullString{String: string(jsonBytes)},
		})
	}
	return videos
}

func (t *Torlock) fetchMagnet(ctx context.Context, videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		feedVideos []*types.FeedVideoBase
	)

	for _, video := range videos {
		wg.Add(1)
		go func(video *types.FeedVideoBase) {
			defer wg.Done()

			if strings.HasPrefix(video.Magnet, "magnet:?x") {
				mu.Lock()
				feedVideos = append(feedVideos, video)
				mu.Unlock()
				return
			}

			// 使用 for 循环代替 goto 进行重试
			var magnet string
			var err error
			for count := 0; count < 3; count++ {
				magnet, err = magnetconvert.FetchMagnetWithHTTPClient(ctx, video.Magnet, t.HTTPClientDynamic(ctx))
				if err == nil {
					break
				}
				if count == 1 {
					log.WithCtx(ctx).Warnf("RETRY: torlock.%s %s http request url is %s retry count:%d, error:%s",
						video.Type, video.TorrentName, video.TorrentURL, count+1, err)
				}
			}

			if err != nil {
				return
			}

			video.Magnet = magnet
			log.WithCtx(ctx).Debugf("Add: torlock.%s   %#v", video.Type, video)

			mu.Lock()
			feedVideos = append(feedVideos, video)
			mu.Unlock()
		}(video)
	}
	wg.Wait()

	return feedVideos
}

func (t *Torlock) fetchMagnetDownLoad(ctx context.Context, videos []*types.FeedVideoBase) []*types.FeedVideoBase {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []*types.FeedVideoBase
	)

	for _, video := range videos {
		wg.Add(1)
		//nolint:noctx
		go func(v *types.FeedVideoBase) {
			defer wg.Done()

			// 使用 for 循环代替 goto 进行重试
			for count := 0; count < 3; count++ {
				resp, err := t.HTTPClientDynamic(ctx).Get(v.TorrentURL)
				if err != nil {
					if count == 1 {
						log.WithCtx(ctx).Debugf("RETRY: torlock.%s %s http request url is %s , retry count: %d",
							v.Type, v.TorrentName, v.TorrentURL, count+1)
					}
					continue
				}

				var buf buffer.Buffer
				_, _ = io.Copy(&buf, resp.Body)
				resp.Body.Close() // 立即关闭，而不是 defer

				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf.Bytes()))
				if err != nil {
					log.WithCtx(ctx).Errorf("torlock.%s %#v url: %s  content：%s  error:%s",
						v.Type, v.TorrentName, v.TorrentURL, buf.String(), err)
					return
				}

				val, exists := doc.Find("body > article > div:nth-child(6) > div > div:nth-child(2) > a").Attr("href")
				if exists {
					v.Magnet = val
					mu.Lock()
					results = append(results, v)
					mu.Unlock()
				}
				return
			}
		}(video)
	}
	wg.Wait()

	return results
}
