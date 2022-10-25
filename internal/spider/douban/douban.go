package douban

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron/v3"
	httpClient2 "movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type DouBan struct {
	doubanUrl  string
	scheduling string
}

func NewSpiderDouBan(doubanUrl, scheduling string) *DouBan {
	return &DouBan{
		doubanUrl,
		scheduling,
	}
}

func (d *DouBan) Crawler() {

	request, err := http.NewRequest(http.MethodGet, d.doubanUrl, nil)
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

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		//return nil, errors.WithMessage(err, "getMovies goquery")
		log.Error(err)
		return
	}
	//fmt.Println(doc.Text())
	doc.Find("#content > div.grid-16-8.clearfix > div.article > div.grid-view> div").Each(func(i int, s *goquery.Selection) {

		s.Each(func(i int, selection *goquery.Selection) {
			doubanVideo := new(types.DouBanVideo)

			// 片名
			nameStr := selection.Find("div.info> ul > li.title > a > em ").Text()

			var strList []string
			if strings.Contains(nameStr, "/") {
				for _, v := range strings.Split(nameStr, "/") {
					if strings.Contains(v, " ") {
						tempName := strings.Trim(v, " ")
						nedName := strings.ReplaceAll(tempName, " ", ".")
						strList = append(strList, nedName)
					}

				}
			} else {
				if strings.Contains(nameStr, " ") {
					tempName := strings.Trim(nameStr, " ")
					nedName := strings.ReplaceAll(tempName, " ", ".")
					strList = append(strList, nedName)
				} else {
					tempName := strings.Trim(nameStr, " ")
					strList = append(strList, tempName)
				}
			}

			by, _ := json.Marshal(strList)
			doubanVideo.Names = string(by)

			//#content > div.grid-16-8.clearfix > div.article > div.grid-view > div:nth-child(1) > div.info > ul > li.title > a
			val, _ := selection.Find("div.info>ul > li.title > a").Attr("href")

			compileRegex := regexp.MustCompile("[0-9]+")
			matchArr := compileRegex.FindStringSubmatch(val)
			doubanVideo.DoubanID = matchArr[0]
			//#content > div.grid-16-8.clearfix > div.article > div.grid-view > div:nth-child(2) > div.info > ul > li.title > span
			Playable := selection.Find(" div.info > ul > li.title > span").Text()
			Playable = strings.ReplaceAll(Playable, "[", "")
			Playable = strings.ReplaceAll(Playable, "]", "")
			doubanVideo.Playable = Playable
			err = model.NewMovieDB().CreatDouBanVideo(doubanVideo)
			if err != nil {
				log.Error(err)
				return
			}
			log.Warnf("DouBan %s 已保存", doubanVideo.Names)
		})

	})

	return
}

func (d *DouBan) Run() {
	if d.scheduling == "" {
		log.Error("DouBan Scheduling is null")
		os.Exit(1)
	}
	log.Infof("DouBan Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	c.AddFunc(d.scheduling, func() { d.Crawler() })
	c.Start()
}
