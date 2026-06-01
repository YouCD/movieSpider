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

// SearchTVHandler 搜索电视剧工具处理函数
func SearchTVHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取电视剧名称参数
		tvName := request.GetString("tv_name", "")
		if tvName == "" {
			return mcp.NewToolResultError("请提供电视剧名称参数 tv_name"), nil
		}

		// 调用搜索服务
		results, err := service.SearchTV(ctx, tvName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("搜索失败: %v", err)), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("未找到电视剧: %s", tvName)), nil
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

// DownloadVideoHandler 下载视频工具处理函数（支持电影和电视剧）
func DownloadVideoHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var result *DownloadResult
		var err error

		// 检查参数优先级：video_id > video_name > magnet
		videoID := request.GetInt("video_id", 0)
		videoName := request.GetString("video_name", "")
		magnet := request.GetString("magnet", "")

		if videoID > 0 {
			// 通过ID下载
			result, err = service.DownloadVideoByID(ctx, int32(videoID))
		} else if videoName != "" {
			// 通过名称下载
			result, err = service.DownloadVideoByName(ctx, videoName)
		} else if magnet != "" {
			// 通过磁力链接下载
			result, err = service.DownloadByMagnet(ctx, magnet)
		} else {
			return mcp.NewToolResultError("请提供 video_id、video_name 或 magnet 参数"), nil
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("下载失败: %v", err)), nil
		}

		if !result.Success {
			return mcp.NewToolResultError(result.Message), nil
		}

		// 格式化输出结果
		output := fmt.Sprintf("✅ %s\n\n", result.Message)
		output += fmt.Sprintf("视频名称: %s\n", result.MovieName)
		if result.TorrentName != "" {
			output += fmt.Sprintf("种子名: %s\n", result.TorrentName)
		}
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

func CheckMovieIsPlayable(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取GID参数
		movieName := request.GetString("movie_name", "")
		if movieName == "" {
			return mcp.NewToolResultError("请提供 movie_name 参数"), nil
		}
		// 调用搜索服务
		ok, err := service.CheckMovieIsPlayable(ctx, movieName)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("搜索失败: %v", err)), nil
		}

		if ok {
			return mcp.NewToolResultText("该电影今日可播放"), nil
		}

		return mcp.NewToolResultText("该电影今日不可播放,有可能输入的名称不正确"), nil
	}
}

// GetWatchlistHandler 获取TMDB收藏列表工具处理函数
func GetWatchlistHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取类型参数
		listType := request.GetString("type", "all")
		page := request.GetInt("page", 1)
		if page <= 0 {
			page = 1
		}

		// 根据类型获取收藏列表
		allItems, totalPages, totalResults, err := fetchWatchlistByType(ctx, service, listType, page)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if len(allItems) == 0 {
			return mcp.NewToolResultText("收藏列表为空"), nil
		}

		// 格式化输出结果
		output := formatWatchlistOutput(allItems, page, totalPages, totalResults)

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(allItems, "", "  ")
		output += "详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}

// fetchWatchlistByType 根据类型获取收藏列表
func fetchWatchlistByType(ctx context.Context, service *MovieService, listType string, page int) ([]WatchlistItem, int, int, error) {
	switch listType {
	case "movie":
		items, tp, tr, err := service.GetWatchlistMovies(ctx, page)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("获取收藏电影列表失败: %w", err)
		}
		return items, tp, tr, nil
	case "tv":
		items, tp, tr, err := service.GetWatchlistTV(ctx, page)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("获取收藏电视剧列表失败: %w", err)
		}
		return items, tp, tr, nil
	default:
		return fetchAllWatchlist(ctx, service, page)
	}
}

// fetchAllWatchlist 获取所有收藏列表（电影和电视剧）
func fetchAllWatchlist(ctx context.Context, service *MovieService, page int) ([]WatchlistItem, int, int, error) {
	// 获取电影
	movies, mp, mr, err := service.GetWatchlistMovies(ctx, page)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("获取收藏电影列表失败: %w", err)
	}
	// 获取电视剧
	tvs, tp, tr, err := service.GetWatchlistTV(ctx, page)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("获取收藏电视剧列表失败: %w", err)
	}
	allItems := make([]WatchlistItem, 0, len(movies)+len(tvs))
	allItems = append(allItems, movies...)
	allItems = append(allItems, tvs...)
	return allItems, mp + tp, mr + tr, nil
}

// formatWatchlistOutput 格式化收藏列表输出
func formatWatchlistOutput(items []WatchlistItem, page, totalPages, totalResults int) string {
	var output string
	output = fmt.Sprintf("TMDB 收藏列表 (第 %d 页，共 %d 页，%d 个结果):\n\n", page, totalPages, totalResults)
	for i, item := range items {
		output += fmt.Sprintf("%d. [%s] ID: %d\n", i+1, item.Type, item.ID)
		output += fmt.Sprintf("   名称: %s\n", item.Title)
		if item.ReleaseDate != "" {
			output += fmt.Sprintf("   上映日期: %s\n", item.ReleaseDate)
		}
		if item.FirstAirDate != "" {
			output += fmt.Sprintf("   首播日期: %s\n", item.FirstAirDate)
		}
		output += fmt.Sprintf("   评分: %.1f\n", item.VoteAverage)
		if len(item.Overview) > 100 {
			output += fmt.Sprintf("   简介: %s...\n", item.Overview[:100])
		} else if item.Overview != "" {
			output += fmt.Sprintf("   简介: %s\n", item.Overview)
		}
		// 显示季信息
		if len(item.SeasonInfo) > 0 {
			output += "   季信息:\n"
			for _, season := range item.SeasonInfo {
				output += fmt.Sprintf("     - 第%d季: %s (%d集, 评分: %.1f)\n",
					season.SeasonNumber, season.Name, season.EpisodeCount, season.VoteAverage)
				if season.AirDate != "" {
					output += fmt.Sprintf("       首播: %s\n", season.AirDate)
				}
			}
		}
		output += "\n"
	}
	return output
}

// CheckTVEpisodeIsPlayableHandler 检查电视剧某季某集是否可播放
func CheckTVEpisodeIsPlayableHandler(service *MovieService) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tvName := request.GetString("tv_name", "")
		if tvName == "" {
			return mcp.NewToolResultError("请提供 tv_name 参数"), nil
		}
		season := request.GetInt("season_number", 0)
		if season <= 0 {
			return mcp.NewToolResultError("请提供正确的 season_number 参数"), nil
		}
		episode := request.GetInt("episode_number", 0)
		if episode <= 0 {
			return mcp.NewToolResultError("请提供正确的 episode_number 参数"), nil
		}

		// 调用服务
		info, err := service.CheckTVEpisodeIsPlayable(ctx, tvName, season, episode)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("检查失败: %v", err)), nil
		}

		// 格式化输出
		var output string
		output += fmt.Sprintf("电视剧: %s\n", info.TVName)
		output += fmt.Sprintf("集数: S%02dE%02d\n", info.SeasonNumber, info.EpisodeNumber)
		if info.IsPlayable {
			output += "状态: 可播放\n"
			if info.Providers != "" {
				output += fmt.Sprintf("提供者: %s\n", info.Providers)
			}
		} else {
			output += "状态: 不可播放\n"
		}

		// 同时返回JSON格式数据
		jsonData, _ := json.MarshalIndent(info, "", "  ")
		output += "\n详细数据(JSON格式):\n" + string(jsonData)

		return mcp.NewToolResultText(output), nil
	}
}
