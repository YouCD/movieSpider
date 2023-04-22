package feedSpider

import (
	"movieSpider/internal/types"
)

type Feeder interface {
	Crawler() ([]*types.FeedVideo, error)
	Run(chan *types.FeedVideo)
}

type FeederAbstractFactory interface {
	CreateFeeder(args ...interface{}) Feeder
}
