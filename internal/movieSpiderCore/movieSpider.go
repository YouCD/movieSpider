package movieSpiderCore

import (
	"movieSpider/internal/bot"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	feed2 "movieSpider/internal/feed"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/report"
	"movieSpider/internal/spider"
	"movieSpider/internal/spider/douban"
	"movieSpider/internal/spider/wmdb"
	"movieSpider/internal/types"
)

type movieSpider struct {
	feeds    []feed2.Feeder
	download *download.Download
	report   *report.Report
	bot      *bot.TGBot
	spiders  []spider.Spider
}

var ms *movieSpider

type Option interface {
	apply(*movieSpider)
}
type optionFunc func(*movieSpider)

func (f optionFunc) apply(ms *movieSpider) {
	f(ms)
}

func WithFeeds() Option {
	// BTBT
	facFeedBTBT := new(feed2.FactoryBTBT)
	feedBTBT := facFeedBTBT.CreateFeeder(config.BTBT.Scheduling)

	//EZTV
	facFeedEZTV := new(feed2.FactoryEZTV)
	feedEZTV := facFeedEZTV.CreateFeeder(config.EZTV.Scheduling)
	//
	facFeedRarbg := new(feed2.FactoryRARBG)
	// rarbg TV
	var feedRarbgTV feed2.Feeder
	// rarbg Movie
	var feedRarbgMovie feed2.Feeder

	for _, r := range config.RARBG {
		if r != nil {
			if r.Typ == types.ResourceTV {
				feedRarbgTV = facFeedRarbg.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
			if r.Typ == types.ResourceMovie {
				feedRarbgMovie = facFeedRarbg.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
		}
	}

	// GLODLS
	facFeedGLODLS := new(feed2.FactoryGLODLS)
	feedGLODLS := facFeedGLODLS.CreateFeeder(config.GLODLS.Scheduling)

	// TGX
	facFeedTGX := new(feed2.FactoryTGX)
	feedTGXS := facFeedTGX.CreateFeeder(config.TGX.Scheduling)

	// TORLOCK
	facFeedTorlock := new(feed2.FactoryTORLOCK)
	var feedTorlockMovie feed2.Feeder
	var feedTorlockTV feed2.Feeder
	for _, r := range config.TORLOCK {
		if r != nil {
			if r.Typ == types.ResourceTV {
				feedTorlockTV = facFeedTorlock.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
			if r.Typ == types.ResourceMovie {
				feedTorlockMovie = facFeedTorlock.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
		}
	}

	return optionFunc(func(ms *movieSpider) {
		ms.feeds = append(ms.feeds,
			feedBTBT,
			feedEZTV,
			feedRarbgTV,
			feedRarbgMovie,
			feedGLODLS,
			feedTGXS,
			feedTorlockMovie,
			feedTorlockTV,
		)

	})
}
func WithConfigFile(configFile string) Option {
	config.InitConfig(configFile)
	model.NewMovieDB()
	return optionFunc(func(ms *movieSpider) {})
}

func WithReport() Option {
	return optionFunc(func(ms *movieSpider) {
		ms.report = report.NewReport("*/1 * * * *")
		go ms.report.Run()
	})
}
func WithDownload() Option {
	return optionFunc(func(ms *movieSpider) {
		ms.download = download.NewDownloader(config.Downloader.Scheduling)
		go ms.download.Run()
	})
}

func NewMovieSpider(options ...Option) *movieSpider {
	ms = new(movieSpider)
	for _, option := range options {
		option.apply(ms)
	}
	return ms
}

func (m *movieSpider) RunWithFeed() {

	for _, feeder := range m.feeds {
		go func(feeder feed2.Feeder) {
			feeder.Run()
		}(feeder)
	}
}
func (m *movieSpider) RunWithTGBot() {
	ms.bot = bot.NewTgBot(config.TG.BotToken, config.TG.TgIDs)
	go ms.bot.StartBot()
}

func (m *movieSpider) RunWithSpider() {
	// Spider
	m.spiders = append(m.spiders, douban.NewSpiderDouBan(config.DouBan.DoubanUrl, config.DouBan.Scheduling))
	m.spiders = append(m.spiders, wmdb.NewSpiderWmdb(config.DouBan.Scheduling))
	for _, s := range m.spiders {
		go func(spider spider.Spider) {
			spider.Run()
		}(s)
	}

}
