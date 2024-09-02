package core

import (
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/feedspider"
	"movieSpider/internal/types"
	"strings"

	"github.com/youcd/toolkit/log"
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
	feedBTBT := feedspider.NewBtbt()

	// EZTV
	feedEZTV := feedspider.NewEztv()

	// GLODLS
	feedGLODLS := feedspider.NewGlodls()

	// TGX
	var TGXRss feedspider.Feeder
	var TGXDump feedspider.Feeder
	var TgxWeb feedspider.Feeder
	for _, tgx := range config.Config.Feed.TGX {
		switch strings.ToLower(tgx.Name) {
		case "rss":
			TGXRss = feedspider.NewTgx(tgx.Scheduling, tgx.Url, tgx.UseIPProxy)
		case "dump":
			TGXDump = feedspider.NewTgxDump(tgx.Scheduling, tgx.Url, tgx.UseIPProxy)
		case "web":
			TgxWeb = feedspider.NewTgxWeb(tgx.Scheduling, tgx.Url, tgx.UseIPProxy)
		}
	}

	// TORLOCK
	var feedTorlockMovie feedspider.Feeder
	var feedTorlockTV feedspider.Feeder
	for _, r := range config.Config.Feed.TORLOCK {
		if r != nil {
			if r.ResourceType == types.VideoTypeTV {
				feedTorlockTV = feedspider.NewTorlock(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
			}
			log.Debug(r)
			if r.ResourceType == types.VideoTypeMovie {
				feedTorlockMovie = feedspider.NewTorlock(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
			}
			log.Debug(r)
		}
	}
	// 1337x
	var feed1337xMovie feedspider.Feeder
	var feed1337xTV feedspider.Feeder
	for _, r := range config.Config.Feed.Web1337x {
		if r != nil {
			if r.ResourceType == types.VideoTypeTV {
				feed1337xTV = feedspider.NewWeb1337x(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
			}
			log.Debug(r)
			if r.ResourceType == types.VideoTypeMovie {
				feed1337xMovie = feedspider.NewWeb1337x(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
			}
			log.Debug(r)
		}
	}

	feedThePirateBay := feedspider.NewThePirateBay()

	return optionFunc(func(ms *MovieSpider) {
		ms.feeds = append(ms.feeds,
			feedBTBT,
			feedEZTV,
			feedGLODLS,
			TGXDump,
			TGXRss,
			TgxWeb,
			feedTorlockMovie,
			feedTorlockTV,
			feed1337xMovie,
			feed1337xTV,
			feedThePirateBay,
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
	return optionFunc(func(_ *MovieSpider) {})
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
	if config.Config.TG == nil {
		return optionFunc(func(_ *MovieSpider) {
			log.Warn("未开启TG通知，无法运行 电影上线 通知job")
		})
	}
	return optionFunc(func(ms *MovieSpider) {
		ms.releaseTimeJob = job.NewReleaseTimeJob("")
		go ms.releaseTimeJob.Run()
	})
}
