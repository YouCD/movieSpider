package feedspider

import (
	"fmt"
	"movieSpider/internal/types"
)

var (
	ErrNoFeedData = fmt.Errorf("no feed data")
)

type FactoryBTBT struct{}

func (f *FactoryBTBT) CreateFeeder(args ...interface{}) *Btbt {
	//nolint:forcetypeassert
	return &Btbt{
		BaseFeeder{
			web:        "btbt",
			url:        urlBtbt,
			scheduling: args[0].(string),
		},
	}
}

type FactoryEZTV struct{}

func (f *FactoryEZTV) CreateFeeder(args ...interface{}) *Eztv {
	var url string
	if len(args) == 0 {
		url = fmt.Sprintf("%s/%s", urlBaseEztv, urlRssURIEztv)
	} else if args[1] != "" {
		url = fmt.Sprintf("%s/%s", args[1], urlRssURIEztv)
	}
	//nolint:forcetypeassert
	return &Eztv{BaseFeeder{
		web:        "eztv",
		url:        url,
		scheduling: args[0].(string),
	}}
}

type FactoryGLODLS struct{}

//nolint:forcetypeassert
func (f *FactoryGLODLS) CreateFeeder(args ...interface{}) *Glodls {
	var url, urlBase string
	if len(args) == 0 {
		url = fmt.Sprintf("%s/%s", urlBaseGlodls, urlRssURIGlodls)
		urlBase = urlBaseGlodls
	} else if args[1] != "" {
		url = fmt.Sprintf("%s/%s", args[1], urlRssURIGlodls)
		urlBase = args[1].(string)
	}

	//nolint:exhaustruct
	return &Glodls{
		urlBase,
		BaseFeeder{
			web:        "glodls",
			url:        url,
			scheduling: args[0].(string),
		},
	}
}

type FactoryTGX struct{}

//nolint:forcetypeassert
func (f *FactoryTGX) CreateFeeder(args ...interface{}) *Tgx {
	var url string
	if len(args) == 0 {
		url = fmt.Sprintf("%s/%s", urlBaseTgx, urlRssURITgx)
	} else if args[1] != "" {
		url = fmt.Sprintf("%s/%s", args[1], urlRssURITgx)
	}

	return &Tgx{
		BaseFeeder{
			web:        "tgx",
			url:        url,
			scheduling: args[0].(string),
		},
	}
}

type FactoryTORLOCK struct{}

//nolint:forcetypeassert
func (f *FactoryTORLOCK) CreateFeeder(args ...interface{}) *Torlock {
	resourceType := args[1].(types.VideoType)

	urlBase := urlBaseTorlock
	if len(args) == 3 && args[2] != 0 {
		urlBase = args[2].(string)
	}

	url := fmt.Sprintf("%s/television/rss.xml", urlBase)
	if resourceType == types.VideoTypeMovie {
		url = fmt.Sprintf("%s/movies/rss.xml", urlBase)
	}
	return &Torlock{
		typ: resourceType,
		BaseFeeder: BaseFeeder{
			web:        "torlock",
			url:        url,
			scheduling: args[0].(string),
		},
	}
}

type FactoryMAGNETDL struct{}

//nolint:forcetypeassert
func (f *FactoryMAGNETDL) CreateFeeder(args ...interface{}) *Magnetdl {
	resourceType := args[1].(types.VideoType)
	var urlBase string

	if len(args) == 2 {
		urlBase = urlBaseMagnetdl
	} else if args[2] != "" {
		urlBase = args[2].(string)
	}

	url := fmt.Sprintf("%s/%s", urlBase, urlRssURITVMagnetdl)
	if resourceType == types.VideoTypeMovie {
		url = fmt.Sprintf("%s/%s", urlBase, urlRssURIMovieMagnetdl)
	}
	return &Magnetdl{
		resourceType,
		BaseFeeder{
			web:        "magnetdl",
			url:        url,
			scheduling: args[0].(string),
		},
	}
}

type FactoryTPBPIRATEPROXY struct{}

//nolint:forcetypeassert
func (f *FactoryTPBPIRATEPROXY) CreateFeeder(args ...interface{}) *Tpbpirateproxy {
	var url string
	if len(args) == 1 {
		url = fmt.Sprintf("%s/%s", urlBaseTpbpirateProxy, urlRssURITpbpirateProxy)
	} else if args[1] != "" {
		url = fmt.Sprintf("%s/%s", args[1], urlRssURITpbpirateProxy)
	}

	return &Tpbpirateproxy{
		BaseFeeder{
			web:        "tpbpirateproxy",
			url:        url,
			scheduling: args[0].(string),
		},
	}
}
