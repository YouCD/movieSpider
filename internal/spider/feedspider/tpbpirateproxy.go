package feedspider

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

//nolint:gosimple,gocritic
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
	//nolint:prealloc
	var videosA []*types.FeedVideo
	for _, v := range fd.Items {
		// 片名
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		// 片名处理
		var name, year string

		switch strings.ToLower(v.Categories[0]) {
		case "tv":
			compileRegex := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(matchArr) == 0 {
				continue
			}
			name = matchArr[1]
		case "movies":
			compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			if len(matchArr) == 0 {
				name = torrentName
			} else {
				name = matchArr[1]
				year = matchArr[2]
			}
		default:
			name = torrentName
		}

		fVideo := new(types.FeedVideo)
		fVideo.Name = fVideo.FormatName(name)
		fVideo.Year = year

		fVideo.Web = g.web
		fVideo.Magnet = v.Link
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)

		fVideo.TorrentURL = v.GUID

		// 处理 资源类型 是 电影 还是电视剧
		typ := strings.ToLower(v.Categories[0])
		if strings.Contains(typ, "tv shows") {
			fVideo.Type = "tv"
		} else if strings.Contains(typ, "movies") {
			fVideo.Type = "movies"
		} else {
			fVideo.Type = typ
		}
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		//nolint:exhaustruct
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
			//if len(matchArr) == 0 {
			//	v.Name = v.Name
			//} else {
			//	v.Name = matchArr[1]
			//	v.Year = strings.ReplaceAll(strings.ReplaceAll(matchArr[2], "(", ""), ")", "")
			//}

			if len(matchArr) > 0 {
				v.Name = matchArr[1]
				v.Year = strings.ReplaceAll(strings.ReplaceAll(matchArr[2], "(", ""), ")", "")
			}
		}
		videos = append(videos, v)
	}
	//nolint:nakedret
	return
}

func (g *tpbpirateproxy) Run(ch chan *types.FeedVideo) {
	if g.scheduling == "" {
		log.Error("tpbpirateproxy Scheduling is null")
		os.Exit(1)
	}
	log.Infof("tpbpirateproxy Scheduling is: [%s]", g.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(g.scheduling, func() {
		videos, err := g.Crawler()
		if err != nil {
			log.Error(err)
			return
		}
		// model.ProxySaveVideo2DB(videos...)
		for _, video := range videos {
			ch <- video
		}
	})
	c.Start()
}
