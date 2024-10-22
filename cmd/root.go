package cmd

import (
	"movieSpider/internal/core"
	"movieSpider/internal/model"
	"os"

	"github.com/spf13/cobra"
	"github.com/youcd/toolkit/log"
)

//nolint:gochecknoglobals
var (
	Name       = "movieSpider"
	configFile string
)

// rootCmd represents the base command when called without any subcommands
//
//nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: Name + "电影助手，自动获取电影种子信息，自动刮取豆瓣电影想看列表，自动下载",

	Run: func(_ *cobra.Command, _ []string) {
		movieSpider := core.NewMovieSpider(
			core.WithConfigFile(configFile),
			core.WithFeeds(),
			core.WithDownload(),
			core.WithReport(),
			core.WithReleaseTimeJob(),
			core.WithDHT(),
		)

		movieSpider.Start()

		model.NewMovieDB().SaveFeedVideoFromChan()
		select {}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config.file", "f", "", "指定配置文件")
}
