package searchSpider

import "movieSpider/internal/types"

type SearchSpider interface {
	Search(name, resolution string) ([]*types.FeedVideo, error)
}
