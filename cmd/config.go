package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"os"
)

var (
	cfgTmp = `MySQL:
  Host: 127.0.0.1
  Port: 3306
  Database: movie
  User: root
  Password: P@ssw0rd

Douban:
  # 豆瓣电影想看清单
  DoubanUrl: "https://movie.douban.com/people/251312920/wish"
  Scheduling: "*/10 * * * *"
  # 豆瓣电影 公共API
  WMDBPrefix: "https://api.wmdb.tv/movie/api?id="
Feed:
  # 代理池 这里使用 https://github.com/jhao104/proxy_pool
  ProxyPool: "http://127.0.0.1:5010"
  BTBT:
    Scheduling: "*/5 * * * *"
  EZTV:
    Scheduling: "*/5 * * * *"
  GLODLS:
    Scheduling: "*/3 * * * *"
  TGX:
    Scheduling: "*/3 * * * *"
  RARBG:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  TORLOCK:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
Global:
  LogLevel: debug
  Report: true

# Downloader 下载
Downloader:
  Scheduling: "*/1 * * * *"
  # 使用哪个 Aria2 下载
  Aria2Label: "home"

# Aria2 下载服务器
Aria2cList:
  - Url: "http://127.0.0.1e:6800"
    Token: 123456
    Label: home
  - Url: "http://127.0.0.1:6801"
    Token: 123456
    Label: nas

# 如果没有Telegram 就请忽略
#TG:
  # Telegram 网络代理
#  Proxy:
#    Url: socks5://127.0.0.1:1080
#    Enable: false
  # Telegram 机器人 token
#  BotToken: "TOKEN"
#   能正常访问机器人的Telegram用户
#  TgIDs: [ 123456 ]

`
	outFile string
	initDB  bool
)
var configCmd = &cobra.Command{
	Use:   "config",
	Short: fmt.Sprintf("generate %s config file.", Name),
	Run: func(cmd *cobra.Command, args []string) {
		switch initDB {
		case false:
			if outFile == "" {
				fmt.Println(cfgTmp)
			}
			if outFile != "" {
				err := ioutil.WriteFile(outFile, []byte(cfgTmp), 0644)
				if err != nil {
					log.Error(err)
					os.Exit(-1)
				}
			}
		case true:
			config.InitConfig(configFile)
			model.NewMovieDB()
			err := model.NewMovieDB().InitDBTable()
			if err != nil {
				log.Error(err)
				return
			}
			log.Infof("db: %s 数据库初始化完毕.", config.MySQL.Database)
		}

	},
}

func init() {
	configCmd.Flags().StringVarP(&outFile, "out.file", "o", "", "指定输出的文件")
	configCmd.Flags().BoolVar(&initDB, "init.db", false, "初始化DB")
}
