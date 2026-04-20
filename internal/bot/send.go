package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"text/template"

	"movieSpider/internal/model"
	"movieSpider/internal/types"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/youcd/toolkit/log"
)

func (t *TGBot) SendReportFeedVideosMsg(ctx context.Context, msgChatID, msgID int64) {
	count, err := model.NewMovieDB().CountFeedVideo()
	if err != nil {
		log.WithCtx(ctx).Error(err)
	}
	reportFeedVideosTmpl := template.New("reportFeedVideosTmpl")
	_, _ = reportFeedVideosTmpl.Parse(`<b>Feed数据统计</b>
{{range .}} <b>{{ .Web }}：</b>   {{ .Count }}
{{end}} 
`)
	b := new(bytes.Buffer)
	_ = reportFeedVideosTmpl.Execute(b, count)

	_, err = t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    msgChatID,
		Text:      b.String(),
		ParseMode: models.ParseModeHTML,
		ReplyParameters: &models.ReplyParameters{
			MessageID: int(msgID),
		},
	})
	if err != nil {
		log.WithCtx(ctx).Error(err)
	}
}

func splitSpace(s string, index int) string {
	parts := strings.Split(s, " ")
	if index < len(parts) {
		return parts[index]
	}
	return ""
}

// msgType 定义消息模板结构体
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

// baseMsgTmpl 基础消息模板，包含所有通知共用的部分
const baseMsgTmpl = `<b>电影名：</b> {{.Name}}
<b>上映时间：</b> {{.DatePublished}}
<a href="https://movie.douban.com{{.MovieURI}}">豆瓣</a>
<b>导演：</b>  {{range .Director}} <a href="https://movie.douban.com{{.URL}}">{{splitSpace .Name 0}}</a> {{end}} 
<b>演员：</b>  {{range .Actor}} <a href="https://movie.douban.com{{.URL}}">{{ splitSpace .Name 0 }}</a> {{end}} 
<b>类型：</b>  {{range .Genre}} {{ . }} {{end}} 
<b>简介：</b>   {{ .Description }}`

// notificationTemplates 定义不同通知类型的模板
var notificationTemplates = map[notifyType]string{
	notifyTypeDownload: `<b>下载通知</b>
` + baseMsgTmpl + `
<b>文件：</b>   {{ .File }}
<b>Gid：</b>   {{ .Gid }}`,
	notifyTypeDatePublished: `<b>上映通知</b>
` + baseMsgTmpl,
	notifyTypeDownloadComplete: `<b>下载完毕通知</b>
` + baseMsgTmpl + `
<b>文件名：</b>   {{ .File }}
<b>大小：</b>   {{ .Size }}
<b>Gid：</b>   {{ .Gid }}`,
}

// SendDatePublishedOrDownloadMsg 发送电影上映消息或下载通知
func (t *TGBot) SendDatePublishedOrDownloadMsg(ctx context.Context, v *types.DownloadNotifyVideo, notify notifyType) {
	video, err := model.NewMovieDB().FetchOneTMDBVideoByImdbID(ctx, v.FeedVideo.ImdbID)
	if err != nil {
		log.WithCtx(ctx).Error(err)
		return
	}

	// 处理原始信息
	var rowData types.RowData
	if err = json.Unmarshal([]byte(video.RowData), &rowData); err != nil {
		log.WithCtx(ctx).Error(err)
		return
	}

	// 处理电影名
	var names []string
	if err = json.Unmarshal([]byte(video.Names), &names); err != nil {
		log.WithCtx(ctx).Error(err)
		return
	}

	// 定义模板结构体
	msg := msgType{
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

	// 获取对应通知类型的模板
	tmplStr, ok := notificationTemplates[notify]
	if !ok {
		log.WithCtx(ctx).Errorf("unknown notify type: %v", notify)
		return
	}

	datePublishedMsgTmpl := template.New("datePublishedMsgTmpl")
	datePublishedMsgTmpl.Funcs(template.FuncMap{"splitSpace": splitSpace})
	_, err = datePublishedMsgTmpl.Parse(tmplStr)
	if err != nil {
		log.WithCtx(ctx).Error(err)
		return
	}

	// 定义缓冲区 用于存储模板渲染后的数据
	b := new(bytes.Buffer)
	if err = datePublishedMsgTmpl.Execute(b, msg); err != nil {
		log.WithCtx(ctx).Error(err)
		return
	}

	image := rowData.Image

	for _, id := range t.IDs {
		_, err = t.bot.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:    int64(id),
			Photo:     &models.InputFileString{Data: image},
			Caption:   b.String(),
			ParseMode: models.ParseModeHTML,
		})
		if err != nil {
			log.WithCtx(ctx).Error(err)
		}
	}
}
