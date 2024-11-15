package config

import (
	"fmt"
	"movieSpider/internal/types"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/youcd/toolkit/log"
)

//nolint:tagliatelle
type downloader struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	Aria2Label string `json:"Aria2Label" yaml:"Aria2Label" validate:"required"`
}

//nolint:tagliatelle
type global struct {
	LogLevel    string `json:"LogLevel" yaml:"LogLevel" validate:"required,oneof=debug info warn error panic fatal"`
	Report      bool   `json:"Report" yaml:"Report" validate:"required"`
	IPProxyPool string `json:"IPProxyPool" yaml:"IPProxyPool" validate:"omitempty,http_url"`
	DHTThread   int    `json:"DHTThread" yaml:"DHTThread"`
}

//nolint:tagliatelle
type aria2 struct {
	URL   string `json:"URL" yaml:"URL" validate:"required,http_url"`
	Token string `json:"Token" yaml:"Token" validate:"required"`
	Label string `json:"Label" yaml:"Label" validate:"required"`
}

//nolint:tagliatelle
type tg struct {
	BotToken string `json:"BotToken" yaml:"BotToken" validate:"required"`
	TgIDs    []int  `json:"TgIDs" yaml:"TgIDs" validate:"required"`
	ProxyURL string `json:"ProxyUrl" yaml:"ProxyUrl" validate:"omitempty,url"`
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
	URL        string    `json:"URL,omitempty" yaml:"URL,omitempty" validate:"omitempty"`
}

/*
//nolint:tagliatelle
type tmDB struct {
	Scheduling string `json:"Scheduling" yaml:"Scheduling" validate:"cron"`
	APIKey     string `json:"APIKey" yaml:"APIKey" validate:"required"`
}

*/

//nolint:tagliatelle
type tgx struct {
	types.BaseFeed `mapstructure:",squash"`
	Name           string `json:"Name" yaml:"Name" validate:"required"`
}

//nolint:tagliatelle
type BaseRT struct {
	types.BaseFeed `mapstructure:",squash"`
	ResourceType   types.VideoType `json:"ResourceType" yaml:"ResourceType" validate:"required,oneof=movie tv"`
}

//nolint:tagliatelle
type config struct {
	// TmDB   *tmDB   `json:"TmDB"`
	MySQL        *mysql   `json:"MySQL" yaml:"MySQL" validate:"required"`
	DouBan       *DouBan  `json:"DouBan" yaml:"DouBan" validate:"required"`
	ExcludeWords []string `json:"ExcludeWords" yaml:"ExcludeWords" validate:"required"`
	Feed         struct {
		BTBT         *types.BaseFeed `json:"BTBT" yaml:"BTBT" validate:"required"`
		EZTV         *types.BaseFeed `json:"EZTV" yaml:"EZTV" validate:"required"`
		GLODLS       *types.BaseFeed `json:"GLODLS" yaml:"GLODLS" validate:"required"`
		TGX          []*tgx          `json:"TGX" yaml:"TGX" validate:"required"`
		TORLOCK      []*BaseRT       `json:"TORLOCK" yaml:"TORLOCK" validate:"required"`
		Web1337x     []*BaseRT       `json:"Web1337x" yaml:"Web1337x" validate:"required"`
		ThePirateBay *types.BaseFeed `json:"ThePirateBay" yaml:"ThePirateBay" validate:"required"`
		Knaben       *types.BaseFeed `json:"Knaben" yaml:"Knaben" validate:"required"`
		Rarbg2       []*BaseRT       `json:"Rarbg2" yaml:"Rarbg2" validate:"required"`
		TheRarbg     []*BaseRT       `json:"TheRarbg" yaml:"TheRarbg" validate:"required"`
	} `json:"Feed" yaml:"Feed" validate:"required"`
	Global     *global     `json:"Global" yaml:"Global" validate:"required"`
	Downloader *downloader `json:"Downloader" yaml:"Downloader" validate:"required"`
	Aria2cList []aria2     `json:"Aria2cList" yaml:"Aria2cList" validate:"required"`
	TG         *tg         `json:"TG" yaml:"TG" validate:"omitempty"`
}

//nolint:gochecknoglobals
var (
	Config *config
	v      = viper.New()
)

func InitConfig(configFile string) {
	v.SetConfigType("yaml")
	v.SetConfigFile(configFile)

	fmt.Printf("config file is %s.\n", configFile)
	err := v.ReadInConfig()
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

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s\n", e.Name)
		c := new(config)
		// 解析配置文件，反序列化
		if err := v.Unmarshal(c); err != nil {
			log.Errorf("Unmarshal yaml faild: %s", err)
			os.Exit(-1)
		}
		if err = ValidateFc(Config); err == nil {
			Config = c
			log.SetLogLevel(Config.Global.LogLevel)
		}

	})
	// 打印 日志级别
	log.Init(true)
	log.SetLogLevel(Config.Global.LogLevel)
	log.Debug("日志级别： ", Config.Global.LogLevel)

	if err = ValidateFc(Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ValidateFc(s interface{}) error {
	validate := validator.New()
	//nolint:errorlint,forcetypeassert,errorlint
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("配置项: %s 条件: %s %v 当前值: %#v\n", err.StructField(), err.Tag(), err.Param(), err.Value())
		}
	}
	return nil
}
