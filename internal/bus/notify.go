package bus

import "movieSpider/internal/types"

var (
	DownloadNotifyChan = make(chan *types.DownloadNotifyVideo)
	FeedVideoChan      = make(chan *types.FeedVideo)
	DatePublishedChan  = make(chan *types.DouBanVideo)
)
