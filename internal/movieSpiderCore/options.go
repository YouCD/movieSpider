package movieSpiderCore

import (
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/feed"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/report"
	"movieSpider/internal/types"
)

type Option interface {
	apply(*movieSpider)
}
type optionFunc func(*movieSpider)

func (f optionFunc) apply(ms *movieSpider) {
	f(ms)
}

//
// WithFeeds
//  @Description: 初始化feeds
//  @param feeds
//  @return Option
//
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
	facFeedTPBPIRATEPROXY := new(feed.FactoryTPBPIRATEPROXY)
	feedTPBPIRATEPROXY := facFeedTPBPIRATEPROXY.CreateFeeder(config.TPBPIRATEPROXY.Scheduling)

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
			feedTPBPIRATEPROXY,
		)
		ms.feeds = append(ms.feeds, feeds...)
	})
}

//
// WithConfigFile
//  @Description: 初始化配置文件
//  @param configFile
//  @return Option
//
func WithConfigFile(configFile string) Option {
	config.InitConfig(configFile)
	model.NewMovieDB()
	return optionFunc(func(ms *movieSpider) {})
}

//
// WithReport
//  @Description: 初始化 report
//  @return Option
//
func WithReport() Option {
	return optionFunc(func(ms *movieSpider) {
		ms.report = report.NewReport("*/1 * * * *")
		go ms.report.Run()
	})
}

//
// WithDownload
//  @Description: 初始化下载器
//  @return Option
//
func WithDownload() Option {
	return optionFunc(func(ms *movieSpider) {
		ms.download = download.NewDownloader(config.Downloader.Scheduling)
		go ms.download.Run()
	})
}
