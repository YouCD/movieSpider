package feedspider

import (
	"movieSpider/internal/types"
	"net/http"
)

type FactoryBTBT struct{}

func (f *FactoryBTBT) CreateFeeder(args ...interface{}) Feeder {
	//nolint:forcetypeassert
	scheduling := args[0].(string)
	return &btbt{
		urlBtbt,
		scheduling,
	}
}

type FactoryEZTV struct{}

func (f *FactoryEZTV) CreateFeeder(args ...interface{}) Feeder {
	//nolint:forcetypeassert
	scheduling := args[0].(string)
	return &eztv{
		scheduling,
		urlEztv,
		"eztv",
	}
}

type FactoryGLODLS struct{}

func (f *FactoryGLODLS) CreateFeeder(args ...interface{}) Feeder {
	//nolint:forcetypeassert
	scheduling := args[0].(string)
	//nolint:exhaustruct
	return &glodls{
		urlGlodls,
		scheduling,
		"glodls",
		&http.Client{},
	}
}

type FactoryTGX struct{}

//nolint:forcetypeassert
func (f *FactoryTGX) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &tgx{
		scheduling: scheduling,
		url:        urlTgx,
		web:        "tgx",
	}
}

type FactoryTORLOCK struct{}

//nolint:forcetypeassert
func (f *FactoryTORLOCK) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	resourceType := args[1].(types.VideoType)

	if resourceType == types.VideoTypeMovie {
		return &torlock{
			resourceType,
			"torlock",
			scheduling,
		}
	}
	return &torlock{
		resourceType,
		"torlock",
		scheduling,
	}
}

type FactoryMAGNETDL struct{}

//nolint:forcetypeassert
func (f *FactoryMAGNETDL) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	resourceType := args[1].(types.VideoType)

	if resourceType == types.VideoTypeMovie {
		return &torlock{
			resourceType,
			"magnetdl",
			scheduling,
		}
	}
	return &torlock{
		resourceType,
		"magnetdl",
		scheduling,
	}
}

type FactoryTPBPIRATEPROXY struct{}

//nolint:forcetypeassert
func (f *FactoryTPBPIRATEPROXY) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &tpbpirateproxy{
		scheduling,
		"tpbpirateproxy",
	}
}
