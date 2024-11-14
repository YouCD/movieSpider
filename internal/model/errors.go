package model

import "errors"

var (
	ErrFeedVideoIsNil            = errors.New("is nil")
	ErrFeedVideoExclude          = errors.New("exclude")
	ErrFeedVideoResolutionTooLow = errors.New("resolution too low")
)
