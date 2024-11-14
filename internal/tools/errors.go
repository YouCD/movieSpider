package tools

import "errors"

var (
	ErrNotMatchTorrentName = errors.New("torrent name not match")
	ErrFeedVideoResolution = errors.New("resolution match")
	ErrFeedVideoYear       = errors.New("year match")
	ErrFeedVideoMovieMatch = errors.New("movie match")
)
