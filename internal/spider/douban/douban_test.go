package douban

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"movieSpider/internal/config"
	httpClient2 "movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"net/http"
	"strings"
	"testing"
)

func TestDouBan_Crawler(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://movie.douban.com/subject/34825964/", nil)
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("User-Agent", "go")
	client := httpClient2.NewHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	if resp == nil {
		log.Warn("未能正常获取豆瓣数据")
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		//return nil, errors.WithMessage(err, "getMovies goquery")
		log.Error(err)
		return
	}
	//<script type="application/ld+json">

	//fmt.Println(doc.Html())
	content := doc.Find("script[type='application/ld+json']").Text()
	content = strings.ReplaceAll(content, "\n", "")

	var d rowData
	err = json.Unmarshal([]byte(content), &d)
	if err != nil {
		//return nil, errors.WithMessage(err, "getMovies goquery")
		log.Error(err)
		return
	}
	fmt.Println(d)
}

func TestDouBan_Crawler1(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")

	d := &DouBan{
		doubanUrl:  "https://movie.douban.com/people/251312920/wish",
		scheduling: "tt.fields.scheduling",
	}
	videos := d.Crawler()
	for _, video := range videos {
		fmt.Println(video)
	}
}
