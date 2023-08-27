package feedSpider

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	urlMagnetdl      = "https://www.magnetdl.com"
	urlMagnetdlTV    = "https://www.magnetdl.com/download/tv/"
	urlMagnetdlMovie = "https://www.magnetdl.com/download/movies/"
)

type magnetdl struct {
	typ        types.VideoType
	web        string
	scheduling string
}

func (m *magnetdl) Crawler() (Videos []*types.FeedVideo, err error) {
	c := httpClient.NewHttpClient()

	switch m.typ {
	case types.VideoTypeMovie:
		req, err := http.NewRequest("GET", urlMagnetdlMovie, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.Do(req)
		if err != nil {
			return nil, errors.WithMessage(err, "getMovies resp")
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, errors.WithMessage(err, "getMovies goquery")
		}
		//#content > div.fill-table > table > tbody > tr:nth-child(1) > td.n > a
		selector := "#content > div.fill-table > table > tbody > tr"
		compileRegex := regexp.MustCompile("(.*)\\.([0-9][0-9][0-9][0-9])\\.")
		//#content > div.fill-table > table > tbody > tr:nth-child(1) > td.n > a
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {

			// 片名
			name := strings.ReplaceAll(s.Text(), " ", ".")
			// magnet 链接
			var magnet, torrentUrl string
			val, exists := s.Find("td>a").Attr("href")
			if exists {
				magnet = val
			} else {
				return
			}

			val, exists = s.Find("td.n > a").Attr("href")
			if exists {
				torrentUrl = urlMagnetdl + val
			}

			fVideo := new(types.FeedVideo)
			matchArr := compileRegex.FindStringSubmatch(name)
			if len(matchArr) > 0 {
				fVideo.Name = fVideo.FormatName(matchArr[1])
			} else {
				fVideo.Name = fVideo.FormatName(name)
			}

			fVideo.Type = "movie"
			fVideo.Magnet = magnet
			fVideo.Web = m.web
			fVideo.TorrentUrl = torrentUrl
			fVideo.TorrentName = fVideo.Name
			Videos = append(Videos, fVideo)
		})
	case types.VideoTypeTV:
		req, err := http.NewRequest("GET", urlMagnetdlTV, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.Do(req)
		if err != nil {
			return nil, errors.WithMessage(err, "getMovies resp")
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, errors.WithMessage(err, "getMovies goquery")
		}
		//#content > div.fill-table > table > tbody > tr:nth-child(1) > td.n > a
		selector := "#content > div.fill-table > table > tbody > tr"
		compileRegex := regexp.MustCompile("(.*)\\.[sS][0-9][0-9]|[Ee][0-9][0-9]?\\.")
		//#content > div.fill-table > table > tbody > tr:nth-child(1) > td.n > a
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {

			// 片名
			name := strings.ReplaceAll(s.Text(), " ", ".")

			// magnet 链接
			var magnet string
			val, exists := s.Find("td>a").Attr("href")
			if exists {
				if !strings.Contains(strings.ToLower(val), "magnet") {
					return
				}
				magnet = val
			} else {
				return
			}

			var torrentUrl string
			val, exists = s.Find("td.n > a").Attr("href")
			if exists {
				torrentUrl = urlMagnetdl + val
			}

			fVideo := new(types.FeedVideo)
			matchArr := compileRegex.FindStringSubmatch(name)
			if len(matchArr) > 0 {
				fVideo.Name = fVideo.FormatName(matchArr[1])
			} else {
				fVideo.Name = fVideo.FormatName(name)
			}

			fVideo.Type = "tv"
			fVideo.Magnet = magnet
			fVideo.Web = m.web
			fVideo.TorrentUrl = torrentUrl
			fVideo.TorrentName = fVideo.Name
			Videos = append(Videos, fVideo)
		})
	}
	return
}

func (m *magnetdl) Run(ch chan *types.FeedVideo) {
	if m.scheduling == "" {
		log.Error("MAGNETDL: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("MAGNETDL: Scheduling is: [%s]", m.scheduling)
	c := cron.New()
	c.AddFunc(m.scheduling, func() {
		videos, err := m.Crawler()
		if err != nil {
			log.Error(err)
			return
		}
		for _, video := range videos {
			ch <- video
		}
		//model.ProxySaveVideo2DB(videos...)
	})
	c.Start()

}
