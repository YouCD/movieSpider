package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/youcd/toolkit/log"
	"os"
)

//nolint:gochecknoglobals
var (
	cfgTmp = `MySQL:
  # 这个地址是docker里面的地址
  Host: 127.0.0.1
  Port: 3306
  Database: movie
  User: root
  Password: P@ssw0rd

Douban:
  # 豆瓣电影想看清单
  DoubanUrl:
    - Url: "https://movie.douban.com/people/251312920/wish"
    - Url: "https://movie.douban.com/people/271517237/wish"
  Scheduling: "*/10 * * * *"

#  排除包含以下关键字的资源
ExcludeWords:
  - 720p
  - dvsux
  - 480p
  - hdr
  - .dv.
  - .dolby.vision
Feed:
  BTBT:
    Scheduling: "*/5 * * * *"
  EZTV:
    Scheduling: "*/5 * * * *"
    MirrorSite: "https://eztvx.to"
  GLODLS:
    Scheduling: "*/3 * * * *"
    MirrorSite: "https://gtso.cc"
  TGX:
    Scheduling: "*/3 * * * *"
    MirrorSite: "https://tgx.rs"
  TORLOCK:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
      MirrorSite: "https://torlock.123unblock.art"
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
      MirrorSite: "https://torlock.123unblock.art"
  MAGNETDL:
    - Scheduling: "*/3 * * * *"
      ResourceType: movie
    - Scheduling: "*/2 * * * *"
      ResourceType: tv
  TPBPIRATEPROXY:
    Scheduling: "*/3 * * * *"
    MirrorSite: "https://thepiratebay10.info"
Global:
  LogLevel: info
  Report: true
# 网络代理
#  Proxy:
#    Url: socks5://127.0.0.1:1080
# Downloader 下载
Downloader:
  Scheduling: "*/60 * * * *"
  # 使用哪个 Aria2 下载
  Aria2Label: "home"

# Aria2 下载服务器
Aria2cList:
  - Url: "http://127.0.0.1e:6800"
    Token: P@ssw0rd
    Label: home

# 如果没有Telegram 就请忽略
#TG:
  # Telegram 机器人 token
#  BotToken: "TOKEN"
#   能正常访问机器人的Telegram用户
#  TgIDs: [ 123456 ]


`
	outFile string
)

//nolint:gochecknoglobals, forbidigo
var configCmd = &cobra.Command{
	Use:   "config",
	Short: fmt.Sprintf("generate %s config file.", Name),
	Run: func(_ *cobra.Command, _ []string) {
		if outFile == "" {
			fmt.Println(cfgTmp)
		}
		if outFile != "" {
			err := os.WriteFile(outFile, []byte(cfgTmp), 0644)
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
		}

	},
}

//nolint:gochecknoinits
func init() {
	configCmd.Flags().StringVarP(&outFile, "out.file", "o", "", "指定输出的文件")
}
