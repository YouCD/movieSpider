package feedspider

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Btbt struct {
	BaseFeeder
	webHost string
}

func NewBtbt() *Btbt {
	parse, _ := url.Parse(config.Config.Feed.BTBT.Url)
	return &Btbt{
		BaseFeeder: BaseFeeder{
			BaseFeed: types.BaseFeed{
				Url:        config.Config.Feed.BTBT.Url,
				Scheduling: config.Config.Feed.BTBT.Scheduling,
				UseIPProxy: config.Config.Feed.BTBT.UseIPProxy,
			},
			web: "btbt",
		},
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}
func (b *Btbt) Crawler() (videos []*types.FeedVideoBase, err error) {
	log.Debug(b.Url)

	resp, err := b.HTTPRequest(b.Url)
	if err != nil {
		return nil, fmt.Errorf("btbt new request,url: %s, err: %w", b.Url, err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}

	// 洗出  下载电影的页面
	var Videos1 []*types.FeedVideoBase
	compileRegex := regexp.MustCompile(`\[.*?\]`)

	doc.Find(".media.thread.tap.hidden-sm").Each(func(_ int, s *goquery.Selection) {
		// 跳过 title 为 站务专区
		val, exists := s.Find(".text-secondary.small.hidden-sm").Attr("title")
		if exists && val == "站务专区" {
			return
		}

		// torrentName
		torrentName := s.Find(".text-title").Text()
		// 跳过 包含 excludeWords 的视频
		if ok := tools.ExcludeVideo(torrentName, config.Config.ExcludeWords); ok {
			return
		}
		year := s.Find("a:nth-child(2)").Text()

		href, exists := s.Find(".text-title").Attr("href")
		if !exists {
			return
		}
		// torrentURL 连接地址
		torrentURL := fmt.Sprintf("%s/%s", b.webHost, href)

		matchArr := compileRegex.FindAllString(torrentName, 2)
		if len(matchArr) < 2 {
			return
		}

		video := &types.FeedVideoBase{
			//Name:        trim(matchArr[1]),
			TorrentName: trim(matchArr[1]),
			TorrentURL:  torrentURL,
			//Magnet:      magnet,
			Year:    year,
			Type:    "movie",
			RowData: sql.NullString{},
			Web:     b.web,
		}

		Videos1 = append(Videos1, video)
	})

	// 洗出  下载电影种子下载的页面
	var wg sync.WaitGroup
	for _, v := range Videos1 {
		wg.Add(1)
		go func(video *types.FeedVideoBase) {
			defer wg.Done()
			magnet, err := b.fetchMagnet(b.webHost, video.TorrentURL)
			if err != nil {
				log.Errorf("BTBT: 获取磁链错误 err:%s", errors.Unwrap(err))
				return
			}
			video.Magnet = magnet
			videos = append(videos, video)
		}(v)
	}
	wg.Wait()

	return videos, nil
}

var (
	ErrDownloadURLIsEmpty = errors.New("downloadURL is empty")
)

func (b *Btbt) fetchMagnet(webHost, torrentURL string) (string, error) {
	downloadURL, err := b.moviePageURL(torrentURL)
	if err != nil {
		return "", fmt.Errorf("BTBT: 获取磁链下载连接错误 err:%w", err)
	}
	if downloadURL == "" {
		return "", ErrDownloadURLIsEmpty
	}

	magnet, err := b.getMagnet(fmt.Sprintf("%s/%s", webHost, downloadURL))
	if err != nil {
		return "", fmt.Errorf("BTBT: 获取磁链错误 err:%w", err)
	}
	return magnet, nil
}

func trim(str string) string {
	str = strings.Trim(str, "[")
	str = strings.Trim(str, "]")
	return str
}

func (b *Btbt) moviePageURL(pageURL string) (url string, err error) {
	resp, err := b.HTTPRequest(pageURL)
	if err != nil {
		return "", fmt.Errorf("连接请求: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", fmt.Errorf("moviePageURL: %w", err)
	}

	doc.Find("#body > div > div:nth-child(2) > div > div.card.card-thread > div > div.message.break-all > fieldset > ul > li > a").Each(func(_ int, selection *goquery.Selection) {
		u, ok := selection.Attr("href")
		if !ok {
			return
		}
		url = u
	})
	return url, nil
}

func (b *Btbt) getMagnet(url string) (magnet string, err error) {
	resp, err := b.HTTPRequest(url)
	if err != nil {
		return "", fmt.Errorf("BTBT: 获取磁链错误 err:%w", err)
	}
	//nolint:wrapcheck
	return magnetconvert.IO2Magnet(bytes.NewReader(resp))
}
