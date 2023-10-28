package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"os"
)

//nolint:tagliatelle
type btbt struct {
	Scheduling string `json:"Scheduling"`
}

//nolint:tagliatelle
type tgx struct {
	Scheduling string `json:"Scheduling"`
}

//nolint:tagliatelle
type torlock struct {
	Scheduling   string `json:"Scheduling"`
	ResourceType string `json:"ResourceType"`
	Typ          types.VideoType
}

//nolint:tagliatelle
type eztv struct {
	Scheduling string `json:"Scheduling"`
}

//nolint:tagliatelle
type glodls struct {
	Scheduling string `json:"Scheduling"`
}

//nolint:tagliatelle
type tpbpirateproxy struct {
	Scheduling string `json:"Scheduling"`
}

//nolint:tagliatelle
type magnetdl struct {
	Scheduling   string `json:"Scheduling"`
	ResourceType string `json:"ResourceType"`
	Typ          types.VideoType
}

//nolint:tagliatelle
type downloader struct {
	Scheduling string `json:"Scheduling"`
	Aria2Label string `json:"Aria2Label"`
}

//nolint:gochecknoglobals
var (
	Global         *global
	Aria2cList     []aria2
	TG             = new(tg)
	MySQL          *mysql
	DouBanList     *DouBan
	BTBT           *btbt
	EZTV           *eztv
	GLODLS         *glodls
	TPBPIRATEPROXY *tpbpirateproxy
	TGX            *tgx
	TORLOCK        []*torlock
	MAGNETDL       []*magnetdl
	Downloader     *downloader
	ProxyPool      string
	TmDB           *tmDB
	ExcludeWords   []string
)

type global struct {
	LogLevel string
	Report   bool
}

type aria2 struct {
	URL   string
	Token string
	Label string
}

type tg struct {
	BotToken string
	TgIDs    []int
	Proxy    struct {
		URL    string
		Enable bool
	}
	Enable bool
}

type mysql struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}
type DouBan struct {
	DouBanList []*DouBan
	Scheduling string
	URL        string
}

//nolint:tagliatelle
type tmDB struct {
	Scheduling string `json:"Scheduling"`
	APIKey     string
}

