package feed

import (
	"database/sql"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"os"
	"regexp"
	"strings"
)

const urlTpbpirateProxy = "https://thepiratebay.party/rss//top100/200"

type tpbpirateproxy struct {
	scheduling string
	web        string
}

func (g *tpbpirateproxy) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(urlTpbpirateProxy)
	if fd == nil {
		return nil, errors.New("tpbpirateproxy: 没有feed数据")
	}
	log.Debugf("Tpbpirateproxy Config: %#v", fd)
	log.Debugf("Tpbpirateproxy Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("Tpbpirateproxy: 没有feed数据")
	}
	var videosA []*types.FeedVideo
	for _, v := range fd.Items {
		// 片名
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		ok := excludeVideo(torrentName)
		if ok {
			continue
		}

		// 片名处理
		var name, year string

		if strings.ToLower(v.Categories[0]) == "tv" {
			compileRegex := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(matchArr) == 0 {
				continue
			}
			name = matchArr[1]
		} else if strings.ToLower(v.Categories[0]) == "movies" {
			compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			if len(matchArr) == 0 {
				name = torrentName
			} else {
				name = matchArr[1]
				year = matchArr[2]
			}

		} else {
			name = torrentName
		}

		fVideo := new(types.FeedVideo)
		fVideo.Name = fVideo.FormatName(name)
		fVideo.Year = year

		fVideo.Web = g.web
		fVideo.Magnet = v.Link
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)

		fVideo.TorrentUrl = v.GUID

		// 处理 资源类型 是 电影 还是电视剧
		typ := strings.ToLower(v.Categories[0])
		if strings.Contains(typ, "tv shows") {
			fVideo.Type = "tv"
		} else if strings.Contains(typ, "movies") {
			fVideo.Type = "movies"
		} else {
			fVideo.Type = typ
		}

		bytes, _ := json.Marshal(v)
		fVideo.RowData = sql.NullString{String: string(bytes)}

		videosA = append(videosA, fVideo)

	}

	tvRegex := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
	moviesRegex := regexp.MustCompile("(.*)\\.(\\(\\d{4}\\))\\.")

	for _, v := range videosA {
		if v.Type == "tv" {
			matchArr := tvRegex.FindStringSubmatch(v.Name)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(matchArr) == 0 {
				continue
			}
			v.Name = matchArr[1]
		}
		if v.Type == "movies" {
			matchArr := moviesRegex.FindStringSubmatch(v.Name)
			if len(matchArr) == 0 {
				v.Name = v.Name
			} else {
				v.Name = matchArr[1]
				v.Year = strings.ReplaceAll(strings.ReplaceAll(matchArr[2], "(", ""), ")", "")
			}
		}
		videos = append(videos, v)
	}
	return
}

//func (g *tpbpirateproxy) fetchMagnet(url string) (magnet string, err error) {
//	request, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil {
//		return "", errors.WithMessage(err, "tpbpirateproxy: 磁链获取错误")
//	}
//	client := httpClient2.NewHttpClient()
//	resp, err := client.Do(request)
//	if err != nil {
//		return "", errors.WithMessage(err, "tpbpirateproxy: 磁链获取错误")
//	}
//	if resp == nil {
//		return "", errors.New("tpbpirateproxy: response is nil")
//	}
//	defer resp.Body.Close()
//
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		return "", errors.WithMessage(err, "tpbpirateproxy: 磁链获取错误")
//	}
//	selector := "#downloadbox > table > tbody > tr > td:nth-child(1) > a:nth-child(2)"
//	magnet, exists := doc.Find(selector).Attr("href")
//	if !exists {
//		return "", errors.WithMessage(err, "tpbpirateproxy: 磁链获取错误")
//	}
//	return magnet, nil
//}
func (g *tpbpirateproxy) Run() {

	if g.scheduling == "" {
		log.Error("tpbpirateproxy Scheduling is null")
		os.Exit(1)
	}
	log.Infof("tpbpirateproxy Scheduling is: [%s]", g.scheduling)
	c := cron.New()
	c.AddFunc(g.scheduling, func() {
		videos, err := g.Crawler()
		if err != nil {
			log.Error(err)
			return
		}
		proxySaveVideo2DB(videos...)
	})
	c.Start()

}
