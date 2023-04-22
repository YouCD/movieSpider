package movieSpiderCore

import (
	"movieSpider/internal/bot"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/feed"
	"movieSpider/internal/report"
	"movieSpider/internal/spider"
	"movieSpider/internal/spider/douban"
)

type movieSpider struct {
	feeds    []feed.Feeder
	download *download.Download
	report   *report.Report
	bot      *bot.TGBot
	spiders  []spider.Spider
}

var ms = new(movieSpider)

//
// NewMovieSpider
//  @Description: 初始化movieSpider
//  @param options
//  @return *movieSpider
//
func NewMovieSpider(options ...Option) *movieSpider {
	for _, option := range options {
		option.apply(ms)
	}
	return ms
}

//
// RunWithFeed
//  @Description: 运行feed
//  @receiver m
//
func (m *movieSpider) RunWithFeed() {
	for _, feeder := range m.feeds {
		go func(feeder feed.Feeder) {
			feeder.Run()
		}(feeder)
	}
}

//
// RunWithTGBot
//  @Description: 运行tgbot
//  @receiver m
//
func (m *movieSpider) RunWithTGBot() {
	if config.TG.Enable {
		ms.bot = bot.NewTgBot(config.TG.BotToken, config.TG.TgIDs)
		go ms.bot.StartBot()
	}

}

//
// RunWithSpider
//  @Description: 运行spider
//  @receiver m
//
func (m *movieSpider) RunWithSpider() {
	// Spider
	m.spiders = append(m.spiders, douban.NewSpiderDouBan(config.DouBan.DoubanUrl, config.DouBan.Scheduling))
	for _, s := range m.spiders {
		go func(spider spider.Spider) {
			spider.Run()
		}(s)
	}

}
