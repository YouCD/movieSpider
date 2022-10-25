package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

func TestNewTgBot(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()
	bot := NewTgBot(config.TG.BotToken, config.TG.TgIDs)

	//bot.StartBot()
	msg := tgbotapi.NewMessage(221941736, "downloadMsg")
	msg.AllowSendingWithoutReply = true

	if _, err := bot.bot.Send(msg); err != nil {
		t.Error(err)
	}
}
