package model

import "errors"

var (
	ErrNotMatchTorrentName = errors.New("torrent name not match")
	ErrFeedVideoIsNil      = errors.New("is nil")
	ErrFeedVideoExclude    = errors.New("exclude")
	ErrFeedVideoResolution = errors.New("resolution match")
	ErrFeedVideoYear       = errors.New("year match")
	ErrFeedVideoMovieMatch = errors.New("movie match")
	ErrFeedVideoTVMatch    = errors.New("tv match")
)
