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
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
	bot      *bot.Bot
	mtx      sync.Mutex
}

// NewTgBot 创建一个TGBot实例
func NewTgBot(botToken string, tgIDs []int) *TGBot {
	once.Do(func() {
		ctx := context.Background()
		opts := []bot.Option{}

		if config.Config.Global.ProxyURL != "" {
			log.WithCtx(ctx).Info("ProxyURL ", config.Config.Global.ProxyURL)
			client := httpclient.NewProxyHTTPClient(ctx)
			opts = append(opts, bot.WithHTTPClient(time.Second*60, client))
		}

		b, err := bot.New(config.Config.TG.BotToken, opts...)
		if err != nil {
			log.WithCtx(ctx).Error(err)
			os.Exit(-1)
		}

		tgBotClient = &TGBot{
			botToken: botToken, IDs: tgIDs, bot: b,
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
func (t *TGBot) StartBot(ctx context.Context) {
	// 发送通知 下载 通知
	t.downloadNotify()
	// 发送通知 上映 通知
	t.datePublishedNotify()
	// 发送通知 下载完毕 通知
	t.downloadCompleteNotify()

	// 获取 bot 信息
	me, err := t.bot.GetMe(ctx)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取 bot 信息失败: %v", err)
	} else {
		log.WithCtx(ctx).Infof("Authorized on account %s", me.Username)
	}

	// 注册命令处理器
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, CMDReportDownload, bot.MatchTypeCommand, t.handleReportDownloadHandler)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, CMDReportFeedVideos, bot.MatchTypeCommand, t.handleReportFeedVideosHandler)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, CMDMoveDownload, bot.MatchTypeCommand, t.handleMovieDownloadHandler)

	// 启动 bot
	t.bot.Start(ctx)
}

// handleReportDownloadHandler 处理下载报告命令
func (t *TGBot) handleReportDownloadHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !t.checkUser(update.Message.Chat.ID, update) {
		log.WithCtx(ctx).Warnf("用户 %s(%d) 没有权限执行命令", update.Message.From.Username, update.Message.From.ID)
		return
	}

	aria2Server, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		log.WithCtx(ctx).Error(err)
		t.sendReplyMessage(update.Message.Chat.ID, update.Message.ID, "aria2下载器服务异常")
		return
	}

	files := aria2Server.CurrentActiveAndStopFiles()
	var sb strings.Builder
	for _, file := range files {
		fileName := truncateString(file.FileName, 40)
		sb.WriteString(fmt.Sprintf("%-40s | %s\n", fileName, file.Completed))
	}

	t.sendReplyMessage(update.Message.Chat.ID, update.Message.ID, sb.String())
}

// handleReportFeedVideosHandler 处理 Feed 视频报告命令
func (t *TGBot) handleReportFeedVideosHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !t.checkUser(update.Message.Chat.ID, update) {
		log.WithCtx(ctx).Warnf("用户 %s(%d) 没有权限执行命令", update.Message.From.Username, update.Message.From.ID)
		return
	}

	t.SendReportFeedVideosMsg(ctx, update.Message.Chat.ID, int64(update.Message.ID))
}

// handleMovieDownloadHandler 处理电影下载命令
func (t *TGBot) handleMovieDownloadHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !t.checkUser(update.Message.Chat.ID, update) {
		log.WithCtx(ctx).Warnf("用户 %s(%d) 没有权限执行命令", update.Message.From.Username, update.Message.From.ID)
		return
	}

	// 获取命令参数
	text := update.Message.Text
	// 移除命令部分，获取参数
	args := strings.TrimPrefix(text, "/"+CMDMoveDownload)
	args = strings.TrimSpace(args)
	pars := tools.RemoveSpaceItem(strings.Split(args, " "))

	if len(pars) < 2 {
		t.sendReplyMessage(update.Message.Chat.ID, update.Message.ID, "参数不足，格式: /movie_download <name> <resolution>")
		return
	}

	downloader := download.NewDownloader(config.Config.Downloader.Scheduling)
	downloadMsg := downloader.DownloadByName(ctx, pars[0], pars[1])
	t.sendReplyMessage(update.Message.Chat.ID, update.Message.ID, downloadMsg)
}

// sendReplyMessage 发送回复消息
func (t *TGBot) sendReplyMessage(chatID int64, messageID int, text string) {
	ctx := context.Background()
	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
		ReplyParameters: &models.ReplyParameters{
			MessageID: messageID,
		},
	})
	if err != nil {
		log.WithCtx(ctx).Error(err)
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
	ctx := context.Background()
	for _, id := range t.IDs {
		_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: int64(id),
			Text:   msg,
		})
		if err != nil {
			log.WithCtx(ctx).Error(err)
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
func (t *TGBot) checkUser(chatID int64, update *models.Update) bool {
	if !inArray(int(chatID), config.Config.TG.TgIDs) {
		t.sendReplyMessage(chatID, update.Message.ID, "您没有权限")
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
				t.SendDatePublishedOrDownloadMsg(context.Background(), video, notifyTypeDownload)
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
				t.SendDatePublishedOrDownloadMsg(context.Background(), &types.DownloadNotifyVideo{
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
						t.SendDatePublishedOrDownloadMsg(context.Background(), video, notifyTypeDownloadComplete)
					}()
				}
			}
		}
	}()
}
