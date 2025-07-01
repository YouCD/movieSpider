package model

import "errors"

var (
	ErrFeedVideoIsNil            = errors.New("is nil")
	ErrFeedVideoExclude          = errors.New("exclude")
	ErrFeedVideoExist            = errors.New("data exist")
	ErrFeedVideoResolutionTooLow = errors.New("resolution too low")
)
