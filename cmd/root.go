package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"movieSpider/internal/log"
	"movieSpider/internal/movieSpiderCore"
	"os"
)

var (
	Name       = "movieSpider"
	configFile string
	runBotFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: fmt.Sprintf("%s 电影助手，自动获取电影种子信息，自动刮取豆瓣电影想看列表，自动下载", Name),

	Run: func(cmd *cobra.Command, args []string) {
		movieSpider := movieSpiderCore.NewMovieSpider(
			movieSpiderCore.WithConfigFile(configFile),
			movieSpiderCore.WithFeeds(),
			movieSpiderCore.WithDownload(),
			movieSpiderCore.WithReport(),
		)

		switch {
		case runBotFlag == true:
			movieSpider.RunWithTGBot()
		}

		movieSpider.RunWithFeed()
		movieSpider.RunWithSpider()
		select {}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config.file", "f", "", "指定配置文件")
	rootCmd.Flags().BoolVar(&runBotFlag, "run.bot", false, "同时运行Telegram bot")
}
