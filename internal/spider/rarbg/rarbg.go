package rarbg

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"net/http"
	"strings"
)

const (
	rarbgUrlPrefix        = "https://rarbg.to/torrents.php?search="
	rarbgTorrentUrlPrefix = "https://rarbg.to"
)

type RarbgVideo struct {
	TorrentName string
	Magnet      string
	TorrentUrl  string
}

func newRarbgReq(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	request.Header.Set("Cookie", "__cfduid=d5555d218c15100c8c9352b7cebf2825f1571727890; gaDts48g=q8h5pp9t; aby=1; skt=D8F9Bz5qm2; skt=D8F9Bz5qm2; gaDts48g=q8h5pp9t")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-User", "?1")

	return request, nil
}

func crawlerRarbg(ImdbID string) (Videos []*RarbgVideo, err error) {
	req, err := newRarbgReq(rarbgUrlPrefix + ImdbID)
	if err != nil {
		return nil, err
	}

	client := httpClient.NewHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp == nil {
		log.Warn("未能正常获取Rarbg数据")
		return nil, errors.New("未能正常获取Rarbg数据")
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		//return nil, errors.WithMessage(err, "getMovies goquery")
		log.Error(err)
		return nil, err
	}

	doc.Find("body > table:nth-child(6) > tbody > tr > td:nth-child(2) > div > table > tbody > tr:nth-child(2) > td > table.lista2t > tbody> tr").Each(func(i int, s *goquery.Selection) {

		s.Each(func(i int, selection *goquery.Selection) {
			var v RarbgVideo
			// 种子名
			nameStr := selection.Find("td:nth-child(2)> a:nth-child(1)").Text()
			// 屏蔽 杜比视界片源 720低码 首行表头
			if strings.Contains(nameStr, "720p") || strings.Contains(nameStr, "DVSUX") || strings.Contains(nameStr, "File") {
				return
			}
			v.TorrentName = nameStr

			val, _ := selection.Find("td:nth-child(2) > a:nth-child(1)").Attr("href")

			v.TorrentUrl = rarbgTorrentUrlPrefix + val

			Videos = append(Videos, &v)
		})

	})

	return
}
func crawlerRarbgMagnet(torrentUrl string) (magnet string, err error) {
	req, err := newRarbgReq(torrentUrl)
	if err != nil {
		return "", err
	}

	client := httpClient.NewHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if resp == nil {
		log.Warn("未能正常获取RarbgTorrent数据")
		return "", errors.New("未能正常获取RarbgTorrent数据")
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	selector := "body > table:nth-child(6) > tbody > tr > td:nth-child(2) > div > table > tbody > tr:nth-child(2) > td > div > table > tbody > tr:nth-child(1) > td.lista > a:nth-child(3)"
	magnet, _ = doc.Find(selector).Attr("href")

	return
}
