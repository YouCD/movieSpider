package searchSpider

import (
	"fmt"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"net/url"
	"os"
)

type FactoryBt4g struct{}

func (f *FactoryBt4g) CreateFeeder(args ...interface{}) *bt4g {
	name := args[0].(string)
	resolution := args[1].(types.Resolution)
	parse, err := url.Parse(urlBt4g)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	bUrl := fmt.Sprintf("%s://%s/search/%s/bysize/1?page=rss", parse.Scheme, parse.Host, name)
	return &bt4g{url: bUrl, resolution: resolution, web: "bt4g"}
}

type FactoryKNABEN struct{}

func (f *FactoryKNABEN) CreateFeeder(args ...interface{}) *knaben {
	name := args[0].(string)
	resolution := args[1].(types.Resolution)
	parse, err := url.Parse(urlKnaben)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	strData := url.QueryEscape(name)

	kUrl := fmt.Sprintf("%s://%s/%s", parse.Scheme, parse.Host, strData)

	return &knaben{url: kUrl, resolution: resolution, web: "knaben"}
}
