package bus

import "movieSpider/internal/types"

var (
	DownloadNotifyChan = make(chan *types.DouBanVideo)
	FeedVideoChan      = make(chan *types.FeedVideo)
	DatePublishedChan  = make(chan *types.DouBanVideo)
)
