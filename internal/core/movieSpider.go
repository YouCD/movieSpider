package core

import (
	"context"
	"errors"
	"movieSpider/internal/bot"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	dhtc_client "movieSpider/internal/dhtclient"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/spider"
	"movieSpider/internal/spider/feedspider"
	"os"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
)

type MovieSpider struct {
	feeds          []feedspider.Feeder
	download       *download.Download
	report         *job.Report
	bot            *bot.TGBot
	spiders        []spider.Spider
	releaseTimeJob *job.ReleaseTimeJob
	DHTThread      int
	wg             sync.WaitGroup
	cronJobs       []*cron.Cron
}

//nolint:gochecknoglobals
var ms = new(MovieSpider)

// NewMovieSpider 初始化movieSpider
func NewMovieSpider(options ...Option) *MovieSpider {
	for _, option := range options {
		option.apply(ms)
	}
	return ms
}

func (m *MovieSpider) Start(ctx context.Context) {
	if config.Config.TG != nil {
		ms.bot = bot.NewTgBot(config.Config.TG.BotToken, config.Config.TG.TgIDs)
		go ms.bot.StartBot(ctx)
	}
	if m.DHTThread > 0 {
		go dhtc_client.Boot(m.DHTThread)
	}

	m.startFeed(ctx)
	m.startSpider(ctx)
}

// Stop 优雅关闭
func (m *MovieSpider) Stop() {
	// 停止所有 cron 任务
	for _, c := range m.cronJobs {
		c.Stop()
	}
	// 等待所有 goroutine 完成
	m.wg.Wait()
}

// startFeed 运行feed
func (m *MovieSpider) startFeed(ctx context.Context) {
	for _, feeder := range m.feeds {
		m.wg.Add(1)
		go func(feeder feedspider.Feeder) {
			defer m.wg.Done()

			if feeder.Scheduling() == "" {
				log.WithCtx(ctx).Errorf("%s Scheduling is null", feeder.WebName())
				os.Exit(1)
			}
			log.WithCtx(ctx).Infof("%s Scheduling is: [%s]", feeder.WebName(), feeder.Scheduling())

			c := cron.New()
			m.cronJobs = append(m.cronJobs, c)

			_, _ = c.AddFunc(feeder.Scheduling(), func() {
				m.processFeed(ctx, feeder)
			})
			c.Start()
		}(feeder)
	}
}

// processFeed 处理单个 feed
func (m *MovieSpider) processFeed(ctx context.Context, feeder feedspider.Feeder) {
	ctx = log.SetRequestId(ctx)
	videos, err := feeder.Crawler(ctx)
	if err != nil {
		if errors.Is(err, feedspider.ErrNoFeedData) {
			log.WithCtx(ctx).Warnf("%s: 没有feed数据, url: %s", strings.ToUpper(feeder.WebName()), feeder.URL())
			return
		}
		log.WithCtx(ctx).Errorf("web: %s, err: %s", feeder.WebName(), err)
		return
	}
	if len(videos) == 0 {
		log.WithCtx(ctx).Warnf("web: %s, url: %s, videos is empty", feeder.WebName(), feeder.URL())
		return
	}
	log.WithCtx(ctx).Infof("web: %s, cont: %d", feeder.WebName(), len(videos))
	for _, video := range videos {
		if video.Magnet == "" {
			log.WithCtx(ctx).Warnf("web: %s, url: %s, Magnet is empty", feeder.WebName(), feeder.URL())
			continue
		}
		bus.FeedVideoChan <- video
	}
}

// startSpider 运行 Spider
func (m *MovieSpider) startSpider(ctx context.Context) {
	// Spider
	for _, s := range m.spiders {
		m.wg.Add(1)
		go func(spider spider.Spider) {
			defer m.wg.Done()
			spider.Run(ctx)
		}(s)
	}
}
