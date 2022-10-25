package feed

import (
	"movieSpider/internal/types"
)

type Feeder interface {
	Crawler() ([]*types.FeedVideo, error)
	Run()
}

type FeederAbstractFactory interface {
	CreateFeeder(args ...interface{}) Feeder
}
