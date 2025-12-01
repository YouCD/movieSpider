package config

import (
	"context"
	"fmt"
	"movieSpider/internal/types"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/natefinch/lumberjack"
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
	LogLevel        string `json:"LogLevel" yaml:"LogLevel" validate:"required,oneof=debug info warn error panic fatal"`
	LogFile         string `json:"LogFile" yaml:"LogFile" validate:"omitempty"`
	Report          bool   `json:"Report" yaml:"Report" validate:"required"`
	IPProxyPool     string `json:"IPProxyPool" yaml:"IPProxyPool" validate:"omitempty,http_url"`
	DHTThread       int    `json:"DHTThread" yaml:"DHTThread"`
	NameParserModel string `json:"NameParserModel" yaml:"NameParserModel"  validate:"http_url"`
	Timeout         int    `json:"Timeout" yaml:"Timeout" validate:"required"`
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
	Host     string `json:"Host" yaml:"Host" validate:"required"`
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
		EZTV          *types.BaseFeed `json:"EZTV" yaml:"EZTV" validate:"required"`
		GLODLS        *types.BaseFeed `json:"GLODLS" yaml:"GLODLS" validate:"required"`
		TORLOCK       []*BaseRT       `json:"TORLOCK" yaml:"TORLOCK" validate:"required"`
		Web1337x      []*BaseRT       `json:"Web1337x" yaml:"Web1337x" validate:"required"`
		ThePirateBay  *types.BaseFeed `json:"ThePirateBay" yaml:"ThePirateBay" validate:"required"`
		Knaben        *types.BaseFeed `json:"Knaben" yaml:"Knaben" validate:"required"`
		TheRarbg      []*BaseRT       `json:"TheRarbg" yaml:"TheRarbg" validate:"required"`
		Uindex        []*BaseRT       `json:"Uindex" yaml:"Uindex" validate:"required"`
		Ilcorsaronero []*BaseRT       `json:"Ilcorsaronero" yaml:"Ilcorsaronero" validate:"required"`
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
		log.WithCtx(context.Background()).Infof("Config file changed: %s\n", e.Name)
		c := new(config)
		// 解析配置文件，反序列化
		err = v.Unmarshal(c)
		if err != nil {
			log.WithCtx(context.Background()).Errorf("Unmarshal yaml faild: %s", err)
			os.Exit(-1)
		}
		err = ValidateFc(Config)
		if err == nil {
			Config = c
			logConfig := &log.Config{
				Stdout: true,
			}
			if Config.Global.LogFile != "" {
				logConfig.LumberjackCfg = &lumberjack.Logger{
					Filename: Config.Global.LogFile,
				}
			}
			log.SetLogLevel(Config.Global.LogLevel)
			log.WithCtx(context.Background()).Debug("日志级别： ", Config.Global.LogLevel)
		}
	})

	logConfig := &log.Config{
		Stdout: true,
	}
	if Config.Global.LogFile != "" {
		logConfig.LumberjackCfg = &lumberjack.Logger{
			Filename: Config.Global.LogFile,
		}
	}
	// 打印 日志级别
	log.Init(logConfig)
	log.SetLogLevel(Config.Global.LogLevel)
	log.WithCtx(context.Background()).Debug("日志级别： ", Config.Global.LogLevel)
	err = ValidateFc(Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ValidateFc(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		//nolint:errorlint
		if _, ok := err.(*validator.InvalidValidationError); ok {
			//nolint:wrapcheck
			return err
		}
		//nolint:errorlint,forcetypeassert
		for _, err := range err.(validator.ValidationErrors) {
			//nolint:err113
			return fmt.Errorf("配置项: %s 条件: %s %v 当前值: %#v\n", err.StructField(), err.Tag(), err.Param(), err.Value())
		}
	}
	return nil
}