//nolint:forbidigo,gosimple,gocognit,gocyclo,maintidx
func InitConfig(config string) {
	v := viper.New()
	v.SetConfigType("yaml")

	fmt.Printf("config file is %s.\n", config)
	b, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Printf("配置文件读取错误,err:%s\n", err.Error())
		os.Exit(1)
	}

	err = v.ReadConfig(bytes.NewReader(b))
	if err != nil {
		fmt.Printf("配置文件错误.")
		os.Exit(1)
	}

	err = v.UnmarshalKey("Global", &Global)
	if err != nil {
		fmt.Println("读取Global配置错误")
		os.Exit(1)
	}
	if Global == nil {
		fmt.Println("配置 Global is nil")
		os.Exit(1)
	}

	if Global.LogLevel == "debug" {
		fmt.Println(string(b))
	}
	log.NewLogger(Global.LogLevel)

	err = v.UnmarshalKey("Aria2cList", &Aria2cList)
	if err != nil {
		fmt.Println("读取Aria2cList配置错误")
		os.Exit(1)
	}
	for index, aria := range Aria2cList {
		Aria2cList[index].URL = aria.URL + "/jsonrpc"
	}
	if Aria2cList == nil {
		fmt.Println("配置 Aria2cList is nil")
		os.Exit(1)
	}

	err = v.UnmarshalKey("TG", &TG)
	if err != nil {
		fmt.Println("读取TG配置错误")
		os.Exit(1)
	}
	if TG.BotToken != "" && TG.TgIDs != nil && TG.Proxy.URL != "" {
		TG.Enable = true
	}

	err = v.UnmarshalKey("MySQL", &MySQL)
	if err != nil {
		fmt.Println("读取MySQL配置错误")
		os.Exit(1)
	}
	if MySQL == nil {
		fmt.Println("配置 MySQL is nil")
		os.Exit(1)
	}
	if err = v.UnmarshalKey("DouBan", &DouBanList); err != nil {
		fmt.Println("读取DouBan配置错误")
		os.Exit(1)
	}

	if DouBanList == nil {
		fmt.Println("配置 DouBanList is nil")
		os.Exit(1)
	}
	// btbt
	if err = v.UnmarshalKey("Feed.BTBT", &BTBT); err != nil {
		fmt.Println("读取BTBT配置错误")
		os.Exit(1)
	}
	if BTBT == nil {
		fmt.Println("配置 BTBT is nil")
		os.Exit(1)
	}
	if err = v.UnmarshalKey("Feed.EZTV", &EZTV); err != nil {
		fmt.Println("读取EZTV配置错误")
		os.Exit(1)
	}
	if EZTV == nil {
		fmt.Println("配置 EZTV is nil")
		os.Exit(1)
	}
	if err = v.UnmarshalKey("Feed.TPBPIRATEPROXY", &TPBPIRATEPROXY); err != nil {
		fmt.Println("读取TPBPIRATEPROXY配置错误")
		os.Exit(1)
	}
	if TPBPIRATEPROXY == nil {
		fmt.Println("配置 TPBPIRATEPROXY is nil")
		os.Exit(1)
	}
	if err = v.UnmarshalKey("Feed.GLODLS", &GLODLS); err != nil {
		fmt.Println("读取GLODLS配置错误")
		os.Exit(1)
	}
	if GLODLS == nil {
		fmt.Println("配置 GLODLS is nil")
		os.Exit(1)
	}

	if err = v.UnmarshalKey("Feed.TGX", &TGX); err != nil {
		fmt.Println("读取TGX配置错误")
		os.Exit(1)
	}
	if TGX == nil {
		fmt.Println("配置 TGX is nil")
		os.Exit(1)
	}

	if err = v.UnmarshalKey("Feed.TORLOCK", &TORLOCK); err != nil {
		fmt.Println("读取TORLOCK配置错误")
		os.Exit(1)
	}
	if TORLOCK == nil {
		fmt.Println("配置 TORLOCK is nil")
		os.Exit(1)
	}
	for _, v := range TORLOCK {
		//nolint:exhaustive
		switch types.Convert2VideoType(v.ResourceType) {
		case types.VideoTypeMovie:
			v.Typ = types.VideoTypeMovie
		case types.VideoTypeTV:
			v.Typ = types.VideoTypeTV
		default:
			v.Typ = types.VideoTypeTV
		}
	}
	if err = v.UnmarshalKey("Feed.MAGNETDL", &MAGNETDL); err != nil {
		fmt.Println("读取MAGNETDL配置错误")
		os.Exit(1)
	}
	if MAGNETDL == nil {
		fmt.Println("配置 MAGNETDL is nil")
		os.Exit(1)
	}
	for _, v := range MAGNETDL {
		switch v.ResourceType {
		case "movie":
			v.Typ = types.VideoTypeMovie
		case "tv":
			v.Typ = types.VideoTypeTV
		default:
			v.Typ = types.VideoTypeTV
		}
	}

	if err = v.UnmarshalKey("Feed.ProxyPool", &ProxyPool); err != nil {
		fmt.Println("读取Feed.ProxyPool配置错误")
		os.Exit(1)
	}
	if ProxyPool == "" {
		fmt.Println("配置 ProxyPool is null")
		os.Exit(1)
	}
	if err = v.UnmarshalKey("Downloader", &Downloader); err != nil {
		fmt.Println("读取Downloader配置错误")
		os.Exit(1)
	}
	if Downloader == nil {
		fmt.Println("配置 Downloader is nil")
		os.Exit(1)
	}

	if err = v.UnmarshalKey("TmDB", &TmDB); err != nil {
		fmt.Println("读取TmDBL配置错误")
		os.Exit(1)
	}
	if err = v.UnmarshalKey("ExcludeWords", &ExcludeWords); err != nil {
		fmt.Println("读取 ExcludeWords 配置错误")
		os.Exit(1)
	}

	return
}
