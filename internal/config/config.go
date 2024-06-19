package config

import (
	"bytes"
	"fmt"
	"movieSpider/internal/types"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/youcd/toolkit/log"
)

//nolint:tagliatelle
type btbt struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
}

//nolint:tagliatelle
type tgx struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	MirrorSite string `json:"MirrorSite,omitempty" yaml:"MirrorSite,omitempty" validate:"omitempty,http_url"`
}

//nolint:tagliatelle
type torlock struct {
	Scheduling   string          `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	MirrorSite   string          `json:"MirrorSite,omitempty" yaml:"MirrorSite,omitempty" validate:"omitempty,http_url"`
	ResourceType types.VideoType `json:"ResourceType" yaml:"ResourceType" validate:"required,oneof=movie tv"`
}

//nolint:tagliatelle
type eztv struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	MirrorSite string `json:"MirrorSite,omitempty" yaml:"MirrorSite,omitempty" validate:"omitempty,http_url"`
}

//nolint:tagliatelle
type glodls struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	MirrorSite string `json:"MirrorSite,omitempty" yaml:"MirrorSite,omitempty" validate:"omitempty,http_url"`
}

//nolint:tagliatelle
type tpbpirateproxy struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	MirrorSite string `json:"MirrorSite,omitempty" yaml:"MirrorSite,omitempty" validate:"omitempty,http_url"`
}

//nolint:tagliatelle
type magnetdl struct {
	Scheduling   string          `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	MirrorSite   string          `json:"MirrorSite,omitempty" yaml:"MirrorSite,omitempty" validate:"omitempty,http_url"`
	ResourceType types.VideoType `json:"ResourceType" yaml:"ResourceType" validate:"required,oneof=movie tv"`
}

//nolint:tagliatelle
type downloader struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	Aria2Label string `json:"Aria2Label" yaml:"Aria2Label" validate:"required"`
}

//nolint:tagliatelle
type global struct {
	LogLevel string `json:"LogLevel" yaml:"LogLevel" validate:"required,oneof=debug info warn error panic fatal"`
	Report   bool   `json:"Report" yaml:"Report" validate:"required"`
	Proxy    struct {
		URL string `json:"Url" yaml:"Url" validate:"required,url"`
	} `json:"Proxy" yaml:"Proxy" validate:"omitempty"`
}

//nolint:tagliatelle
type aria2 struct {
	URL   string `json:"Url" yaml:"Url" validate:"required,http_url"`
	Token string `json:"Token" yaml:"Token" validate:"required"`
	Label string `json:"Label" yaml:"Label" validate:"required"`
}

//nolint:tagliatelle
type tg struct {
	BotToken string `json:"BotToken" yaml:"BotToken" validate:"required"`
	TgIDs    []int  `json:"TgIDs" yaml:"TgIDs" validate:"required"`
}

//nolint:tagliatelle
type mysql struct {
	Host     string `json:"Host" yaml:"Host" validate:"required,ip"`
	Port     int    `json:"Port" yaml:"Port" validate:"gte=0,lte=65535"`
	Database string `json:"Database" yaml:"Database" validate:"required"`
	User     string `json:"User" yaml:"User" validate:"required"`
	Password string `json:"Password" yaml:"Password" validate:"required"`
}

//nolint:tagliatelle
type DouBan struct {
	DouBanList []*DouBan `json:"DouBanList,omitempty" yaml:"DouBanList,omitempty" validate:"required,omitempty"`
	Scheduling string    `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	URL        string    `json:"Url,omitempty" yaml:"Url,omitempty" validate:"omitempty"`
}

/*
//nolint:tagliatelle
type tmDB struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	APIKey     string `json:"APIKey" yaml:"APIKey" validate:"required"`
}

*/

//nolint:tagliatelle
type config struct {
	// TmDB   *tmDB   `json:"TmDB"`
	MySQL        *mysql   `json:"MySQL" yaml:"MySQL" validate:"required"`
	DouBan       *DouBan  `json:"DouBan" yaml:"DouBan" validate:"required"`
	ExcludeWords []string `json:"ExcludeWords" yaml:"ExcludeWords" validate:"required"`
	Feed         struct {
		BTBT           *btbt           `json:"BTBT" yaml:"BTBT" validate:"required"`
		EZTV           *eztv           `json:"EZTV" yaml:"EZTV" validate:"required"`
		GLODLS         *glodls         `json:"GLODLS" yaml:"GLODLS" validate:"required"`
		TGX            *tgx            `json:"TGX" yaml:"TGX" validate:"required"`
		TORLOCK        []*torlock      `json:"TORLOCK" yaml:"TORLOCK" validate:"required"`
		MagnetDL       []*magnetdl     `json:"MagnetDL" yaml:"MagnetDL" validate:"required"`
		TPBPIRATEPROXY *tpbpirateproxy `json:"TPBPIRATEPROXY" yaml:"TPBPIRATEPROXY" validate:"required"`
	} `json:"Feed" yaml:"Feed" validate:"required"`
	Global     *global     `json:"Global" yaml:"Global" validate:"required"`
	Downloader *downloader `json:"Downloader" yaml:"Downloader" validate:"required"`
	Aria2cList []aria2     `json:"Aria2cList" yaml:"Aria2cList" validate:"required"`
	TG         *tg         `json:"TG" yaml:"TG" validate:"omitempty"`
}

//nolint:gochecknoglobals
var (
	Config config
)

func InitConfig(config string) {
	v := viper.New()
	v.SetConfigType("yaml")

	fmt.Printf("config file is %s.\n", config)
	b, err := os.ReadFile(config)
	if err != nil {
		fmt.Printf("配置文件读取错误,err:%s\n", err.Error())
		os.Exit(1)
	}

	err = v.ReadConfig(bytes.NewReader(b))
	if err != nil {
		fmt.Printf("配置文件错误.")
		os.Exit(1)
	}

	err = v.Unmarshal(&Config)
	if err != nil {
		fmt.Println("读取配置错误")
		os.Exit(1)
	}

	//  设置豆瓣列表的调度时间
	for _, ban := range Config.DouBan.DouBanList {
		if ban.Scheduling == "" {
			ban.Scheduling = Config.DouBan.Scheduling
		}
	}

	// 打印 日志级别
	log.Init(true)
	log.SetLogLevel(Config.Global.LogLevel)
	log.Debug("日志级别： ", Config.Global.LogLevel)

	ValidateFc(Config)
}

func ValidateFc(s interface{}) {
	validate := validator.New()
	//nolint:errorlint,forcetypeassert,errorlint
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("配置项: %s 条件: %s %v 当前值: %#v\n", err.StructField(), err.Tag(), err.Param(), err.Value())
		}
		os.Exit(1)
	}
}
