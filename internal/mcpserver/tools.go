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
			mcp.Description("要搜索的电影名称，如：'The.Matrix',只支持英文、数字、点"),
		),
	)
	s.AddTool(searchMovieTool, SearchMovieHandler(movieService))

	// 注册下载电影工具
	downloadMovieTool := mcp.NewTool("download_movie",
		mcp.WithDescription("下载电影资源，通过电影ID或名称获取种子链接并传给aria2下载器进行下载"),
		mcp.WithNumber("movie_id",
			mcp.Description("电影资源ID（可选，与movie_name二选一）"),
		),
		mcp.WithString("movie_name",
			mcp.Description("电影名称（可选，与movie_id二选一）"),
		),
	)
	s.AddTool(downloadMovieTool, DownloadMovieHandler(movieService))

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
}
