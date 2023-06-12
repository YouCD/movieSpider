package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

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
}

//
// NewTgBot
//  @Description: 创建一个TGBot实例
//  @param BotToken
//  @param TgIDs
//  @return *TGBot
//
func NewTgBot(BotToken string, TgIDs []int) *TGBot {
	once.Do(func() {
		client := httpClient.NewHttpClient()
		bot, err := tgbotapi.NewBotAPIWithClient(config.TG.BotToken, "https://api.telegram.org/bot%s/%s", client)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		tgBotClient = &TGBot{
			BotToken, TgIDs, bot,
		}
	})
	return tgBotClient
}

//
// StartBot
//  @Description: 启动bot
//  @receiver t
//
func (t *TGBot) StartBot() {

	// 发送通知 下载通知
	go func() {
		for {
			notifyChan, ok := <-bus.NotifyChan
			if ok {
				log.Infof(notifyChan)
				t.SendStrMsg(notifyChan)
			} else {
				return
			}
		}
	}()
	// 发送通知 上映 通知
	go func() {
		for {
			video, ok := <-bus.DatePublishedChan
			if ok {
				t.SendDatePublishedMsg(video)
			} else {
				return
			}
		}
	}()

	log.Infof("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
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
				aria2Server, err := aria2.NewAria2(config.Downloader.Aria2Label)
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

				count, err := model.NewMovieDB().CountFeedVideo()
				if err != nil {
					log.Error(err)
					continue
				}
				var s string
				var Total int
				for _, reportCount := range count {
					Total += reportCount.Count
					s += fmt.Sprintf("%s: %d ", reportCount.Web, reportCount.Count)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Total: %d  %s", Total, s))
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := t.bot.Send(msg); err != nil {
					log.Error(err)
				}

			// movie_download 指令
			case strings.Contains(update.Message.Text, CMDMoveDownload):
				// 如果参数长度不够直接continue 防止地址越界
				pars, ok := t.checkPars(update.Message.Text, update.Message.Chat.ID, update, CMDMoveDownload)
				if !ok {
					continue
				}
				downloader := download.NewDownloader(config.Downloader.Scheduling)
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

//
// SendStrMsg
//  @Description: 发送字符串消息
//  @receiver t
//  @param msg
//
func (t *TGBot) SendStrMsg(msg string) {
	for _, id := range t.IDs {
		tgMsg := tgbotapi.NewMessage(int64(id), msg)
		if _, err := t.bot.Send(tgMsg); err != nil {
			log.Error(err)
		}
	}
}

const msg = `<b>电影名：</b> %s
<b>上映时间：</b> %s
`

//
// SendDatePublishedMsg
//  @Description: 发送电影上映消息
//  @receiver t
//  @param msg
//
func (t *TGBot) SendDatePublishedMsg(v *types.DouBanVideo) {

	var rowData types.RowData
	err := json.Unmarshal([]byte(v.RowData), &rowData)
	if err != nil {
		log.Error(err)

	}
	image := rowData.Image
	var names []string
	err = json.Unmarshal([]byte(v.Names), &names)
	if err != nil {
		log.Error(err)
	}
	for _, id := range t.IDs {
		photo := tgbotapi.NewPhoto(int64(id), tgbotapi.FileURL(image))
		photo.Caption = fmt.Sprintf(msg, names[0], v.DatePublished)
		photo.ParseMode = tgbotapi.ModeHTML
		if _, err := t.bot.Send(photo); err != nil {
			log.Error(err)
		}
	}
}

//
// getMovieID
//  @Description: 获取电影id
//  @param str
//  @return int
//  @return error
//
func getMovieID(str string) (int, error) {
	sile := strings.Split(str, " ")
	if len(sile) < 2 {
		return 0, errors.New("getMovieID id is 0")
	} else {
		movieID, err := strconv.Atoi(sile[2])
		if err != nil {
			return 0, err
		}
		return movieID, nil
	}

}

//
// getMovieInlineKeyboardMarkup
//  @Description: 获取电影内联键盘
//  @return *tgbotapi.InlineKeyboardMarkup
//
func getMovieInlineKeyboardMarkup() *tgbotapi.InlineKeyboardMarkup {
	if *pageNum <= 1 {
		if *pageNum == 1 {
			Markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("上一页", "0"),
					tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("下一页(第%d页)", *pageNum+1), fmt.Sprintf("%d", *pageNum+1)),
				),
			)
			return &Markup
		}
		if *pageNum == 0 {
			Markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("上一页", "1"),
					tgbotapi.NewInlineKeyboardButtonData("下一页(第2页)", "2"),
				),
			)
			return &Markup
		}
	} else if *pageNum > 1 {
		Markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("上一页(第%d页)", *pageNum-1), fmt.Sprintf("%d", *pageNum-1)),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("下一页(第%d页)", *pageNum+1), fmt.Sprintf("%d", *pageNum+1)),
			),
		)
		return &Markup
	}

	return nil
}

//
// inArray
//  @Description: 判断数组中是否存在某个值
//  @param val
//  @param array
//  @return ok
//  @return i
//
func inArray(val int, array []int) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

//
// checkUser
//  @Description: 检查用户是否有权限
//  @receiver t
//  @param ChatID
//  @param update
//  @return bool
//
func (t *TGBot) checkUser(ChatID int64, update tgbotapi.Update) bool {
	ok, _ := inArray(int(ChatID), config.TG.TgIDs)
	if !ok {
		msg := tgbotapi.NewMessage(ChatID, "您没有权限")
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := t.bot.Send(msg); err != nil {
			log.Error(err)
			return false
		}
		return false
	}
	return ok
}

//
// checkPars
//  @Description: 检查参数
//  @receiver t
//  @param pars
//  @param ChatID
//  @param update
//  @param cmd
//  @return []string
//  @return bool
//
func (t *TGBot) checkPars(pars string, ChatID int64, update tgbotapi.Update, cmd string) ([]string, bool) {
	log.Infof("Msg: %s", update.Message.Text)
	cmdAndargs := tools.RemoveSpaceItem(strings.Split(pars, " "))
	switch cmd {
	case CMDMoveDownload:
		flag := t.checkArgsLen(ChatID, update, cmdAndargs, 2)
		return cmdAndargs, flag
	case CMDReportFeedVideos:
		return cmdAndargs, true
	case CMDReportDownload:
		return cmdAndargs, true
	default:
		return cmdAndargs, false
	}

}

//
// checkArgsLen
//  @Description: 检查参数长度
//  @receiver t
//  @param ChatID
//  @param update
//  @param cmdAndargs
//  @param length
//  @return bool
//
func (t *TGBot) checkArgsLen(ChatID int64, update tgbotapi.Update, cmdAndargs []string, length int) bool {
	if len(cmdAndargs) < length {
		msg := tgbotapi.NewMessage(ChatID, "参数长度不够")
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
