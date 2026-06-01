package mcpserver

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools 注册所有MCP工具
func RegisterTools(s *server.MCPServer, movieService *MovieService) {
	// 注册搜索电影工具
	searchMovieTool := mcp.NewTool("search_movie",
		mcp.WithDescription("搜索电影资源，返回电影的名称、ID和分辨率信息"),
		mcp.WithString("movie_name",
			mcp.Required(),
			mcp.Description("要搜索的电影名称，如：'The Matrix',请提供英文名称"),
		),
	)
	s.AddTool(searchMovieTool, SearchMovieHandler(movieService))

	// 注册搜索电视剧工具
	searchTVTool := mcp.NewTool("search_tv",
		mcp.WithDescription("搜索电视剧资源，返回电视剧的名称、ID和分辨率信息"),
		mcp.WithString("tv_name",
			mcp.Required(),
			mcp.Description("要搜索的电视剧名称，如：'Breaking Bad',请提供英文名称"),
		),
	)
	s.AddTool(searchTVTool, SearchTVHandler(movieService))

	// 注册下载视频工具（支持电影和电视剧）
	downloadVideoTool := mcp.NewTool("download_video",
		mcp.WithDescription("下载电影或电视剧资源，支持通过ID、名称或磁力链接下载"),
		mcp.WithNumber("video_id",
			mcp.Description("视频资源ID（可选，与video_name或magnet三选一）"),
		),
		mcp.WithString("video_name",
			mcp.Description("视频名称（可选，与video_id或magnet三选一）"),
		),
		mcp.WithString("magnet",
			mcp.Description("磁力链接（可选，与video_id或video_name三选一）"),
		),
	)
	s.AddTool(downloadVideoTool, DownloadVideoHandler(movieService))

	// 注册获取所有下载任务进度工具
	allDownloadProgressTool := mcp.NewTool("all_download_progress",
		mcp.WithDescription("获取aria2当前所有的下载任务进度（包括活动和已停止的任务）"),
	)
	s.AddTool(allDownloadProgressTool, AllDownloadProgressHandler(movieService))

	// 注册删除下载任务工具
	removeDownloadTool := mcp.NewTool("remove_download",
		mcp.WithDescription("删除aria2下载任务"),
		mcp.WithString("gid",
			mcp.Required(),
			mcp.Description("要删除的下载任务GID"),
		),
		mcp.WithBoolean("force",
			mcp.Description("是否强制删除（默认false）"),
		),
	)
	s.AddTool(removeDownloadTool, RemoveDownloadHandler(movieService))

	// 注册最近24小时内更新为可播放状态的视频
	playableTodayMovieTvTool := mcp.NewTool("playable_today_movie_tv",
		mcp.WithDescription("获取关注列表中电影或电视剧的可播放状态，返回可播放的电影或电视剧的名称信息"),
	)
	s.AddTool(playableTodayMovieTvTool, PlayableTodayMovieTV(movieService))

	// 注册最近24小时内更新为可播放状态的视频
	checkMovieIsPlayableTool := mcp.NewTool("check_movie_is_playable",
		mcp.WithDescription("检查电影是否可播放"),
		mcp.WithString("movie_name",
			mcp.Required(),
			mcp.Description("请提供电影的英文名称"),
		),
	)
	s.AddTool(checkMovieIsPlayableTool, CheckMovieIsPlayable(movieService))

	// 注册获取TMDB收藏列表工具
	getWatchlistTool := mcp.NewTool("get_tmdb_watchlist",
		mcp.WithDescription("获取TMDB收藏列表（包括电影和电视剧）"),
		mcp.WithString("type",
			mcp.Description("列表类型：'movie'获取电影，'tv'获取电视剧，'all'或不填获取全部（默认all）"),
		),
		mcp.WithNumber("page",
			mcp.Description("页码，默认为1"),
		),
	)
	s.AddTool(getWatchlistTool, GetWatchlistHandler(movieService))

	// 注册检查电视剧某季某集是否可播放工具
	checkTVEpisodePlayableTool := mcp.NewTool("check_tv_episode_is_playable",
		mcp.WithDescription("检查电视剧某季某集是否可播放"),
		mcp.WithString("tv_name",
			mcp.Required(),
			mcp.Description("电视剧的名称"),
		),
		mcp.WithNumber("season_number",
			mcp.Required(),
			mcp.Description("第几季"),
		),
		mcp.WithNumber("episode_number",
			mcp.Required(),
			mcp.Description("第几集"),
		),
	)
	s.AddTool(checkTVEpisodePlayableTool, CheckTVEpisodeIsPlayableHandler(movieService))
}
