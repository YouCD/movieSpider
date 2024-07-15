package feedspider

import (
	"errors"
	"movieSpider/internal/types"
)

var (
	ErrNoFeedData   = errors.New("no feed data")
	ErrFeedParseURL = errors.New("feed url解析失败")
)

//nolint:inamedparam
type Feeder interface {
	Scheduling() string
	WebName() string
	URL() string
	Crawler() ([]*types.FeedVideo, error)
}
type Crawler func() ([]*types.FeedVideo, error)

type BaseFeeder struct {
	web        string
	url        string
	scheduling string
}

func (b *BaseFeeder) Crawler() ([]*types.FeedVideo, error) {
	return nil, nil
}

func (b *BaseFeeder) URL() string {
	return b.url
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
