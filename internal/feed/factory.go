package feed

import (
	"fmt"
	"movieSpider/internal/log"
	types2 "movieSpider/internal/types"
	"net/http"
	"net/url"
	"os"
)

type FactoryBTBT struct{}

func (f *FactoryBTBT) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &btbt{
		urlBtbt,
		scheduling,
	}
}

type FactoryBt4g struct{}

func (f *FactoryBt4g) CreateFeeder(args ...interface{}) Feeder {
	name := args[0].(string)
	resolution := args[1].(types2.Resolution)
	parse, err := url.Parse(urlBt4g)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	bUrl := fmt.Sprintf("%s://%s/search/%s/bysize/1?page=rss", parse.Scheme, parse.Host, name)
	return &bt4g{url: bUrl, resolution: resolution, web: "bt4g"}
}

type FactoryEZTV struct{}

func (f *FactoryEZTV) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &eztv{
		scheduling,
		urlEztv,
		"eztv",
	}
}

type FactoryGLODLS struct{}

func (f *FactoryGLODLS) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &glodls{
		urlGlodls,
		scheduling,
		"glodls",
	}
}

type FactoryKNABEN struct{}

func (f *FactoryKNABEN) CreateFeeder(args ...interface{}) Feeder {
	name := args[0].(string)
	resolution := args[1].(types2.Resolution)
	parse, err := url.Parse(urlKnaben)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	strData := url.QueryEscape(name)

	kUrl := fmt.Sprintf("%s://%s/%s", parse.Scheme, parse.Host, strData)

	return &knaben{url: kUrl, resolution: resolution, web: "knaben"}
}

type FactoryRARBG struct{}

func (f *FactoryRARBG) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	resourceType := args[1].(types2.Resource)

	if resourceType == types2.ResourceMovie {
		return &rarbg{
			resourceType,
			"rarbg",
			scheduling,
			&http.Client{},
		}
	} else {
		return &rarbg{
			resourceType,
			"rarbg",
			scheduling,
			&http.Client{},
		}
	}

}

type FactoryTGX struct{}

func (f *FactoryTGX) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &tgx{
		scheduling: scheduling,
		url:        urlTgx,
		web:        "tgx",
	}
}

type FactoryTORLOCK struct{}

func (f *FactoryTORLOCK) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	resourceType := args[1].(types2.Resource)

	if resourceType == types2.ResourceMovie {
		return &torlock{
			resourceType,
			"torlock",
			scheduling,
		}
	} else {
		return &torlock{
			resourceType,
			"torlock",
			scheduling,
		}
	}

}

type FactoryMAGNETDL struct{}

func (f *FactoryMAGNETDL) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	resourceType := args[1].(types2.Resource)

	if resourceType == types2.ResourceMovie {
		return &torlock{
			resourceType,
			"magnetdl",
			scheduling,
		}
	} else {
		return &torlock{
			resourceType,
			"magnetdl",
			scheduling,
		}
	}

}

type FactoryTPBPIRATEPROXY struct{}

func (f *FactoryTPBPIRATEPROXY) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &tpbpirateproxy{
		scheduling,
		"tpbpirateproxy",
	}
}
