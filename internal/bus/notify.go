package bus

import "movieSpider/internal/types"

//nolint:gochecknoglobals
var (
	DownloadNotifyChan = make(chan *types.DownloadNotifyVideo)
	FeedVideoChan      = make(chan *types.FeedVideo)
	DatePublishedChan  = make(chan *types.DouBanVideo)
)
