package searchspider

import (
	"fmt"
	"movieSpider/internal/types"
	"net/url"
	"os"

	"github.com/youcd/toolkit/log"
)

type FactoryBt4g struct{}

//nolint:forcetypeassert
func (f *FactoryBt4g) CreateFeeder(args ...interface{}) *BT4g {
	name := args[0].(string)
	resolution := args[1].(types.Resolution)
	parse, err := url.Parse(urlBt4g)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	bURL := fmt.Sprintf("%s://%s/search/%s/bysize/1?page=rss", parse.Scheme, parse.Host, name)
	return &BT4g{url: bURL, resolution: resolution, web: "BT4g"}
}
