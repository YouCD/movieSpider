package model

import "errors"

var (
	ErrNotMatchTorrentName = errors.New("torrent name not match")
	ErrFeedVideoIsNil      = errors.New("feedVideo is nil")
	ErrFeedVideoExclude    = errors.New("feedVideo exclude")
	ErrFeedVideoResolution = errors.New("feedVideo resolution match")
	ErrFeedVideoYear       = errors.New("feedVideo year match")
)
