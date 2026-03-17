package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// SearchMovieHandler 搜索电影工具处理函数
func SearchMovieHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取电影名称参数

		movieName := request.GetString("movie_name", "")
		if movieName == "" {
			return mcp.NewToolResultError("请提供电影名称参数 movie_name"), nil
		}

		// 调用搜索服务
		results, err := service.SearchMovie(ctx, movieName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("搜索失败: %v", err)), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("未找到电影: %s", movieName)), nil
		}

		// 格式化输出结果
		var output string
		output = fmt.Sprintf("找到 %d 个搜索结果:\n\n", len(results))
		for i, r := range results {
			output += fmt.Sprintf("%d. ID: %d\n", i+1, r.ID)
			output += fmt.Sprintf("   名称: %s\n", r.Name)
			output += fmt.Sprintf("   种子名: %s\n", r.TorrentName)
			output += fmt.Sprintf("   分辨率: %s\n", r.Resolution)
			output += fmt.Sprintf("   类型: %s\n", r.Type)
			output += fmt.Sprintf("   来源: %s\n\n", r.Web)
		}

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		output += "详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}

// DownloadMovieHandler 下载电影工具处理函数
func DownloadMovieHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var result *DownloadResult
		var err error

		// 检查是否提供了movie_id参数
		movieID := request.GetInt("movie_id", 0)

		if movieID > 0 {
			result, err = service.DownloadMovieByID(ctx, int32(movieID))
		} else if movieName := request.GetString("movie_name", ""); movieName != "" {
			// 通过名称下载
			result, err = service.DownloadMovieByName(ctx, movieName)
		} else {
			return mcp.NewToolResultError("请提供 movie_id 或 movie_name 参数"), nil
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("下载失败: %v", err)), nil
		}

		if !result.Success {
			return mcp.NewToolResultError(result.Message), nil
		}

		// 格式化输出结果
		output := fmt.Sprintf("✅ %s\n\n", result.Message)
		output += fmt.Sprintf("电影名称: %s\n", result.MovieName)
		output += fmt.Sprintf("种子名: %s\n", result.TorrentName)
		output += fmt.Sprintf("下载任务GID: %s\n", result.GID)

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		output += "\n详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}

// DownloadProgressHandler 获取下载进度工具处理函数
func DownloadProgressHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取可选的GID参数
		gid := request.GetString("gid", "")
		if gid == "" {
			return mcp.NewToolResultError("请提供 gid 参数"), nil
		}
		// 获取下载进度
		results, err := service.GetDownloadProgress(ctx, gid)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取下载进度失败: %v", err)), nil
		}

		if len(results) == 0 {
			if gid != "" {
				return mcp.NewToolResultText(fmt.Sprintf("未找到GID为 %s 的下载任务", gid)), nil
			}
			return mcp.NewToolResultText("当前没有正在进行的下载任务"), nil
		}

		// 格式化输出结果
		var output string
		output = fmt.Sprintf("共有 %d 个下载任务:\n\n", len(results))
		for i, r := range results {
			output += fmt.Sprintf("%d. GID: %s\n", i+1, r.GID)
			output += fmt.Sprintf("   文件名: %s\n", r.FileName)
			output += fmt.Sprintf("   状态: %s\n", r.Status)
			output += fmt.Sprintf("   总大小: %s\n", r.TotalSize)
			output += fmt.Sprintf("   进度: %s\n\n", r.Progress)
		}

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		output += "详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}

// AllDownloadProgressHandler 获取所有下载任务进度工具处理函数
func AllDownloadProgressHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取所有下载进度
		results, err := service.GetAllDownloadProgress(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取下载进度失败: %v", err)), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultText("当前没有下载任务"), nil
		}

		// 格式化输出结果
		var output string
		output = fmt.Sprintf("共有 %d 个下载任务:\n\n", len(results))
		for i, r := range results {
			output += fmt.Sprintf("%d. GID: %s\n", i+1, r.GID)
			output += fmt.Sprintf("   文件名: %s\n", r.FileName)
			output += fmt.Sprintf("   状态: %s\n", r.Status)
			output += fmt.Sprintf("   总大小: %s\n", r.TotalSize)
			output += fmt.Sprintf("   进度: %s\n", r.Progress)
			if r.ErrorMsg != "" {
				output += fmt.Sprintf("   错误信息: %s\n", r.ErrorMsg)
			}
			output += "\n"
		}

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		output += "详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}

// RemoveDownloadHandler 删除下载任务工具处理函数
func RemoveDownloadHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取GID参数
		gid := request.GetString("gid", "")
		if gid == "" {
			return mcp.NewToolResultError("请提供 gid 参数"), nil
		}

		// 获取是否强制删除参数
		force := request.GetBool("force", false)

		// 删除任务
		if err := service.RemoveDownloadTask(ctx, gid, force); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("删除下载任务失败: %v", err)), nil
		}

		forceStr := ""
		if force {
			forceStr = "(强制)"
		}
		return mcp.NewToolResultText(fmt.Sprintf("✅ 已成功删除下载任务%s，GID: %s", forceStr, gid)), nil
	}
}

// PlayableTodayMovieTV 获取最近24小时内更新为可播放状态的电影或电视剧
func PlayableTodayMovieTV(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		// 调用搜索服务
		results, err := service.PlayableTodayMovieTV(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("搜索失败: %v", err)), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultText("今日无可播放的电影或电视剧，稍后再试"), nil
		}

		// 格式化输出结果
		var output string
		output = fmt.Sprintf("找到 %d 个搜索结果:\n\n", len(results))
		for _, r := range results {
			output += fmt.Sprintf("   名称: %s\n", r.Name)
			output += fmt.Sprintf("   类型: %s\n", r.Type)
		}

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(results, "", "  ")
		output += "详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}
