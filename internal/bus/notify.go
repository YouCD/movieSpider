package bus

import "movieSpider/internal/types"

//nolint:gochecknoglobals
var (
	// 有缓冲通道，避免单个消费者缓慢时阻塞所有生产者 goroutine
	DownloadNotifyChan = make(chan *types.DownloadNotifyVideo, 1024)
	FeedVideoChan      = make(chan *types.FeedVideoBase, 4096)
	DatePublishedChan  = make(chan *types.TMDBVideo, 1024)
)
