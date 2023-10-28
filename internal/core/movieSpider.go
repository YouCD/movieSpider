package core

import (
	"movieSpider/internal/bot"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/spider"
	"movieSpider/internal/spider/douban"
	"movieSpider/internal/spider/feedspider"
)

type MovieSpider struct {
	feeds          []feedspider.Feeder
	download       *download.Download
	report         *job.Report
	bot            *bot.TGBot
	spiders        []spider.Spider
	releaseTimeJob *job.ReleaseTimeJob
}

//nolint:gochecknoglobals
var ms = new(MovieSpider)

// NewMovieSpider
//
//	@Description: 初始化movieSpider
//	@param options
//	@return *MovieSpider
func NewMovieSpider(options ...Option) *MovieSpider {
	for _, option := range options {
		option.apply(ms)
	}
	return ms
}

// RunWithFeed
//
//	@Description: 运行feed
//	@receiver m
func (m *MovieSpider) RunWithFeed() {
	for _, feeder := range m.feeds {
		go func(feeder feedspider.Feeder) {
			feeder.Run(bus.FeedVideoChan)
		}(feeder)
	}
}

// RunWithTGBot
//
//	@Description: 运行tgbot
//	@receiver m
func (m *MovieSpider) RunWithTGBot() {
	if config.TG.Enable {
		ms.bot = bot.NewTgBot(config.TG.BotToken, config.TG.TgIDs)
		go ms.bot.StartBot()
	}
}

// RunWithFeedSpider
//
//	@Description: 运行feedSpider
//	@receiver m
func (m *MovieSpider) RunWithFeedSpider() {
	// Spider
	m.spiders = append(m.spiders, douban.NewSpiderDouBan(config.DouBanList)...)
	for _, s := range m.spiders {
		go func(spider spider.Spider) {
			spider.Run()
		}(s)
	}
}
