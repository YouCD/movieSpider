package bot

import (
	"context"
	"fmt"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/youcd/toolkit/log"
)

//nolint:gochecknoglobals,unused
var (
	tgBotClient *TGBot
	once        sync.Once
)

const (
	CMDMoveDownload     = "movie_download"
	CMDReportDownload   = "report_download"
	CMDReportFeedVideos = "report_feedvioes"
)

type TGBot struct {
	botToken string
	IDs      []int
	bot      *tgbotapi.BotAPI
	mtx      sync.Mutex
}

// NewTgBot 创建一个TGBot实例
func NewTgBot(botToken string, tgIDs []int) *TGBot {
	once.Do(func() {
		ctx := context.Background()
		client := http.DefaultClient
		if config.Config.Global.ProxyURL != "" {
			log.WithCtx(ctx).Info("ProxyURL ", config.Config.Global.ProxyURL)
			client = httpclient.NewProxyHTTPClient(ctx)
		}
		bot, err := tgbotapi.NewBotAPIWithClient(config.Config.TG.BotToken, "https://api.telegram.org/bot%s/%s", client)
		if err != nil {
			log.WithCtx(ctx).Error(err)
			os.Exit(-1)
		}

		tgBotClient = &TGBot{
			botToken: botToken, IDs: tgIDs, bot: bot,
		}
	})
	return tgBotClient
}

type notifyType int

const (
	notifyTypeDownload notifyType = iota + 1
	notifyTypeDownloadComplete
	notifyTypeDatePublished
)

// StartBot 启动bot
//
//nolint:gocognit
func (t *TGBot) StartBot() {
	// 发送通知 下载 通知
	t.downloadNotify()
	// 发送通知 上映 通知
	t.datePublishedNotify()
	// 发送通知 下载完毕 通知
	t.downloadCompleteNotify()
	log.WithCtx(context.Background()).Infof("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	updates := t.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		if !t.checkUser(update.Message.Chat.ID, update) {
			log.WithCtx(context.Background()).Warnf("用户 %s(%d) 没有权限执行命令 %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Command())
			continue
		}

		ctx := context.Background()
		switch update.Message.Command() {
		case CMDReportDownload:
			t.handleReportDownload(ctx, update)
		case CMDReportFeedVideos:
			t.SendReportFeedVideosMsg(update.Message.Chat.ID, int64(update.Message.MessageID))
		case CMDMoveDownload:
			t.handleMovieDownload(ctx, update)
		default:
			t.sendReplyMessage(update.Message.Chat.ID, update.Message.MessageID, "不支持此指令")
		}
	}
}

// handleReportDownload 处理下载报告命令
func (t *TGBot) handleReportDownload(ctx context.Context, update tgbotapi.Update) {
	aria2Server, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		log.WithCtx(ctx).Error(err)
		t.sendReplyMessage(update.Message.Chat.ID, update.Message.MessageID, "aria2下载器服务异常")
		return
	}

	files := aria2Server.CurrentActiveAndStopFiles()
	var sb strings.Builder
	for _, file := range files {
		fileName := truncateString(file.FileName, 40)
		sb.WriteString(fmt.Sprintf("%-40s | %s\n", fileName, file.Completed))
	}

	t.sendReplyMessage(update.Message.Chat.ID, update.Message.MessageID, sb.String())
}

// handleMovieDownload 处理电影下载命令
func (t *TGBot) handleMovieDownload(ctx context.Context, update tgbotapi.Update) {
	update.Message.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 15}}

	arguments := update.Message.CommandArguments()
	pars := tools.RemoveSpaceItem(strings.Split(arguments, " "))

	if len(pars) < 3 {
		t.sendReplyMessage(update.Message.Chat.ID, update.Message.MessageID, "参数不足，格式: /movie_download <name> <resolution>")
		return
	}

	downloader := download.NewDownloader(config.Config.Downloader.Scheduling)
	downloadMsg := downloader.DownloadByName(ctx, pars[1], pars[2])
	t.sendReplyMessage(update.Message.Chat.ID, update.Message.MessageID, downloadMsg)
}

// sendReplyMessage 发送回复消息
func (t *TGBot) sendReplyMessage(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	_, err := t.bot.Send(msg)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	nameRune := []rune(s)
	return string(nameRune[:maxLen])
}

// SendStrMsg 发送字符串消息
func (t *TGBot) SendStrMsg(msg string) {
	for _, id := range t.IDs {
		tgMsg := tgbotapi.NewMessage(int64(id), msg)
		_, err := t.bot.Send(tgMsg)
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
		}
	}
}

// inArray 判断数组中是否存在某个值
func inArray(val int, array []int) bool {
	for _, v := range array {
		if v == val {
			return true
		}
	}
	return false
}

// checkUser 检查用户是否有权限
func (t *TGBot) checkUser(chatID int64, update tgbotapi.Update) bool {
	if !inArray(int(chatID), config.Config.TG.TgIDs) {
		t.sendReplyMessage(chatID, update.Message.MessageID, "您没有权限")
		return false
	}
	return true
}

// downloadNotify 下载通知
func (t *TGBot) downloadNotify() {
	go func() {
		for {
			video, ok := <-bus.DownloadNotifyChan
			if ok {
				t.SendDatePublishedOrDownloadMsg(video, notifyTypeDownload)
			} else {
				return
			}
		}
	}()
}

// datePublishedNotify 上映通知
func (t *TGBot) datePublishedNotify() {
	go func() {
		for {
			v, ok := <-bus.DatePublishedChan
			if ok {
				t.SendDatePublishedOrDownloadMsg(&types.DownloadNotifyVideo{
					TMDBVideo: v,
				}, notifyTypeDatePublished)
			} else {
				return
			}
		}
	}()
}

// downloadCompleteNotify 下载完成通知
func (t *TGBot) downloadCompleteNotify() {
	downLoadChan := make(chan *types.DownloadNotifyVideo)
	go func() {
		defer close(downLoadChan)
		aria2Server, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			return
		}
		for {
			time.Sleep(time.Second * 1)
			t.mtx.Lock()
			aria2Server.Subscribe(downLoadChan)

			select {
			case video, ok := <-downLoadChan:
				if ok {
					func() {
						defer t.mtx.Unlock()
						t.SendDatePublishedOrDownloadMsg(video, notifyTypeDownloadComplete)
					}()
				}
			}
		}
	}()
}
