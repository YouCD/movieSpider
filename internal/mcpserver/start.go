package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/youcd/toolkit/log"
)

type McpConfig struct {
	MovieConfig *MovieServiceConfig
	Host        string
	ApiKey      string
}

func Start(ctx context.Context, config *McpConfig) error {
	// MCP 服务对外暴露, 必须配置非空 API key
	if config.ApiKey == "" {
		return errors.New("MCP ApiKey 不能为空")
	}

	// 创建 MCP Server
	s := server.NewMCPServer(
		"Movie Spider MCP Server",
		"1.0.0",
		server.WithLogging(),
		server.WithInstructions("新闻咨询助手，请输入指令进行操作。"),
	)

	// 创建服务实例
	movieService, err := NewMovieService(config.MovieConfig)
	if err != nil {
		return err
	}

	// 注册所有工具
	RegisterTools(s, movieService)

	log.WithCtx(ctx).Info("Movie Spider MCP Server 正在启动...")
	slogLogger := slog.New(&log.ZapSlogAdapter{Logger: log.GetLogger()})
	stream := server.NewStreamableHTTPServer(s, server.WithStreamableHTTPLogger(slogLogger))

	// StreamableHTTPServer 本身实现了 http.Handler, 挂载到自定义 server 上,
	// 这样可以正常执行 Shutdown 优雅关闭
	mux := http.NewServeMux()
	mux.Handle("/mcp", tokenAuth(stream, config.ApiKey))
	httpSrv := &http.Server{
		Addr:    config.Host,
		Handler: mux,
	}

	// 创建错误通道用于接收服务器错误
	errCh := make(chan error, 1)

	go func() {
		log.WithCtx(ctx).Infof("MCP server started at %s", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// 等待上下文取消或服务器错误
	select {
	case <-ctx.Done():
		log.WithCtx(ctx).Info("正在关闭 MCP server...")
		// 使用独立的超时上下文 (ctx 已取消, 不能直接派生)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("MCP server 关闭失败: %w", err)
		}
	case err := <-errCh:
		return err
	}
	return nil
}
