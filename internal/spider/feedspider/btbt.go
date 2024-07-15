package feedspider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/youcd/toolkit/log"
)

type Btbt struct {
	BaseFeeder
	webHost string
}

func NewBtbt(scheduling string, siteURL string) *Btbt {
	parse, _ := url.Parse(siteURL)
	return &Btbt{
		BaseFeeder: BaseFeeder{
			web:        "btbt",
			url:        siteURL,
			scheduling: scheduling,
		},
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}
func (b *Btbt) Crawler() (videos []*types.FeedVideo, err error) {
	timeCtx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()
	newContext, _ := chromedp.NewContext(timeCtx)
	log.Debug(b.url)
	var htmlStr string
	err = chromedp.Run(newContext,
		chromedp.Navigate(b.url),
		chromedp.WaitVisible(`.media.thread.tap.hidden-sm`),
		chromedp.Sleep(time.Second*10),
		chromedp.InnerHTML(`body`, &htmlStr),
	)
	if err != nil {
		return nil, fmt.Errorf("btbt new request,err: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery,err: %w", err)
	}

	// 洗出  下载电影的页面
	var Videos1 []*types.FeedVideo
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

		video := &types.FeedVideo{
			Name:        trim(matchArr[1]),
			TorrentName: torrentName,
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
		go func(video *types.FeedVideo) {
			defer wg.Done()
			magnet, err := fetchMagnet(b.webHost, video.TorrentURL)
			if err != nil {
				log.Errorf("BTBT: 获取磁链错误 err:%s", err)
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

func fetchMagnet(webHost, torrentURL string) (string, error) {
	downloadURL, err := moviePageURL(torrentURL)
	if err != nil {
		return "", fmt.Errorf("BTBT: 获取磁链错误 err:%w", err)
	}
	if downloadURL == "" {
		return "", ErrDownloadURLIsEmpty
	}

	magnet, err := getMagnet(fmt.Sprintf("%s/%s", webHost, downloadURL))
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

func moviePageURL(pageURL string) (url string, err error) {
	resp, err := client(pageURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
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

func getMagnet(url string) (magnet string, err error) {
	resp, err := client(url)
	if err != nil {
		return "", fmt.Errorf("BTBT: 获取磁链错误 err:%w", err)
	}
	defer resp.Body.Close()
	//nolint:wrapcheck
	return magnetconvert.IO2Magnet(resp.Body)
}

func client(url string) (resp *http.Response, err error) {
	c := httpclient.NewHTTPClient()
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client new request,err: %w", err)
	}
	resp, err = c.Do(req)
	return
}
