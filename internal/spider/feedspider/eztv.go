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

const urlEztv = "https://eztv.re/ezrss.xml"

type eztv struct {
	scheduling string
	url        string
	web        string
}

//nolint:gosimple
func (f *eztv) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(f.url)
	if fd == nil {
		return nil, errors.New("EZTV: 没有feed数据")
	}
	log.Debugf("EZTV Config: %#v", fd)
	log.Debugf("EZTV Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("EZTV: 没有feed数据")
	}
	for _, v := range fd.Items {
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		var name, year string
		compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
		matchArr := compileRegex.FindStringSubmatch(torrentName)
		if len(matchArr) == 0 {
			tvReg := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			TVNameArr := tvReg.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(TVNameArr) == 0 {
				continue
			}
			name = TVNameArr[1]
		} else {
			year = matchArr[2]
			name = matchArr[1]
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = f.web
		fVideo.Year = year

		// 片名
		fVideo.Name = fVideo.FormatName(name)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)
		fVideo.TorrentURL = v.Link
		fVideo.Magnet = v.Extensions["torrent"]["magnetURI"][0].Value
		//nolint:errchkjson
		bytes, _ := json.Marshal(v)
		fVideo.Type = strings.ToLower(v.Categories[0])
		//nolint:exhaustruct
		fVideo.RowData = sql.NullString{String: string(bytes)}

		videos = append(videos, fVideo)
	}
	//nolint:nakedret
	return
}
func (f *eztv) Run(ch chan *types.FeedVideo) {
	if f.scheduling == "" {
		log.Error("EZTV Scheduling is null")
		os.Exit(1)
	}
	log.Infof("EZTV Scheduling is: [%s]", f.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(f.scheduling, func() {
		videos, err := f.Crawler()
		if err != nil {
			log.Error(err)
			return
		}
		for _, video := range videos {
			ch <- video
		}
		// model.ProxySaveVideo2DB(videos...)
	})
	c.Start()
}
