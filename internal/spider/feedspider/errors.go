package feedspider

import "errors"

var (
	ErrNoFeedData         = errors.New("no feed data")
	ErrFeedParseURL       = errors.New("feed url解析失败")
	ErrDownloadURLIsEmpty = errors.New("downloadURL is empty")
	ErrMagnetIsEmpty      = errors.New("magnet is empty")
)
