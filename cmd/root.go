package cmd

import (
	"context"
	"movieSpider/internal/config"
	"movieSpider/internal/core"
	"movieSpider/internal/mcpserver"
	"movieSpider/internal/model"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	Run: func(cmd *cobra.Command, _ []string) {
		// 创建可取消的上下文
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		// 设置信号监听
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		movieSpider := core.NewMovieSpider(
			core.WithConfigFile(configFile),
			core.WithFeeds(),
			core.WithDownload(ctx),
			core.WithReport(),
			core.WithReleaseTimeJob(ctx),
			core.WithDHT(),
			core.WithTMDBSpider(config.Config.TMDB.AccountID, config.Config.TMDB.BearerToken),
		)

		movieSpider.Start(ctx)
		mcpCfg := &mcpserver.McpConfig{
			MovieConfig: &mcpserver.MovieServiceConfig{
				AccountID:   config.Config.TMDB.AccountID,
				BearerToken: config.Config.TMDB.BearerToken,
			},
			Host:   config.Config.MCP.HostPort,
			ApiKey: config.Config.MCP.ApiKey,
		}
		// 启动 MCP server
		if config.Config.MCP != nil {
			go func() {
				err := mcpserver.Start(ctx, mcpCfg)
				if err != nil {
					log.WithCtx(ctx).Errorf("MCP server error: %v", err)
				}
			}()
		}

		// 启动保存视频数据的 goroutine
		go model.NewMovieDB().SaveFeedVideoFromChan(ctx)

		// 等待关闭信号
		sig := <-sigCh
		log.WithCtx(ctx).Infof("接收到信号 %v，正在优雅关闭...", sig)

		// 创建关闭超时上下文
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// 取消主上下文，通知所有组件停止
		cancel()

		// 停止 MovieSpider
		movieSpider.Stop()

		// 等待关闭完成或超时
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.WithCtx(ctx).Warn("优雅关闭超时，强制退出")
		} else {
			log.WithCtx(ctx).Info("优雅关闭完成")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.WithCtx(rootCmd.Context()).Error(err)
		os.Exit(1)
	}
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config.file", "f", "", "指定配置文件")
}
