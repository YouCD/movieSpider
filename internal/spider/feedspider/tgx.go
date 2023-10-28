package feedspider

import (
	"database/sql"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/log"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"os"
	"regexp"
	"strings"
	"sync"
)

const urlTgx = "https://tgx.rs/rss"

type tgx struct {
	scheduling string
	url        string
	web        string
}

//nolint:gosimple,gocritic,,rowserrcheck,ineffassign
func (t *tgx) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(t.url)
	if fd == nil {
		return nil, errors.New("TGx: 没有feed数据")
	}
	log.Debugf("TGx Config: %#v", fd)
	log.Debugf("TGx Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("TGx: 没有feed数据")
	}
	//nolint:prealloc,ineffassign
	var videos1 []*types.FeedVideo
	for _, v := range fd.Items {
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		var name, year, typ string

		compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
		matchArr := compileRegex.FindStringSubmatch(torrentName)
		if len(matchArr) == 0 {
			tvReg := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			TVNameArr := tvReg.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(TVNameArr) == 0 {
				continue
			}
			//nolint:wastedassign
			name = TVNameArr[1]
		}
		year = matchArr[2]
		name = matchArr[1]

		// 过滤掉 其他类型的种子
		if strings.HasPrefix(strings.ToLower(v.Categories[0]), "tv :") {
			typ = "tv"
		} else if strings.HasPrefix(strings.ToLower(v.Categories[0]), "movies :") {
			typ = "movie"
		} else {
			continue
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = t.web
		fVideo.Year = year

		// 片名
		fVideo.Name = fVideo.FormatName(name)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)
		fVideo.Type = typ
		fVideo.TorrentURL = v.Link
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		//nolint:exhaustruct
		fVideo.RowData = sql.NullString{String: string(bytes)}
		videos1 = append(videos1, fVideo)
	}
	var wg sync.WaitGroup
	for _, video := range videos1 {
		wg.Add(1)
		magnet, err := magnetconvert.FetchMagnet(video.TorrentURL)
		if err != nil {
			wg.Done()
			return nil, errors.WithMessage(err, "FetchMagnet")
		}
		video.Magnet = magnet
		videos = append(videos, video)
		wg.Done()
	}
	wg.Wait()
	//nolint:nakedret
	return
}

func (t *tgx) Run(ch chan *types.FeedVideo) {
	if t.scheduling == "" {
		log.Error("TGx Scheduling is null")
		os.Exit(1)
	}
	log.Infof("TGx Scheduling is: [%s]", t.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(t.scheduling, func() {
		videos, err := t.Crawler()
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
