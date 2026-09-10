package core

import (
	"context"

	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/feedspider"
	tmdbspider "movieSpider/internal/spider/tmdb"
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

// WithFeeds 初始化feeds
func WithFeeds(feeds ...feedspider.Feeder) Option {
	// EZTV
	feedEZTV := feedspider.NewEztv()

	// Knaben
	feedKnaben := feedspider.NewFeedKnaben()

	// TORLOCK
	feedTorlockTV, feedTorlockMovie := createFeederWithURLs(config.Config.Feed.TORLOCK, feedspider.NewTorlock)
	// 1337x
	// feed1337xTV, feed1337xMovie := createFeederWithURLs(config.Config.Feed.Web1337x, feedspider.NewWeb1337x)

	// therarbg
	feedTheRarbg2TV, feedTheRarbg2Movie := createFeederWithURLs(config.Config.Feed.TheRarbg, feedspider.NewTheRarbg)

	feedThePirateBay := feedspider.NewThePirateBay()

	// Uindex
	uindexTv, uindexMovie := createFeederWithURLs(config.Config.Feed.Uindex, feedspider.NewUindex)

	feedYts := feedspider.NewYts()
	return optionFunc(func(ms *MovieSpider) {
		ms.feeds = append(
			ms.feeds,
			feedEZTV,
			feedTorlockMovie,
			feedTorlockTV,
			// feed1337xMovie,
			// feed1337xTV,
			feedThePirateBay,
			feedKnaben,
			feedTheRarbg2TV,
			feedTheRarbg2Movie,
			uindexTv,
			uindexMovie,
			feedYts,
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

// WithConfigFile 初始化配置文件
func WithConfigFile(configFile string) Option {
	config.InitConfig(configFile)
	model.NewMovieDB()
	return optionFunc(func(_ *MovieSpider) {})
}

// WithReport 初始化 report
func WithReport() Option {
	return optionFunc(func(ms *MovieSpider) {
		ms.report = job.NewReport("*/1 * * * *")
		go ms.report.Run()
	})
}

// WithDownload 初始化下载器
func WithDownload(ctx context.Context) Option {
	return optionFunc(func(ms *MovieSpider) {
		ms.download = download.NewDownloader(config.Config.Downloader.Scheduling)
		go ms.download.Run(ctx)
	})
}

// WithReleaseTimeJob 初始化下载器
func WithReleaseTimeJob(ctx context.Context) Option {
	if config.Config.TG == nil {
		return optionFunc(func(_ *MovieSpider) {
			log.WithCtx(context.Background()).Warn("未开启TG通知，无法运行 电影上线 通知job")
		})
	}
	return optionFunc(func(ms *MovieSpider) {
		ms.releaseTimeJob = job.NewReleaseTimeJob("")
		go ms.releaseTimeJob.Run(ctx)
	})
}

func WithDHT() Option {
	return optionFunc(func(ms *MovieSpider) {
		if config.Config.Global.DHTThread > 0 {
			ms.DHTThread = config.Config.Global.DHTThread
		}
	})
}

// WithTMDBSpider 初始化TMDB爬虫
func WithTMDBSpider(accountID int, bearerToken string) Option {
	return optionFunc(func(ms *MovieSpider) {
		tmdbSpider, err := tmdbspider.NewTMDBSpider(accountID, bearerToken)
		if err != nil {
			log.WithCtx(context.Background()).Errorf("创建TMDB爬虫失败: %s", err)
			return
		}

		ms.spiders = append(ms.spiders, tmdbSpider)
		log.WithCtx(context.Background()).Info("TMDB爬虫已初始化，每10分钟执行一次")
	})
}
