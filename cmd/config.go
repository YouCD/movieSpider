package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"movieSpider/internal/log"
	"os"
)

var (
	cfgTmp = `MySQL:
  IP: 127.0.0.1
  Port: 3306
  Database: movie
  User: root
  Password: P@ssw0rd

#  排除包含以下关键字的资源
ExcludeWords:
  - 720p
  - dvsux
  - 480p
  - hdr
  - .dv.
  - .dolby.vision
DouBan:
  # 豆瓣电影想看清单
  Scheduling: "*/10 * * * *"
  DouBanList:
    - Url: "https://movie.douban.com/people/251312920/wish"
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
  TORLOCK:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  MAGNETDL:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  TPBPIRATEPROXY:
    Scheduling: "*/3 * * * *"
Global:
  LogLevel: debug
  Report: true

# Downloader 下载
Downloader:
  Scheduling: "*/60 * * * *"
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
)
var configCmd = &cobra.Command{
	Use:   "config",
	Short: fmt.Sprintf("generate %s config file.", Name),
	Run: func(cmd *cobra.Command, args []string) {
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

	},
}

func init() {
	configCmd.Flags().StringVarP(&outFile, "out.file", "o", "", "指定输出的文件")
}
