package feedspider

import (
	"errors"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"regexp"
	"strings"

	"github.com/youcd/toolkit/log"
)

type ThePirateBay struct {
	BaseFeeder
}

var (
	ErrNotMatchTorrentName = errors.New("torrent name not match")
)

func NewThePirateBay() *ThePirateBay {
	return &ThePirateBay{
		BaseFeeder{
			web: "ThePirateBay",
			BaseFeed: types.BaseFeed{
				Url:        config.Config.Feed.ThePirateBay.Url,
				Scheduling: config.Config.Feed.ThePirateBay.Scheduling,
				UseIPProxy: config.Config.Feed.ThePirateBay.UseIPProxy,
			},
		},
	}
}
func (t *ThePirateBay) Crawler() ([]*types.FeedVideoBase, error) {
	fd, err := t.FeedParser().ParseURL(t.Url)
	if err != nil {
		return nil, ErrFeedParseURL
	}
	log.Debugf("%s Data: %#v", strings.ToUpper(t.web), fd.String())

	videos := make([]*types.FeedVideoBase, 0)
	for _, v := range fd.Items {
		if len(v.Categories) < 1 {
			continue
		}
		video := new(types.FeedVideoBase)
		if strings.Contains(strings.ToLower(v.Categories[0]), "movie") {
			video.Type = types.VideoTypeMovie.String()
		}
		if strings.Contains(strings.ToLower(v.Categories[0]), "tv") {
			video.Type = types.VideoTypeTV.String()
		}

		video.TorrentName = v.Title
		video.Magnet = v.Link
		video.TorrentURL = v.GUID
		video.Web = t.web
		videos = append(videos, video)
	}
	return videos, nil
	// ctx, cancel := chromedp.NewContext(
	// 	context.Background(),
	// 	chromedp.WithLogf(log.Infof),
	// 	// chromedp.WithDebugf(log.Printf),
	// )
	// defer cancel()
	//
	// ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()
	//
	// var html string
	// if err := chromedp.Run(ctx,
	//	chromedp.Navigate(t.Url),
	//	chromedp.WaitVisible(`body > main`),
	//	chromedp.Click(`#f_1080p`, chromedp.NodeVisible),
	//	chromedp.Click(`#f_2160p`, chromedp.NodeVisible),
	//	chromedp.InnerHTML(`#torrents`, &html),
	// ); err != nil {
	// 	return nil, fmt.Errorf("chromedp.Run err: %w", err)
	// }
	//
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	// if err != nil {
	// 	return nil, fmt.Errorf("goquery.NewDocumentFromReader err: %w", err)
	// }
	//
	// var videos []*types.FeedVideo
	// doc.Find("#st").Each(func(_ int, s *goquery.Selection) {
	// 	//  名字
	// 	torrentName := s.Find("span.list-item.item-name.item-title").Text()
	// 	if torrentName == "" {
	//		log.Warn("text is empty")
	//		return
	//	}
	//	// 连接地址
	//	magnet, exists := s.Find("span.item-icons > a").Attr("href")
	//	if !exists {
	//		log.Warn("magnet is empty")
	//		return
	//	}
	//
	//	// 类型
	//	// 过滤掉 其他类型的种子
	//	typStr := s.Find("span.list-item.item-type > a:nth-child(2)").Text()
	//	var typ string
	//	switch {
	//	case strings.Contains(strings.ToLower(typStr), "tv-shows"):
	//		typ = "tv"
	//	case strings.Contains(strings.ToLower(typStr), "movies"):
	//		typ = "movie"
	//	default:
	//		log.Warn("typStr is empty: ", typStr)
	//		return
	//	}
	//
	//	name, _, year, err := torrentName2info(torrentName)
	//	if err != nil {
	//		log.Warnf("torrentName2info err: %s", err)
	//		return
	//	}
	//	video := &types.FeedVideo{
	//		Name:        name,
	//		TorrentName: name,
	//		TorrentURL:  "",
	//		Magnet:      magnet,
	//		Year:        year,
	//		Type:        typ,
	//		RowData:     sql.NullString{},
	//		Web:         t.web,
	//		DoubanID:    "",
	//	}
	// 	video.Name = video.FormatName(video.Name)
	// 	videos = append(videos, video)
	// })
	// return videos, nil
}

func torrentName2info(torrentName string) (string, string, string, error) {
	// 去除空格
	newTorrentName := strings.ReplaceAll(torrentName, " ", ".")
	// 去除 []
	newTorrentName = strings.ReplaceAll(newTorrentName, "[", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, "]", "")

	var name, resolution, year string
	compileRegex := regexp.MustCompile(`(.*)\.(\d{4}[p|P])\.`)
	matchArr := compileRegex.FindStringSubmatch(newTorrentName)
	if len(matchArr) < 3 {
		return "", "", "", fmt.Errorf("torrentName:%s,第一次匹配失败:%w", torrentName, ErrNotMatchTorrentName)
	}
	resolution = matchArr[2]
	if len(matchArr) == 0 {
		tvReg := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9][eE][0-9][0-9])`)
		TVNameArr := tvReg.FindStringSubmatch(torrentName)
		// 如果 正则匹配过后 没有结果直接 过滤掉
		if len(TVNameArr) == 0 {
			return "", "", "", fmt.Errorf("第二次 匹配失败:%w", ErrNotMatchTorrentName)
		}
		name = TVNameArr[1]
	} else {
		name = matchArr[1]
	}
	compileYearRegex := regexp.MustCompile(`(\d{4})`)
	matchArrYear := compileYearRegex.FindStringSubmatch(newTorrentName)
	if len(matchArr) < 2 {
		return "", "", "", fmt.Errorf("第一次匹配失败:%w", ErrNotMatchTorrentName)
	}
	year = matchArrYear[0]

	return name, resolution, year, nil
}
