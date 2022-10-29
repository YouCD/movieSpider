package movieSpiderCore

import (
	"movieSpider/internal/bot"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/feed"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/report"
	"movieSpider/internal/spider"
	"movieSpider/internal/spider/douban"
	"movieSpider/internal/spider/wmdb"
	"movieSpider/internal/types"
)

type movieSpider struct {
	feeds    []feed.Feeder
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

func WithFeeds(feeds ...feed.Feeder) Option {
	// BTBT
	facFeedBTBT := new(feed.FactoryBTBT)
	feedBTBT := facFeedBTBT.CreateFeeder(config.BTBT.Scheduling)

	//EZTV
	facFeedEZTV := new(feed.FactoryEZTV)
	feedEZTV := facFeedEZTV.CreateFeeder(config.EZTV.Scheduling)
	//
	facFeedRarbg := new(feed.FactoryRARBG)
	// rarbg TV
	var feedRarbgTV feed.Feeder
	// rarbg Movie
	var feedRarbgMovie feed.Feeder

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
	facFeedGLODLS := new(feed.FactoryGLODLS)
	feedGLODLS := facFeedGLODLS.CreateFeeder(config.GLODLS.Scheduling)

	// TGX
	facFeedTGX := new(feed.FactoryTGX)
	feedTGXS := facFeedTGX.CreateFeeder(config.TGX.Scheduling)

	// TORLOCK
	facFeedTorlock := new(feed.FactoryTORLOCK)
	var feedTorlockMovie feed.Feeder
	var feedTorlockTV feed.Feeder
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
	// MAGNETDL
	facFeedMagnetdl := new(feed.FactoryMAGNETDL)
	var feedMagnetdlMovie feed.Feeder
	var feedMagnetdlTV feed.Feeder
	for _, r := range config.TORLOCK {
		if r != nil {
			if r.Typ == types.ResourceTV {
				feedMagnetdlTV = facFeedMagnetdl.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
			if r.Typ == types.ResourceMovie {
				feedMagnetdlMovie = facFeedMagnetdl.CreateFeeder(r.Scheduling, r.Typ)
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
			feedMagnetdlMovie,
			feedMagnetdlTV,
		)
		ms.feeds = append(ms.feeds, feeds...)
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
		go func(feeder feed.Feeder) {
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
