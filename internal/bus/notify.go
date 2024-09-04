package bus

import "movieSpider/internal/types"

//nolint:gochecknoglobals
var (
	DownloadNotifyChan = make(chan *types.DownloadNotifyVideo)
	FeedVideoChan      = make(chan *types.FeedVideoBase)
	DatePublishedChan  = make(chan *types.DouBanVideo)
)
