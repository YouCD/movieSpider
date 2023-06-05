package config

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
	"io/ioutil"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"os"
)

type btbt struct {
	Scheduling string `json:"Scheduling"`
}
type tgx struct {
	Scheduling string `json:"Scheduling"`
}
type torlock struct {
	Scheduling   string `json:"Scheduling"`
	ResourceType string `json:"ResourceType"`
	Typ          types.Resource
}

type eztv struct {
	Scheduling string `json:"Scheduling"`
}
type glodls struct {
	Scheduling string `json:"Scheduling"`
}
type tpbpirateproxy struct {
	Scheduling string `json:"Scheduling"`
}

type magnetdl struct {
	Scheduling   string `json:"Scheduling"`
	ResourceType string `json:"ResourceType"`
	Typ          types.Resource
}

type downloader struct {
	Scheduling string `json:"Scheduling"`
	Aria2Label string `json:"Aria2Label"`
}

var (
	Global         *global
	Aria2cList     []aria2
	TG             = new(tg)
	MySQL          *mysql
	DouBan         *douban
	BTBT           *btbt
	EZTV           *eztv
	GLODLS         *glodls
	TPBPIRATEPROXY *tpbpirateproxy
	TGX            *tgx
	TORLOCK        []*torlock
	MAGNETDL       []*magnetdl
	Downloader     *downloader
	ProxyPool      string
)

type global struct {
	LogLevel string
	Report   bool
}

type aria2 struct {
	Url   string
	Token string
	Label string
}

type tg struct {
	BotToken string
	TgIDs    []int
	Proxy    struct {
		Url    string
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
type douban struct {
	DoubanUrl  string
	Scheduling string
	Cookie     string
}

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
		Aria2cList[index].Url = aria.Url + "/jsonrpc"
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
	if TG.BotToken != "" && TG.TgIDs != nil && TG.Proxy.Url != "" {
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
	if err = v.UnmarshalKey("Douban", &DouBan); err != nil {
		fmt.Println("读取Douban配置错误")
		os.Exit(1)
	}
	if !govalidator.IsURL(DouBan.DoubanUrl) {
		DouBan.DoubanUrl = ""
	}
	if DouBan == nil {
		fmt.Println("配置 DouBan is nil")
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
		switch v.ResourceType {
		case types.VideoTypeMovie:
			v.Typ = types.ResourceMovie
		case types.VideoTypeTV:
			v.Typ = types.ResourceTV
		default:
			v.Typ = types.ResourceTV
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
			v.Typ = types.ResourceMovie
		case "tv":
			v.Typ = types.ResourceTV
		default:
			v.Typ = types.ResourceTV
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
	return

}
