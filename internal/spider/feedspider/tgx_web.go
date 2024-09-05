package feedspider

import (
	"context"
	"database/sql"
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/youcd/toolkit/log"
)

type TgxWeb struct {
	webHost string
	BaseFeeder
}

func NewTgxWeb(scheduling, siteURL string, useIPProxy bool) *TgxWeb {
	parse, _ := url.Parse(siteURL)
	return &TgxWeb{
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
		BaseFeeder: BaseFeeder{
			web:      "tgx_web",
			BaseFeed: types.BaseFeed{Url: siteURL, Scheduling: scheduling, UseIPProxy: useIPProxy},
		},
	}
}

func (t *TgxWeb) Crawler() ([]*types.FeedVideoBase, error) {
	//nolint:govet
	timeCtx, _ := context.WithTimeout(context.TODO(), 300*time.Second)
	// allocator, _ := chromedp.NewRemoteAllocator(timeCtx, "http://127.0.0.1:9222/")
	newContext, _ := chromedp.NewContext(timeCtx, chromedp.WithLogf(log.Infof))
	defer newContext.Done()
	log.Debug(t.Url)
	var htmlStr string
	err := chromedp.Run(newContext,
		chromedp.Navigate(t.Url),
		chromedp.WaitVisible(`.tgxtable`),
		chromedp.Sleep(time.Second*10),
		chromedp.InnerHTML(`body`, &htmlStr),
	)

	if err != nil {
		log.Info(err)
		return nil, fmt.Errorf("采集失败：%w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("goquery.NewDocumentFromReader goquery, err:%w", err)
	}

	videosTemp := make([]*types.FeedVideoBase, 0)
	doc.Find(".tgxtablerow.txlight").Each(func(_ int, s *goquery.Selection) {
		// 类别
		typStr := s.Find(".tgxtablecell.shrink.rounded.txlight").Text()
		var typ string
		switch {
		case strings.Contains(strings.ToLower(typStr), "tv"):
			typ = "tv"
		case strings.Contains(strings.ToLower(typStr), "movie"):
			typ = "movie"
		default:
			log.Warn("typStr is empty", typStr, s.Text())
			return
		}

		//  名字
		torrentName := s.Find(".tgxtablecell.clickable-row.click.textshadow.rounded.txlight").Text()
		u, ok := s.Find("#click > div > a.txlight").Attr("href")
		if !ok {
			return
		}

		video := &types.FeedVideoBase{
			TorrentName: torrentName,
			TorrentURL:  t.webHost + u,
			Magnet:      "",
			Type:        typ,
			RowData:     sql.NullString{},
			Web:         t.web,
		}
		videosTemp = append(videosTemp, video)
	})

	videos := make([]*types.FeedVideoBase, 0)
	var wg sync.WaitGroup
	for _, video := range videosTemp {
		wg.Add(1)
		go func(v *types.FeedVideoBase) {
			defer wg.Done()
			magnet, match := t.fetchMagnet(v.TorrentURL)
			if !match {
				log.Warnf("magnet is empty: %s", v.TorrentURL)
				return
			}
			v.Magnet = magnet
			videos = append(videos, v)
		}(video)
	}
	wg.Wait()
	return videos, nil
}

func (t *TgxWeb) fetchMagnet(urlStr string) (string, bool) {
	timeCtx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	newContext, _ := chromedp.NewContext(timeCtx, chromedp.WithLogf(log.Infof))
	chromedp.Stop()
	defer cancel()
	// defer newContext.Done()
	var magnet string
	var match bool

	err := chromedp.Run(newContext,
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(`#covercell > center > a.btn.btn-danger.lift.txlight`),
		chromedp.Sleep(time.Second*20), // #covercell > center > a.btn.btn-danger.lift.txlight
		chromedp.AttributeValue(`#covercell > center > a.btn.btn-danger.lift.txlight`, "href", &magnet, &match),
	)

	if err != nil {
		log.Error(err)
		match = false
	}
	return magnet, match
}
