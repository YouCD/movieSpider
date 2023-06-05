package feedSpider

import (
	types2 "movieSpider/internal/types"
	"net/http"
)

type FactoryBTBT struct{}

func (f *FactoryBTBT) CreateFeeder(args ...interface{}) Feeder {
	scheduling := args[0].(string)
	return &btbt{
		urlBtbt,
		scheduling,
	}
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
		&http.Client{},
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
