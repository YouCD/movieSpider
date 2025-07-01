package core

import (
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/feedspider"
	"movieSpider/internal/types"

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
	//// BTBT
	//feedBTBT := feedspider.NewBtbt()

	// EZTV
	feedEZTV := feedspider.NewEztv()

	// GLODLS
	feedGLODLS := feedspider.NewGlodls()
	// Knaben
	feedKnaben := feedspider.NewFeedKnaben()

	// TORLOCK
	feedTorlockTV, feedTorlockMovie := createFeederWithURLs(config.Config.Feed.TORLOCK, feedspider.NewTorlock)
	// 1337x
	feed1337xTV, feed1337xMovie := createFeederWithURLs(config.Config.Feed.Web1337x, feedspider.NewWeb1337x)

	// therarbg
	feedTheRarbg2TV, feedTheRarbg2Movie := createFeederWithURLs(config.Config.Feed.TheRarbg, feedspider.NewTheRarbg)

	feedThePirateBay := feedspider.NewThePirateBay()

	// Uindex
	uindexTv, uindexMovie := createFeederWithURLs(config.Config.Feed.Uindex, feedspider.NewUindex)

	// NewIlcorsaronero
	IlcorsaroneroTv, IlcorsaroneroMovie := createFeederWithURLs(config.Config.Feed.Uindex, feedspider.NewIlcorsaronero)
	return optionFunc(func(ms *MovieSpider) {
		ms.feeds = append(ms.feeds,
			feedEZTV,
			feedGLODLS,
			feedTorlockMovie,
			feedTorlockTV,
			feed1337xMovie,
			feed1337xTV,
			feedThePirateBay,
			feedKnaben,
			feedTheRarbg2TV,
			feedTheRarbg2Movie,
			uindexTv,
			uindexMovie,
			IlcorsaroneroTv,
			IlcorsaroneroMovie,
		)
		ms.feeds = append(ms.feeds, feeds...)
	})
}

type createFunc func(scheduling string, resourceType types.VideoType, siteURL string, useIPProxy bool) feedspider.Feeder

func createFeederWithURLs(urls []*config.BaseRT, create createFunc) (feedspider.Feeder, feedspider.Feeder) {
	var tv, movie feedspider.Feeder
	for _, r := range urls {
		if r.ResourceType == types.VideoTypeTV {
			tv = create(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
		}
		if r.ResourceType == types.VideoTypeMovie {
			movie = create(r.Scheduling, r.ResourceType, r.Url, r.UseIPProxy)
		}
	}
	return tv, movie
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

func WithDHT() Option {
	return optionFunc(func(ms *MovieSpider) {
		if config.Config.Global.DHTThread > 0 {
			ms.DHTThread = config.Config.Global.DHTThread
		}
	})
}
