package feedspider

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/types"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/youcd/toolkit/log"
)

type Web1337x struct {
	typ types.VideoType
	BaseFeeder
	webHost string
}

func NewWeb1337x(scheduling string, resourceType types.VideoType, siteURL string) *Web1337x {
	parse, _ := url.Parse(siteURL)
	return &Web1337x{
		typ: resourceType,
		BaseFeeder: BaseFeeder{
			web:        "1337x",
			url:        siteURL,
			scheduling: scheduling,
		},
		webHost: fmt.Sprintf("%s://%s", parse.Scheme, parse.Host),
	}
}

func (w *Web1337x) Crawler() ([]*types.FeedVideo, error) {
	videosTemp, err := w.crawler()
	if err != nil {
		return nil, err
	}
	//nolint:exhaustive
	return videosTemp, nil
}

func (w *Web1337x) crawler() ([]*types.FeedVideo, error) {
	videosTemp := make([]*types.FeedVideo, 0)
	data, err := fetchHTMLData(w.url)
	if err != nil {
		return nil, fmt.Errorf("fetchHTMLData, err:%w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("getMovies goquery, err:%w", err)
	}

	selector := "body > main > div > div > div.featured-list > div > table > tbody > tr"

	compileRegex := regexp.MustCompile(`/torrent/\d+/(.*)/.*`)
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		var torrentName, torrentURL, magnet, typ string
		// href
		href, exists := s.Find("td.coll-1.name > a:nth-child(2)").Attr("href")
		if !exists {
			return
		}
		torrentURL = fmt.Sprintf("%s%s", w.webHost, href)

		matchArr := compileRegex.FindStringSubmatch(href)
		if len(matchArr) <= 1 {
			return
		}
		// 种子名
		torrentName = strings.ReplaceAll(matchArr[1], "-", ".")

		// magnet 链接
		data, err = fetchHTMLData(torrentURL)
		if err != nil {
			log.Error(err)
			return
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
		if err != nil {
			log.Error(err)
			return
		}

		// #openPopup
		magnet, exists = doc.Find("#openPopup").Attr("href")
		if !exists {
			return
		}

		name, _, year, err := torrentName2info(torrentName)
		if err != nil {
			log.Warnf("torrentName2info err: %s", err)
			return
		}

		switch w.typ {
		case types.VideoTypeMovie:
			typ = "movie"
		case types.VideoTypeTV:
			typ = "tv"
		}

		video := &types.FeedVideo{
			Name:        name,
			TorrentName: name,
			TorrentURL:  torrentURL,
			Magnet:      magnet,
			Year:        year,
			Type:        typ,
			RowData:     sql.NullString{},
			Web:         w.web,
			DoubanID:    "",
		}
		videosTemp = append(videosTemp, video)
	})
	return videosTemp, nil
}

func fetchHTMLData(urlStr string) ([]byte, error) {
	c := httpclient.NewHTTPClient()
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("magnetdl req,err:%w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.Do,err:%w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	//nolint:wrapcheck
	return io.ReadAll(resp.Body)
}
