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
	"github.com/pkg/errors"

	"github.com/youcd/toolkit/log"

	"github.com/chromedp/chromedp"
)

type TgxWeb struct {
	web string
	BaseFeeder
}

func NewTgxWeb(scheduling, siteURL string) *TgxWeb {
	parse, _ := url.Parse(siteURL)
	return &TgxWeb{
		web: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
		BaseFeeder: BaseFeeder{
			web:        "tgx_web",
			url:        siteURL,
			scheduling: scheduling,
		},
	}
}

func (t *TgxWeb) Crawler() ([]*types.FeedVideo, error) {
	timeCtx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()
	newContext, _ := chromedp.NewContext(timeCtx)
	log.Debug(t.url)
	var htmlStr string
	err := chromedp.Run(newContext,
		chromedp.Navigate(t.url),
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
		return nil, errors.WithMessage(err, "goquery.NewDocumentFromReader goquery")
	}

	videosTemp := make([]*types.FeedVideo, 0)
	doc.Find(".tgxtablerow.txlight").Each(func(_ int, s *goquery.Selection) {
		// 类别
		typStr := s.Find(".tgxtablecell.shrink.rounded.txlight").Text()
		var typ string
		switch {
		case strings.Contains(strings.ToLower(typStr), "tv"):
			typ = "tv"
		case strings.Contains(strings.ToLower(typStr), "movies"):
			typ = "movie"
		default:
			log.Warn("typStr is empty", typStr, s.Text())
			return
		}

		//  名字
		torrentName := s.Find(".tgxtablecell.clickable-row.click.textshadow.rounded.txlight").Text()
		name, _, year, err := torrentName2info(torrentName)
		if err != nil {
			log.Warnf("torrentName2info err: %s", err)
			return
		}

		u, ok := s.Find("#click > div > a.txlight").Attr("href")
		if !ok {
			return
		}

		video := &types.FeedVideo{
			Name:        name,
			TorrentName: name,
			TorrentURL:  t.web + u,
			Magnet:      "",
			Year:        year,
			Type:        typ,
			RowData:     sql.NullString{},
			Web:         t.web,
			DoubanID:    "",
		}
		videosTemp = append(videosTemp, video)
	})

	videos := make([]*types.FeedVideo, 0)
	var wg sync.WaitGroup
	for _, video := range videosTemp {
		wg.Add(1)
		go func(v *types.FeedVideo) {
			defer wg.Done()
			magnet, match := t.fetchMagnet(video.TorrentURL)
			if !match {
				log.Warnf("magnet is empty: %s", video.TorrentURL)
				return
			}
			video.Magnet = magnet
			videos = append(videos, video)
		}(video)
	}
	wg.Wait()
	return videos, nil
}

func (t *TgxWeb) fetchMagnet(urlStr string) (string, bool) {
	timeCtx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	newContext, _ := chromedp.NewContext(timeCtx)
	defer cancel()
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
