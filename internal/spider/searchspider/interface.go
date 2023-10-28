package searchspider

import "movieSpider/internal/types"

type SearchSpider interface {
	Search() (videos []*types.FeedVideo, err error)
}
