package bot

import (
	"bytes"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"strings"
	"text/template"
)

func (t *TGBot) SendReportFeedVideosMsg(msgChatID, msgID int64) {
	count, err := model.NewMovieDB().CountFeedVideo()
	if err != nil {
		log.Error(err)
	}
	reportFeedVideosTmpl := template.New("reportFeedVideosTmpl")
	reportFeedVideosTmpl.Parse(`<b>Feed数据统计</b>
{{range .}} <b>{{ .Web }}：</b>   {{ .Count }}
{{end}} 
`)
	b := new(bytes.Buffer)
	reportFeedVideosTmpl.Execute(b, count)
	msg := tgbotapi.NewMessage(msgChatID, b.String())
	msg.ReplyToMessageID = int(msgID)
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := t.bot.Send(msg); err != nil {
		log.Error(err)
	}

}

func splitSpace(s string, index int) string {
	return strings.Split(s, " ")[index]
}

// 定义模板 结构体
type msgType struct {
	Name          string
	DatePublished string
	MovieUri      string
	Director      []struct {
		Type string `json:"type"`
		Url  string `json:"url"`
		Name string `json:"name"`
	}
	Actor []struct {
		Type string `json:"type"`
		Url  string `json:"url"`
		Name string `json:"name"`
	}
	Genre       []string
	Description string
	File        string
	Size        string
}

// SendDatePublishedOrDownloadMsg
//  @Description: 发送电影上映消息
//  @receiver t
//  @param msg
//
func (t *TGBot) SendDatePublishedOrDownloadMsg(v *types.DouBanVideo, notify notifyType, args ...string) {
	// 处理原始信息
	var rowData types.RowData
	err := json.Unmarshal([]byte(v.RowData), &rowData)
	if err != nil {
		log.Error(err)
	}
	// 处理电影名
	var names []string
	err = json.Unmarshal([]byte(v.Names), &names)
	if err != nil {
		log.Error(err)
	}
	// 定义模板 结构体
	var msg = msgType{
		Name:          names[0],
		DatePublished: v.DatePublished,
		MovieUri:      rowData.Url,
		Director:      rowData.Director,
		Actor:         rowData.Actor,
		Genre:         rowData.Genre,
		Description:   rowData.Description,
	}
	if len(args) > 1 {
		msg.File = args[0]
		msg.Size = args[1]
	}

	datePublishedMsgTmpl := template.New("datePublishedMsgTmpl")
	datePublishedMsgTmpl.Funcs(template.FuncMap{"splitSpace": splitSpace})
	switch notify {
	// 电影下载通知
	case notifyTypeDownload:
		datePublishedMsgTmpl.Parse(`<b>下载通知</b>
<b>电影名：</b> {{.Name}} 
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieUri}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.Url}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.Url}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}
`)
	//	上映通知
	case notifyTypeDatePublished:
		datePublishedMsgTmpl.Parse(`<b>上映通知</b>
<b>电影名：</b> {{.Name}}
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieUri}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.Url}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.Url}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}
`)
	case notifyTypeDownloadComplete:
		datePublishedMsgTmpl.Parse(`<b>下载完毕通知</b>
<b>电影名：</b> {{.Name}}
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieUri}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.Url}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.Url}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}
<b>文件名：</b>   {{ .File }}
<b>大小：</b>   {{ .Size }}
`)
	}

	// 定义缓冲区 用于存储模板渲染后的数据
	b := new(bytes.Buffer)
	err = datePublishedMsgTmpl.Execute(b, msg)
	if err != nil {
		log.Error(err)
	}

	image := rowData.Image

	for _, id := range t.IDs {
		photo := tgbotapi.NewPhoto(int64(id), tgbotapi.FileURL(image))
		photo.Caption = b.String()
		photo.ParseMode = tgbotapi.ModeHTML
		if _, err := t.bot.Send(photo); err != nil {
			log.Error(err)
		}
	}
}
