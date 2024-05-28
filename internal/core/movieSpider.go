package core

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
	"movieSpider/internal/bot"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/download"
	"movieSpider/internal/job"
	"movieSpider/internal/spider"
	"movieSpider/internal/spider/douban"
	"movieSpider/internal/spider/feedspider"
	"os"
	"strings"
)

type MovieSpider struct {
	feeds          []feedspider.Feeder
	download       *download.Download
	report         *job.Report
	bot            *bot.TGBot
	spiders        []spider.Spider
	releaseTimeJob *job.ReleaseTimeJob
}

//nolint:gochecknoglobals
var (
	ms = new(MovieSpider)
)

// NewMovieSpider
//
//	@Description: 初始化movieSpider
//	@param options
//	@return *MovieSpider
func NewMovieSpider(options ...Option) *MovieSpider {
	for _, option := range options {
		option.apply(ms)
	}
	return ms
}

// RunWithFeed
//
//	@Description: 运行feed
//	@receiver m
func (m *MovieSpider) RunWithFeed() {
	for _, feeder := range m.feeds {
		go func(feeder feedspider.Feeder) {
			if feeder.Scheduling() == "" {
				log.Errorf("%s Scheduling is null", feeder.WebName())
				os.Exit(1)
			}
			log.Infof("%s Scheduling is: [%s]", feeder.WebName(), feeder.Scheduling())
			c := cron.New()
			_, _ = c.AddFunc(feeder.Scheduling(), func() {
				videos, err := feeder.Crawler()
				if err != nil {
					if errors.Is(err, feedspider.ErrNoFeedData) {
						log.Warnf("%s: 没有feed数据, url: %s", strings.ToUpper(feeder.WebName()), feeder.URL())
						return
					}
					log.Error(err)
					return
				}
				for _, video := range videos {
					bus.FeedVideoChan <- video
				}
			})
			c.Start()
		}(feeder)
	}
}

// RunWithTGBot
//
//	@Description: 运行tgbot
//	@receiver m
func (m *MovieSpider) RunWithTGBot() {
	if config.Config.TG != nil {
		ms.bot = bot.NewTgBot(config.Config.TG.BotToken, config.Config.TG.TgIDs)
		go ms.bot.StartBot()
	}
}

// RunWithFeedSpider
//
//	@Description: 运行feedSpider
//	@receiver m
func (m *MovieSpider) RunWithFeedSpider() {
	// Spider
	m.spiders = append(m.spiders, douban.NewSpiderDouBan(config.Config.DouBan)...)
	for _, s := range m.spiders {
		go func(spider spider.Spider) {
			spider.Run()
		}(s)
	}
}
