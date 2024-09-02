package bot

import (
	"fmt"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/model"
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
	pageNum     *int
	tgBotClient *TGBot
	once        sync.Once
)

const (
	CMDMoveDownload     = "/movie_download"
	CMDReportDownload   = "/report_download"
	CMDReportFeedVideos = "/report_feedvioes"
)

type TGBot struct {
	botToken string
	IDs      []int
	bot      *tgbotapi.BotAPI
	mtx      sync.Mutex
}

// NewTgBot
//
//	@Description: 创建一个TGBot实例
//	@param BotToken
//	@param TgIDs
//	@return *TGBot
func NewTgBot(botToken string, tgIDs []int) *TGBot {
	once.Do(func() {
		client := http.DefaultClient
		if config.Config.TG.ProxyURL != "" {
			log.Info(config.Config.TG.ProxyURL)
			client = httpclient.NewProxyHTTPClient(config.Config.TG.ProxyURL)
		}
		bot, err := tgbotapi.NewBotAPIWithClient(config.Config.TG.BotToken, "https://api.telegram.org/bot%s/%s", client)
		if err != nil {
			log.Error(err)
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

// StartBot
//
//	@Description: 启动bot
//	@receiver t
//
//nolint:gocognit
func (t *TGBot) StartBot() {
	// 发送通知 下载 通知
	t.downloadNotify()
	// 发送通知 上映 通知
	t.datePublishedNotify()
	// 发送通知 下载完毕 通知
	t.downloadCompleteNotify()
	log.Infof("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	updates := t.bot.GetUpdatesChan(u)
	for update := range updates {
		index := 1
		pageNum = &index
		if update.Message != nil {
			switch {
			case strings.Contains(update.Message.Text, CMDReportDownload):
				// 如果参数长度不够直接continue 防止地址越界
				_, ok := t.checkPars(update.Message.Text, update.Message.Chat.ID, update, CMDReportDownload)
				if !ok {
					continue
				}
				aria2Server, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
				if err != nil {
					log.Error(err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "aria2下载器服务异常")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := t.bot.Send(msg); err != nil {
						log.Error(err)
					}
					continue
				}

				files := aria2Server.CurrentActiveAndStopFiles()
				var bs string
				for _, file := range files {
					if utf8.RuneCountInString(file.FileName) > 40 {
						nameRune := []rune(file.FileName)
						bs += fmt.Sprintf("%-40s | %s\n", string(nameRune[0:40]), file.Completed)
					} else {
						bs += fmt.Sprintf("%-40s | %s\n", file.FileName, file.Completed)
					}
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, bs)
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := t.bot.Send(msg); err != nil {
					log.Error(err)
				}

			case strings.Contains(update.Message.Text, CMDReportFeedVideos):
				_, ok := t.checkPars(update.Message.Text, update.Message.Chat.ID, update, CMDReportFeedVideos)
				if !ok {
					continue
				}
				t.SendReportFeedVideosMsg(update.Message.Chat.ID, int64(update.Message.MessageID))

			// movie_download 指令
			case strings.Contains(update.Message.Text, CMDMoveDownload):
				// 如果参数长度不够直接continue 防止地址越界
				pars, ok := t.checkPars(update.Message.Text, update.Message.Chat.ID, update, CMDMoveDownload)
				if !ok {
					continue
				}
				downloader := download.NewDownloader(config.Config.Downloader.Scheduling)
				downloadMsg := downloader.DownloadByName(pars[1], pars[2])
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, downloadMsg)
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := t.bot.Send(msg); err != nil {
					log.Error(err)
				}
			}
		}
	}
}

// SendStrMsg
//
//	@Description: 发送字符串消息
//	@receiver t
//	@param msg
func (t *TGBot) SendStrMsg(msg string) {
	for _, id := range t.IDs {
		tgMsg := tgbotapi.NewMessage(int64(id), msg)
		if _, err := t.bot.Send(tgMsg); err != nil {
			log.Error(err)
		}
	}
}

// inArray
//
//	@Description: 判断数组中是否存在某个值
//	@param val
//	@param array
//	@return ok
//	@return i
//
//nolint:unused
func inArray(val int, array []int) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

// checkUser
//
//	@Description: 检查用户是否有权限
//	@receiver t
//	@param ChatID
//	@param update
//	@return bool
//
//nolint:unused
func (t *TGBot) checkUser(chatID int64, update tgbotapi.Update) bool {
	ok, _ := inArray(int(chatID), config.Config.TG.TgIDs)
	if !ok {
		msg := tgbotapi.NewMessage(chatID, "您没有权限")
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := t.bot.Send(msg); err != nil {
			log.Error(err)
			return false
		}
		return false
	}
	return ok
}

// checkPars
//
//	@Description: 检查参数
//	@receiver t
//	@param pars
//	@param ChatID
//	@param update
//	@param cmd
//	@return []string
//	@return bool
func (t *TGBot) checkPars(pars string, chatID int64, update tgbotapi.Update, cmd string) ([]string, bool) {
	log.Infof("Msg: %s", update.Message.Text)
	cmdAndargs := tools.RemoveSpaceItem(strings.Split(pars, " "))
	switch cmd {
	case CMDMoveDownload:
		flag := t.checkArgsLen(chatID, update, cmdAndargs, 2)
		return cmdAndargs, flag
	case CMDReportFeedVideos:
		return cmdAndargs, true
	case CMDReportDownload:
		return cmdAndargs, true
	default:
		return cmdAndargs, false
	}
}

// checkArgsLen
//
//	@Description: 检查参数长度
//	@receiver t
//	@param ChatID
//	@param update
//	@param cmdAndargs
//	@param length
//	@return bool
func (t *TGBot) checkArgsLen(chatID int64, update tgbotapi.Update, cmdAndargs []string, length int) bool {
	if len(cmdAndargs) < length {
		msg := tgbotapi.NewMessage(chatID, "参数长度不够")
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := t.bot.Send(msg); err != nil {
			log.Error(err)
			return false
		}
		log.Warnf("参数长度不够")
		return false
	}
	return true
}

// downloadNotify
//
//	@Description: 下载通知
//	@receiver t
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

// datePublishedNotify
//
//	@Description: 上映通知
//	@receiver t
//

func (t *TGBot) datePublishedNotify() {
	go func() {
		for {
			v, ok := <-bus.DatePublishedChan
			if ok {
				video, err := model.NewMovieDB().FetchOneDouBanVideoByDouBanID(v.DoubanID)
				if err != nil {
					log.Error(err)
				}

				t.SendDatePublishedOrDownloadMsg(&types.DownloadNotifyVideo{
					DouBanVideo: video,
				}, notifyTypeDatePublished)
			} else {
				return
			}
		}
	}()
}

// downloadCompleteNotify
//
//	@Description: 下载完成通知
//	@receiver t
func (t *TGBot) downloadCompleteNotify() {
	downLoadChan := make(chan *types.DownloadNotifyVideo)
	go func() {
		defer close(downLoadChan)
		aria2Server, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
		if err != nil {
			log.Error(err)
			return
		}
		for {
			time.Sleep(time.Second * 1)
			t.mtx.Lock()
			aria2Server.Subscribe(downLoadChan)
			//nolint:gosimple
			select {
			case video, ok := <-downLoadChan:
				if ok {
					t.SendDatePublishedOrDownloadMsg(video, notifyTypeDownloadComplete)
					t.mtx.Unlock()
				}
			}
		}
	}()
}
