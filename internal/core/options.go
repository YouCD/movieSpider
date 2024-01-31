package core

import (
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/feedspider"
	"movieSpider/internal/types"
)

//nolint:inamedparam
type Option interface {
	apply(*MovieSpider)
}
type optionFunc func(*MovieSpider)

func (f optionFunc) apply(ms *MovieSpider) {
	f(ms)
}

// WithFeeds
//
//	@Description: 初始化feeds
//	@param feeds
//	@return Option
func WithFeeds(feeds ...feedspider.Feeder) Option {
	// BTBT
	feedBTBT := feedspider.NewBtbt(config.Config.Feed.BTBT.Scheduling)

	// EZTV
	feedEZTV := feedspider.NewEztv(config.Config.Feed.EZTV.Scheduling, config.Config.Feed.EZTV.MirrorSite)

	// GLODLS
	feedGLODLS := feedspider.NewGlodls(config.Config.Feed.GLODLS.Scheduling, config.Config.Feed.GLODLS.MirrorSite)

	// TGX
	feedTGXS := feedspider.NewTgx(config.Config.Feed.TGX.Scheduling, config.Config.Feed.TGX.MirrorSite)

	// TORLOCK
	var feedTorlockMovie feedspider.Feeder
	var feedTorlockTV feedspider.Feeder
	for _, r := range config.Config.Feed.TORLOCK {
		if r != nil {
			if r.ResourceType == types.VideoTypeTV {
				feedTorlockTV = feedspider.NewTorlock(r.Scheduling, r.ResourceType, r.MirrorSite)
			}
			log.Debug(r)
			if r.ResourceType == types.VideoTypeMovie {
				feedTorlockMovie = feedspider.NewTorlock(r.Scheduling, r.ResourceType, r.MirrorSite)
			}
			log.Debug(r)
		}
	}
	// MAGNETDL
	var feedMagnetdlMovie feedspider.Feeder
	var feedMagnetdlTV feedspider.Feeder
	for _, r := range config.Config.Feed.MagnetDL {
		if r != nil {
			if r.ResourceType == types.VideoTypeTV {
				feedMagnetdlTV = feedspider.NewMagnetdl(r.Scheduling, r.ResourceType, r.MirrorSite)
			}
			log.Debug(r)
			if r.ResourceType == types.VideoTypeMovie {
				feedMagnetdlMovie = feedspider.NewMagnetdl(r.Scheduling, r.ResourceType, r.MirrorSite)
			}
			log.Debug(r)
		}
	}

	feedTPBPIRATEPROXY := feedspider.NewTpbpirateproxy(config.Config.Feed.TPBPIRATEPROXY.Scheduling, config.Config.Feed.TPBPIRATEPROXY.MirrorSite)

	return optionFunc(func(ms *MovieSpider) {
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

// WithConfigFile
//
//	@Description: 初始化配置文件
//	@param configFile
//	@return Option
func WithConfigFile(configFile string) Option {
	config.InitConfig(configFile)
	model.NewMovieDB()
	return optionFunc(func(ms *MovieSpider) {})
}

// WithReport
//
//	@Description: 初始化 report
//	@return Option
func WithReport() Option {
	return optionFunc(func(ms *MovieSpider) {
		ms.report = job.NewReport("*/1 * * * *")
		go ms.report.Run()
	})
}

// WithDownload
//
//	@Description: 初始化下载器
//	@return Option
func WithDownload() Option {
	return optionFunc(func(ms *MovieSpider) {
		ms.download = download.NewDownloader(config.Config.Downloader.Scheduling)
		go ms.download.Run()
	})
}

// WithReleaseTimeJob
//
//	@Description: 初始化下载器
//	@return Option
func WithReleaseTimeJob() Option {
	if !config.Config.TG.Enable {
		return optionFunc(func(ms *MovieSpider) {
			log.Warn("未开启TG通知，无法运行 电影上线 通知job")
		})
	}
	return optionFunc(func(ms *MovieSpider) {
		ms.releaseTimeJob = job.NewReleaseTimeJob("")
		go ms.releaseTimeJob.Run()
	})
}
