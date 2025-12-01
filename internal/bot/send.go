package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"strings"
	"text/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/youcd/toolkit/log"
)

func (t *TGBot) SendReportFeedVideosMsg(msgChatID, msgID int64) {
	count, err := model.NewMovieDB().CountFeedVideo()
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}
	reportFeedVideosTmpl := template.New("reportFeedVideosTmpl")
	_, _ = reportFeedVideosTmpl.Parse(`<b>Feed数据统计</b>
{{range .}} <b>{{ .Web }}：</b>   {{ .Count }}
{{end}} 
`)
	b := new(bytes.Buffer)
	_ = reportFeedVideosTmpl.Execute(b, count)
	msg := tgbotapi.NewMessage(msgChatID, b.String())
	msg.ReplyToMessageID = int(msgID)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = t.bot.Send(msg)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}
}

func splitSpace(s string, index int) string {
	return strings.Split(s, " ")[index]
}

// 定义模板 结构体
type msgType struct {
	Name          string
	DatePublished string
	MovieURI      string
	Director      []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	Actor []struct {
		Type string `json:"type"`
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	Genre       []string
	Description string
	File        string
	Size        string
	Gid         string
}

// SendDatePublishedOrDownloadMsg
//
//	@Description: 发送电影上映消息
//	@receiver t
//	@param msg
func (t *TGBot) SendDatePublishedOrDownloadMsg(v *types.DownloadNotifyVideo, notify notifyType) {
	video, err := model.NewMovieDB().FetchOneDouBanVideoByDouBanID(v.FeedVideo.DoubanID)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
		return
	}

	// 处理原始信息
	var rowData types.RowData
	err = json.Unmarshal([]byte(video.RowData), &rowData)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}
	// 处理电影名
	var names []string
	err = json.Unmarshal([]byte(video.Names), &names)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}
	// 定义模板 结构体
	var msg = msgType{
		Name:          names[0],
		DatePublished: video.DatePublished,
		MovieURI:      rowData.URL,
		Director:      rowData.Director,
		Actor:         rowData.Actor,
		Genre:         rowData.Genre,
		Description:   rowData.Description,
		File:          v.File,
		Size:          v.Size,
		Gid:           v.Gid,
	}

	datePublishedMsgTmpl := template.New("datePublishedMsgTmpl")
	datePublishedMsgTmpl.Funcs(template.FuncMap{"splitSpace": splitSpace})
	switch notify {
	// 电影下载通知
	case notifyTypeDownload:
		_, _ = datePublishedMsgTmpl.Parse(`<b>下载通知</b>
<b>电影名：</b> {{.Name}} 
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieURI}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.URL}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.URL}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}
<b>文件：</b>   {{ .File }}
<b>Gid：</b>   {{ .Gid }}
`)
	//	上映通知
	case notifyTypeDatePublished:
		_, _ = datePublishedMsgTmpl.Parse(`<b>上映通知</b>
<b>电影名：</b> {{.Name}}
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieURI}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.URL}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.URL}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}
`)
	case notifyTypeDownloadComplete:
		_, _ = datePublishedMsgTmpl.Parse(`<b>下载完毕通知</b>
<b>电影名：</b> {{.Name}}
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieURI}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.URL}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.URL}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}
<b>文件名：</b>   {{ .File }}
<b>大小：</b>   {{ .Size }}
<b>Gid：</b>   {{ .Gid }}
`)
	}

	// 定义缓冲区 用于存储模板渲染后的数据
	b := new(bytes.Buffer)
	err = datePublishedMsgTmpl.Execute(b, msg)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}

	image := rowData.Image

	for _, id := range t.IDs {
		photo := tgbotapi.NewPhoto(int64(id), tgbotapi.FileURL(image))
		photo.Caption = b.String()
		photo.ParseMode = tgbotapi.ModeHTML
		_, err = t.bot.Send(photo)
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
		}
	}
}
