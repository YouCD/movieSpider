package douban

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"movieSpider/internal/config"
	httpClient2 "movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"net/http"
	"strings"
	"testing"
)

func TestDouBan_Crawler(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://movie.douban.com/subject/26634250/", nil)
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("User-Agent", "go")
	client := httpClient2.NewHTTPClient()
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

	var d types.RowData
	err = json.Unmarshal([]byte(content), &d)
	if err != nil {
		//return nil, errors.WithMessage(err, "getMovies goquery")
		log.Error(err)
		return
	}
	fmt.Println(content)
}

func TestDouBan_Crawler1(t *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")

	d := &DouBan{
		url:        "https://movie.douban.com/people/251312920/wish",
		scheduling: "tt.fields.scheduling",
	}
	videos := d.Crawler()
	for _, video := range videos {
		fmt.Println(video.Names, video.DatePublished)
	}
}

func TestNewSpiderDouBan(t *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	marshal, err := json.Marshal(config.DouBanList)
	if err != nil {
		log.Error(err)
		return
	}
	t.Log(string(marshal))

	douBan := NewSpiderDouBan(config.DouBanList)
	t.Log(douBan)

}
