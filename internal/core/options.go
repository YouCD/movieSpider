package core

import (
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	feed2 "movieSpider/internal/spider/feedspider"
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
func WithFeeds(feeds ...feed2.Feeder) Option {
	// BTBT
	feedBTBT := feed2.NewBtbt(config.BTBT.Scheduling)

	// EZTV
	feedEZTV := feed2.NewEztv(config.EZTV.Scheduling, config.EZTV.MirrorSite)

	// GLODLS
	feedGLODLS := feed2.NewGlodls(config.GLODLS.Scheduling, config.GLODLS.MirrorSite)

	// TGX
	feedTGXS := feed2.NewTgx(config.TGX.Scheduling, config.TGX.MirrorSite)

	// TORLOCK
	var feedTorlockMovie feed2.Feeder
	var feedTorlockTV feed2.Feeder
	for _, r := range config.TORLOCK {
		if r != nil {
			if r.Typ == types.VideoTypeTV {
				feedTorlockTV = feed2.NewTorlock(r.Scheduling, r.Typ, r.MirrorSite)
			}
			log.Debug(r)
			if r.Typ == types.VideoTypeMovie {
				feedTorlockMovie = feed2.NewTorlock(r.Scheduling, r.Typ, r.MirrorSite)
			}
			log.Debug(r)
		}
	}
	// MAGNETDL
	var feedMagnetdlMovie feed2.Feeder
	var feedMagnetdlTV feed2.Feeder
	for _, r := range config.MAGNETDL {
		if r != nil {
			if r.Typ == types.VideoTypeTV {
				feedMagnetdlTV = feed2.NewMagnetdl(r.Scheduling, r.Typ, r.MirrorSite)
			}
			log.Debug(r)
			if r.Typ == types.VideoTypeMovie {
				feedMagnetdlMovie = feed2.NewMagnetdl(r.Scheduling, r.Typ, r.MirrorSite)
			}
			log.Debug(r)
		}
	}

	feedTPBPIRATEPROXY := feed2.NewTpbpirateproxy(config.TPBPIRATEPROXY.Scheduling, config.TPBPIRATEPROXY.MirrorSite)

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
		ms.download = download.NewDownloader(config.Downloader.Scheduling)
		go ms.download.Run()
	})
}

// WithReleaseTimeJob
//
//	@Description: 初始化下载器
//	@return Option
func WithReleaseTimeJob() Option {
	if !config.TG.Enable {
		return optionFunc(func(ms *MovieSpider) {
			log.Warn("未开启TG通知，无法运行 电影上线 通知job")
		})
	}
	return optionFunc(func(ms *MovieSpider) {
		ms.releaseTimeJob = job.NewReleaseTimeJob("")
		go ms.releaseTimeJob.Run()
	})
}
