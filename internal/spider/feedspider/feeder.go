package feedspider

import (
	"movieSpider/internal/types"
)

//nolint:inamedparam
type Feeder interface {
	Crawler() ([]*types.FeedVideo, error)
	Run(chan *types.FeedVideo)
}

type FeederAbstractFactory interface {
	CreateFeeder(args ...interface{}) Feeder
}
