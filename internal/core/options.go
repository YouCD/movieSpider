package core

import (
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	feed2 "movieSpider/internal/spider/feedSpider"
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
func WithFeeds(feeds ...feed2.Feeder) Option {
	// BTBT
	facFeedBTBT := new(feed2.FactoryBTBT)
	feedBTBT := facFeedBTBT.CreateFeeder(config.BTBT.Scheduling)

	//EZTV
	facFeedEZTV := new(feed2.FactoryEZTV)
	feedEZTV := facFeedEZTV.CreateFeeder(config.EZTV.Scheduling)

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
	// MAGNETDL
	facFeedMagnetdl := new(feed2.FactoryMAGNETDL)
	var feedMagnetdlMovie feed2.Feeder
	var feedMagnetdlTV feed2.Feeder
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
	facFeedTPBPIRATEPROXY := new(feed2.FactoryTPBPIRATEPROXY)
	feedTPBPIRATEPROXY := facFeedTPBPIRATEPROXY.CreateFeeder(config.TPBPIRATEPROXY.Scheduling)

	return optionFunc(func(ms *movieSpider) {
		ms.feeds = append(ms.feeds,
			feedBTBT,
			feedEZTV,
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
		ms.report = job.NewReport("*/1 * * * *")
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

//
// WithReleaseTimeJob
//  @Description: 初始化下载器
//  @return Option
//
func WithReleaseTimeJob() Option {
	if !config.TG.Enable {
		return optionFunc(func(ms *movieSpider) {
			log.Warn("未开启TG通知，无法运行 电影上线 通知job")
		})
	}
	return optionFunc(func(ms *movieSpider) {
		ms.releaseTimeJob = job.NewReleaseTimeJob("")
		go ms.releaseTimeJob.Run()
	})
}
