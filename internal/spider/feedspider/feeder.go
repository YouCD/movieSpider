package feedspider

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"os"
	"strings"
)

//nolint:inamedparam
type Feeder interface {
	//
	// Crawler
	//  @Description: 爬取视频
	//  @return []*types.FeedVideo
	//  @return error
	//
	Crawler() ([]*types.FeedVideo, error)
	//
	// Run
	//  @Description: cron运行
	//  @param chan
	//
	Run(chan *types.FeedVideo)
	//
	// Scheduling
	//  @Description: 获取cron时间表
	//  @return string
	//
	Scheduling() string
	//
	// WebName
	//  @Description: 获取网站名
	//  @return string
	//
	WebName() string
}

type BaseFeeder struct {
	web        string
	url        string
	scheduling string
}

func (b *BaseFeeder) Crawler() ([]*types.FeedVideo, error) {
	return nil, nil
}

func (b *BaseFeeder) Run(videosChan chan *types.FeedVideo) {
	if b.scheduling == "" {
		log.Errorf("%s Scheduling is null", b.web)
		os.Exit(1)
	}
	log.Infof("%s Scheduling is: [%s]", b.web, b.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(b.scheduling, func() {
		videos, err := b.Crawler()
		if err != nil {
			if errors.Is(err, ErrNoFeedData) {
				log.Warnf("%s: 没有feed数据, url: %s", strings.ToUpper(b.web), b.url)
				return
			}
			log.Error(err)
			return
		}
		for _, video := range videos {
			videosChan <- video
		}
	})
	c.Start()
}

func (b *BaseFeeder) Scheduling() string {
	return b.scheduling
}

func (b *BaseFeeder) WebName() string {
	return b.web
}

type FeederAbstractFactory interface {
	CreateFeeder(args ...interface{}) Feeder
}
