package feedspider

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	httpClient2 "movieSpider/internal/httpclient"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const urlGlodls = "http://glodls.to/rss.php?cat=1,41"

type glodls struct {
	url        string
	scheduling string
	web        string
	httpClient *http.Client
}

//nolint:gosimple,ineffassign,goconst
func (g *glodls) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fp.Client = g.httpClient
	fd, err := fp.ParseURL(g.url)
	if fd == nil {
		return nil, errors.New("GLODLS: 没有feed数据")
	}
	log.Debugf("GLODLS Config: %#v", fd)
	log.Debugf("GLODLS Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("GLODLS: 没有feed数据")
	}
	//nolint:prealloc
	var videosA []*types.FeedVideo
	for _, v := range fd.Items {
		// 片名
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		// 片名处理
		var name, year string
		if strings.ToLower(v.Categories[0]) == "tv" {
			compileRegex := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(matchArr) == 0 {
				continue
			}
			//nolint:wastedassign
			name = matchArr[1]
		} else if strings.ToLower(v.Categories[0]) == "movies" {
			compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			if len(matchArr) == 0 {
				//nolint:wastedassign
				name = torrentName
			} else {
				//nolint:wastedassign
				name = matchArr[1]
				year = matchArr[2]
			}
		}
		name = torrentName

		fVideo := new(types.FeedVideo)
		fVideo.Name = fVideo.FormatName(name)
		fVideo.Year = year

		fVideo.Web = g.web
		parse, _ := url.Parse(v.Link)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)

		if len(parse.Query()["id"]) == 0 {
			log.Error("没有ID")
		}
		id := parse.Query()["id"][0]
		all := strings.ReplaceAll(v.Title, " ", "-")

		TorrentURL := fmt.Sprintf("http://glodls.to/%s-f-%s.html", strings.ToLower(all), id)

		fVideo.TorrentURL = TorrentURL

		// 处理 资源类型 是 电影 还是电视剧
		typ := strings.ToLower(v.Categories[0])
		if typ == "movies" {
			fVideo.Type = "movie"
		} else {
			fVideo.Type = typ
		}
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		//nolint:exhaustruct
		fVideo.RowData = sql.NullString{String: string(bytes)}

		videosA = append(videosA, fVideo)
	}
	var wg sync.WaitGroup
	for _, v := range videosA {
		wg.Add(1)
		go func(video *types.FeedVideo) {
			defer wg.Done()
			magnet, err := g.fetchMagnet(video.TorrentURL)
			if err != nil {
				err := errors.Unwrap(err)
				log.Error(err)
			}
			if magnet == "" {
				return
			}
			video.Magnet = magnet
			videos = append(videos, video)
		}(v)
	}
	wg.Wait()
	//nolint:nakedret
	return
}

func (g *glodls) fetchMagnet(url string) (magnet string, err error) {
	request, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "GLODLS: 创建新的请求")
	}
	g.httpClient.Timeout = 20 * time.Second
	resp, err := g.httpClient.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "GLODLS: 请求错误")
	}
	if resp == nil {
		return "", errors.New("GLODLS: response is nil")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "GLODLS: Document 查找出错")
	}
	selector := "#downloadbox > table > tbody > tr > td:nth-child(1) > a:nth-child(2)"
	magnet, exists := doc.Find(selector).Attr("href")
	if !exists {
		return "", errors.Wrap(err, "GLODLS: 查找href出错")
	}
	return magnet, nil
}

func (g *glodls) Run(ch chan *types.FeedVideo) {
	if g.scheduling == "" {
		log.Error("GLODLS Scheduling is null")
		os.Exit(1)
	}
	log.Infof("GLODLS Scheduling is: [%s]", g.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(g.scheduling, func() {
		log.Info("GLODLS: is working...")
		for {
		Start:
			videos, err := g.Crawler()
			if err != nil {
				log.Error(err)
			}
			if len(videos) == 0 || videos == nil {
				log.Info("GLODLS: 切换代理")
				g.httpClient = httpClient2.NewProxyHTTPClient("http")
				log.Info("GLODLS: crawler agan...")
				goto Start
			}
			for _, video := range videos {
				ch <- video
			}
			break
		}
	})
	c.Start()
}
