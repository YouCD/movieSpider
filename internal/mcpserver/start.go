package mcpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/youcd/toolkit/log"
)

func Start(ctx context.Context, host, apiKey string) error {
	// 创建 MCP Server
	s := server.NewMCPServer(
		"Movie Spider MCP Server",
		"1.0.0",
		server.WithLogging(),
		server.WithInstructions("新闻咨询助手，请输入指令进行操作。"),
	)

	// 创建服务实例
	movieService := NewMovieService()

	// 注册所有工具
	RegisterTools(s, movieService)

	// 启动服务器
	log.WithCtx(ctx).Info("Movie Spider MCP Server 正在启动...")
	stream := server.NewStreamableHTTPServer(s, server.WithLogger(log.GetLogger()))
	mux := http.NewServeMux()
	mux.Handle("/mcp", tokenAuth(stream, apiKey))
	httpSrv := &http.Server{
		Addr:    host,
		Handler: mux,
	}
	// 创建错误通道用于接收服务器错误
	errCh := make(chan error, 1)

	go func() {
		log.WithCtx(ctx).Infof("MCP server started at %s", httpSrv.Addr)
		if err := stream.Start(httpSrv.Addr); err != nil {
			errCh <- err
		}
	}()
	// 等待上下文取消或服务器错误
	select {
	case <-ctx.Done():
		log.WithCtx(ctx).Info("正在关闭 MCP server...")
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			return err
		}
	case err := <-errCh:
		return err
	}
	return nil
}
