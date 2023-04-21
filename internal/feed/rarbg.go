package feed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/ipProxy"
	"movieSpider/internal/log"
	types2 "movieSpider/internal/types"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	urlRarbgMovie = "http://rarbg.to/rssdd.php?categories=14;15;16;17;21;22;42;44;45;46;47;48"
	urlRarbgTV    = "http://rarbg.to/rssdd.php?categories=18;19;41"
)

type rarbg struct {
	typ        types2.Resource
	web        string
	scheduling string
	httpClient *http.Client
}

func (r *rarbg) Crawler() (Videos []*types2.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fp.Client = r.httpClient
	fp.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	if r.typ == types2.ResourceMovie {
		fd, err := fp.ParseURL(urlRarbgMovie)
		if err != nil {
			log.Error(err)
		}
		if fd == nil {
			return nil, errors.New(fmt.Sprintf("RARBG.%s feed is nil.", r.typ.Typ()))
		}
		log.Debugf("RARBG.movie Data: %#v", fd.String())
		if len(fd.Items) == 0 {
			return nil, errors.New(fmt.Sprintf("RARBG.movie: 没有feed数据."))
		}
		compileRegex := regexp.MustCompile("(.*)\\.([0-9][0-9][0-9][0-9])\\.")
		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")
			ok := excludeVideo(name)
			if ok {
				continue
			}

			var fVideo types2.FeedVideo
			fVideo.Web = r.web
			fVideo.TorrentName = name
			fVideo.Magnet = v.Link
			fVideo.Type = "movie"

			// 原始数据
			bytes, _ := json.Marshal(v)
			fVideo.RowData = sql.NullString{String: string(bytes)}

			// 片名
			matchArr := compileRegex.FindStringSubmatch(name)
			if len(matchArr) > 0 {
				fVideo.Name = fVideo.FormatName(matchArr[1])
			} else {
				fVideo.Name = fVideo.FormatName(name)
			}
			// 年份
			if len(matchArr) > 0 {
				fVideo.Year = matchArr[2]
			}
			Videos = append(Videos, &fVideo)
		}
	}
	if r.typ == types2.ResourceTV {
		fd, err := fp.ParseURL(urlRarbgTV)
		if err != nil {
			log.Error("RARBG.tv:", err)
		}
		if fd == nil {
			return nil, errors.New("RARBG.tv: 没有feed数据")
		}
		log.Debugf("RARBG.tv Data: %#v", fd.String())
		if len(fd.Items) == 0 {
			return nil, errors.New("RARBG.tv: 没有feed数据")
		}
		compileRegex := regexp.MustCompile("(.*)\\.[sS][0-9][0-9]|[Ee][0-9][0-9]?\\.")
		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")
			ok := excludeVideo(name)
			if ok {
				continue
			}

			matchArr := compileRegex.FindStringSubmatch(name)

			var fVideo types2.FeedVideo
			fVideo.TorrentName = fVideo.FormatName(name)
			fVideo.Magnet = v.Link
			fVideo.Type = "tv"
			// 原始数据
			bytes, _ := json.Marshal(v)
			fVideo.RowData = sql.NullString{String: string(bytes)}
			fVideo.Web = r.web
			// 片名
			if len(matchArr) > 0 {
				fVideo.Name = fVideo.FormatName(matchArr[1])
			} else {
				fVideo.Name = fVideo.FormatName(name)
			}
			Videos = append(Videos, &fVideo)
		}
	}
	return
}
func (r *rarbg) Run() {
	if r.scheduling == "" {
		log.Errorf("RARBG %s: Scheduling is null", r.typ.Typ())
		os.Exit(1)
	}
	log.Infof("RARBG %s: Scheduling is: [%s]", r.typ.Typ(), r.scheduling)
	c := cron.New()
	c.AddFunc(r.scheduling, func() {
		videos, err := r.Crawler()
		if err != nil {
			//for {
			//	r.switchClient()
			//	videos, err = r.Crawler()
			//	if err != nil {
			//		log.Error(err)
			//		return
			//	}
			//	if len(videos) == 0 {
			//		continue
			//	} else {
			//		r.proxySaveVideo2DB(videos)
			//		break
			//	}
			//
			//}
			r.switchClient()
			videos, err = r.Crawler()
			if err != nil {
				log.Error(err)
				return
			}
			proxySaveVideo2DB(videos...)
		}
		proxySaveVideo2DB(videos...)
	})
	c.Start()

}

func (r *rarbg) useProxyClient() {
	proxyStr := ipProxy.FetchProxy("")
	if proxyStr == "" {
		log.Error("useProxyClient: proxy is null")
		return
	}

	proxyUrl, err := url.Parse(proxyStr)
	if err != nil {
		log.Error(err)
		return
	}
	if proxyUrl != nil {
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		log.Errorf("useProxyClient: use proxy %#v", proxyUrl.String())
		httpClient := &http.Client{Transport: transport, Timeout: time.Minute * 5}

		r.httpClient = httpClient
	}
	return
}
func (r *rarbg) switchClient() {
	if r.httpClient.Transport == nil {

		proxyStr := ipProxy.FetchProxy("")
		if proxyStr == "" {
			log.Infof("RARBG.%s: proxy is null.", r.typ.Typ())
			return
		}
		proxyUrl, err := url.Parse(proxyStr)
		if err != nil {
			log.Error(err)
		}
		if proxyUrl != nil {
			proxy := http.ProxyURL(proxyUrl)
			transport := &http.Transport{Proxy: proxy}
			httpClient := &http.Client{Transport: transport, Timeout: time.Minute * 5}
			r.httpClient = httpClient
			log.Infof("RARBG.%s: 添加代理. proxy: %s", r.typ.Typ(), proxyUrl)
		} else {
			log.Warnf("RARBG.%s: 请添加Global.Proxy.Url配置", r.typ.Typ())
		}

	} else {
		r.httpClient = &http.Client{}
		log.Infof("RARBG.%s: 删除代理.", r.typ.Typ())
	}
}
