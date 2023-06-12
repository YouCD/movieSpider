package bus

import "movieSpider/internal/types"

var (
	NotifyChan        = make(chan string)
	FeedVideoChan     = make(chan *types.FeedVideo)
	DatePublishedChan = make(chan *types.DouBanVideo)
)
