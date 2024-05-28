package feedspider

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"net/http"
	"regexp"
	"strings"
)

const (
	urlBaseMagnetdl        = "https://www.magnetdl.com"
	urlRssURITVMagnetdl    = "download/tv/"
	urlRssURIMovieMagnetdl = "download/movies/"
)

type Magnetdl struct {
	typ types.VideoType
	BaseFeeder
}

func NewMagnetdl(scheduling string, typ types.VideoType, mirrorSite string) *Magnetdl {
	resourceType := typ
	urlBase := urlBaseMagnetdl
	if mirrorSite != "" {
		urlBase = mirrorSite
	}

	url := fmt.Sprintf("%s/%s", urlBase, urlRssURITVMagnetdl)
	if resourceType == types.VideoTypeMovie {
		url = fmt.Sprintf("%s/%s", urlBase, urlRssURIMovieMagnetdl)
	}
	return &Magnetdl{
		resourceType,
		BaseFeeder{
			web:        "magnetdl",
			url:        url,
			scheduling: scheduling,
		},
	}
}

//nolint:gosimple,gocognit,gocritic
func (m *Magnetdl) Crawler() (Videos []*types.FeedVideo, err error) {
	c := httpclient.NewHTTPClient()
	//nolint:exhaustive
	switch m.typ {
	case types.VideoTypeMovie:
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, m.url, nil)
		if err != nil {
			return nil, errors.WithMessage(err, "magnetdl req")
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
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			// 片名
			name := strings.ReplaceAll(s.Text(), " ", ".")
			// magnet 链接
			var magnet, torrentURL string
			val, exists := s.Find("td>a").Attr("href")
			if exists {
				magnet = val
			} else {
				return
			}

			val, exists = s.Find("td.n > a").Attr("href")
			if exists {
				torrentURL = urlBaseMagnetdl + val
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
			fVideo.TorrentURL = torrentURL
			fVideo.TorrentName = fVideo.Name
			Videos = append(Videos, fVideo)
		})
	case types.VideoTypeTV:
		log.Infof("%s working, type:%s ,url: %s", strings.ToUpper(m.web), m.typ.String(), m.url)
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, m.url, nil)
		if err != nil {
			return nil, errors.WithMessage(err, "getMovies newRequest")
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
		doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
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

			var torrentURL string
			val, exists = s.Find("td.n > a").Attr("href")
			if exists {
				torrentURL = urlBaseMagnetdl + val
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
			fVideo.TorrentURL = torrentURL
			fVideo.TorrentName = fVideo.Name
			Videos = append(Videos, fVideo)
		})
	}
	//nolint:nakedret
	return
}
