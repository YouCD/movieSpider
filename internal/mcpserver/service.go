package mcpserver

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/tmdb"
	"movieSpider/internal/types"

	"github.com/youcd/toolkit/log"
)

// MovieService 提供电影搜索和下载服务
type MovieService struct {
	aria2Client *aria2.Aria2
	once        sync.Once
	apiKey      string
	tmdbClient  *tmdb.Client
}

type MovieServiceConfig struct {
	AccountID   int
	BearerToken string
}

// NewMovieService 创建新的MovieService实例
func NewMovieService(cfg *MovieServiceConfig) (*MovieService, error) {
	client, err := tmdb.NewClient(cfg.AccountID, cfg.BearerToken)
	if err != nil {
		return nil, fmt.Errorf("初始化 TMDB 客户端失败: %w", err)
	}

	return &MovieService{tmdbClient: client}, nil
}

// getAria2Client 获取aria2客户端（懒加载）
func (s *MovieService) getAria2Client() (*aria2.Aria2, error) {
	var err error
	s.once.Do(func() {
		s.aria2Client, err = aria2.NewAria2(config.Config.Downloader.Aria2Label)
	})
	if err != nil {
		return nil, fmt.Errorf("aria2客户端初始化失败: %w", err)
	}
	return s.aria2Client, nil
}

// MovieResult 电影搜索结果
type MovieResult struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	TorrentName string `json:"torrent_name,omitempty"`
	Resolution  string `json:"resolution,omitempty"`
	Type        string `json:"type"`
	Web         string `json:"web,omitempty"`
	Magnet      string `json:"magnet,omitempty"`
}

// SearchMovie 搜索电影
func (s *MovieService) SearchMovie(ctx context.Context, movieName string) ([]MovieResult, error) {
	log.WithCtx(ctx).Infof("开始搜索电影: %s", movieName)

	// 从数据库获取搜索结果
	videos, err := model.NewMovieDB().GetFeedVideoMovieByNames(ctx, movieName)
	if err != nil {
		if err == model.ErrMoviesNotFound {
			return nil, fmt.Errorf("未找到电影: %s", movieName)
		}
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}
	if len(videos) == 0 {
		//todo 在线搜索
		//// 使用BT4g搜索
		//feedBt4g := searchspider.NewFeedBt4g(movieName, types.ResolutionOther)
		//_, err := feedBt4g.Search()
		//if err != nil {
		//	log.WithCtx(ctx).Errorf("BT4g搜索失败: %v", err)
		//}
	}
	// 转换结果
	var results []MovieResult
	for _, v := range videos {
		resolution := parseResolution(v.TorrentName)
		results = append(results, MovieResult{
			ID:          v.ID,
			Name:        v.Name,
			TorrentName: v.TorrentName,
			Resolution:  resolution,
			Type:        v.Type,
			Web:         v.Web,
			Magnet:      v.Magnet,
		})
	}

	log.WithCtx(ctx).Infof("搜索完成，找到 %d 个结果", len(results))
	return results, nil
}

// SearchTV 搜索电视剧
func (s *MovieService) SearchTV(ctx context.Context, tvName string) ([]MovieResult, error) {
	log.WithCtx(ctx).Infof("开始搜索电视剧: %s", tvName)

	// 从数据库获取搜索结果
	videos, err := model.NewMovieDB().GetFeedVideoTVByNames(ctx, tvName)
	if err != nil {
		if err == model.ErrTVNotFound {
			return nil, fmt.Errorf("未找到电视剧: %s", tvName)
		}
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}
	if len(videos) == 0 {
		return nil, fmt.Errorf("未找到电视剧: %s", tvName)
	}

	// 转换结果
	var results []MovieResult
	for _, v := range videos {
		resolution := parseResolution(v.TorrentName)
		results = append(results, MovieResult{
			ID:          v.ID,
			Name:        v.Name,
			TorrentName: v.TorrentName,
			Resolution:  resolution,
			Type:        v.Type,
			Web:         v.Web,
			Magnet:      v.Magnet,
		})
	}

	log.WithCtx(ctx).Infof("搜索完成，找到 %d 个结果", len(results))
	return results, nil
}

// DownloadResult 下载结果
type DownloadResult struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	GID         string `json:"gid,omitempty"`
	MovieName   string `json:"movie_name,omitempty"`
	TorrentName string `json:"torrent_name,omitempty"`
}

// DownloadMovieByID 通过电影ID下载
func (s *MovieService) DownloadMovieByID(ctx context.Context, id int32) (*DownloadResult, error) {
	log.WithCtx(ctx).Infof("开始下载电影，ID: %d", id)

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 从数据库获取视频信息
	video, err := s.getVideoByID(id)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("获取视频信息失败: %v", err),
		}, nil
	}

	if video.Magnet == "" {
		return &DownloadResult{
			Success: false,
			Message: "该视频没有磁力链接",
		}, nil
	}

	// 开始下载
	gid, err := client.DownloadByWithVideo(ctx, video, video.Magnet)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("添加下载任务失败: %v", err),
		}, nil
	}

	// 更新下载状态
	if err = model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 1); err != nil {
		log.WithCtx(ctx).Errorf("更新下载状态失败: %v", err)
	}

	return &DownloadResult{
		Success:     true,
		Message:     "下载任务已添加",
		GID:         gid,
		MovieName:   video.Name,
		TorrentName: video.TorrentName,
	}, nil
}

// DownloadVideoByID 通过视频ID下载（支持电影和电视剧）
func (s *MovieService) DownloadVideoByID(ctx context.Context, id int32) (*DownloadResult, error) {
	log.WithCtx(ctx).Infof("开始下载视频，ID: %d", id)

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 从数据库获取视频信息
	video, err := s.getVideoByID(id)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("获取视频信息失败: %v", err),
		}, nil
	}

	if video.Magnet == "" {
		return &DownloadResult{
			Success: false,
			Message: "该视频没有磁力链接",
		}, nil
	}

	// 开始下载
	gid, err := client.DownloadByWithVideo(ctx, video, video.Magnet)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("添加下载任务失败: %v", err),
		}, nil
	}

	// 更新下载状态
	if err = model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 1); err != nil {
		log.WithCtx(ctx).Errorf("更新下载状态失败: %v", err)
	}

	return &DownloadResult{
		Success:     true,
		Message:     "下载任务已添加",
		GID:         gid,
		MovieName:   video.Name,
		TorrentName: video.TorrentName,
	}, nil
}

// DownloadVideoByName 通过视频名称下载（同时搜索电影和电视剧）
func (s *MovieService) DownloadVideoByName(ctx context.Context, name string) (*DownloadResult, error) {
	log.WithCtx(ctx).Infof("开始下载视频，名称: %s", name)

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 先搜索电影
	videos, err := model.NewMovieDB().GetFeedVideoMovieByNames(ctx, name)
	if err != nil && err != model.ErrMoviesNotFound {
		return nil, fmt.Errorf("查询电影数据库失败: %w", err)
	}

	// 如果没有找到电影，搜索电视剧
	if len(videos) == 0 {
		videos, err = model.NewMovieDB().GetFeedVideoTVByNames(ctx, name)
		if err != nil && err != model.ErrTVNotFound {
			return nil, fmt.Errorf("查询电视剧数据库失败: %w", err)
		}
	}

	if len(videos) == 0 {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("未找到视频: %s", name),
		}, nil
	}

	// 选择第一个视频进行下载
	video := videos[0]

	if video.Magnet == "" {
		return &DownloadResult{
			Success: false,
			Message: "该视频没有磁力链接",
		}, nil
	}

	// 开始下载
	gid, err := client.DownloadByWithVideo(ctx, video, video.Magnet)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("添加下载任务失败: %v", err),
		}, nil
	}

	// 更新下载状态
	if err = model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 1); err != nil {
		log.WithCtx(ctx).Errorf("更新下载状态失败: %v", err)
	}

	return &DownloadResult{
		Success:     true,
		Message:     "下载任务已添加",
		GID:         gid,
		MovieName:   video.Name,
		TorrentName: video.TorrentName,
	}, nil
}

// DownloadByMagnet 通过磁力链接下载
func (s *MovieService) DownloadByMagnet(ctx context.Context, magnet string) (*DownloadResult, error) {
	log.WithCtx(ctx).Infof("开始通过磁力链接下载")

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 创建一个虚拟的视频对象用于下载
	video := &types.FeedVideo{
		Name: "磁力链接下载",
		FeedVideoBase: types.FeedVideoBase{
			Magnet: magnet,
		},
	}

	// 开始下载
	gid, err := client.DownloadByWithVideo(ctx, video, magnet)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("添加下载任务失败: %v", err),
		}, nil
	}

	return &DownloadResult{
		Success:   true,
		Message:   "下载任务已添加",
		GID:       gid,
		MovieName: "磁力链接下载",
	}, nil
}

// DownloadMovieByName 通过电影名称下载
func (s *MovieService) DownloadMovieByName(ctx context.Context, name string) (*DownloadResult, error) {
	log.WithCtx(ctx).Infof("开始下载电影，名称: %s", name)

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 搜索视频
	videos, err := model.NewMovieDB().GetFeedVideoMovieByNames(ctx, name)
	if err != nil {
		if err == model.ErrMoviesNotFound {
			return &DownloadResult{
				Success: false,
				Message: fmt.Sprintf("未找到电影: %s", name),
			}, nil
		}
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}

	if len(videos) == 0 {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("未找到可下载的电影: %s", name),
		}, nil
	}

	// 选择第一个视频进行下载
	video := videos[0]

	if video.Magnet == "" {
		return &DownloadResult{
			Success: false,
			Message: "该视频没有磁力链接",
		}, nil
	}

	// 开始下载
	gid, err := client.DownloadByWithVideo(ctx, video, video.Magnet)
	if err != nil {
		return &DownloadResult{
			Success: false,
			Message: fmt.Sprintf("添加下载任务失败: %v", err),
		}, nil
	}

	// 更新下载状态
	if err = model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 1); err != nil {
		log.WithCtx(ctx).Errorf("更新下载状态失败: %v", err)
	}

	return &DownloadResult{
		Success:     true,
		Message:     "下载任务已添加",
		GID:         gid,
		MovieName:   video.Name,
		TorrentName: video.TorrentName,
	}, nil
}

// DownloadProgress 下载进度信息
type DownloadProgress struct {
	GID           string `json:"gid"`
	FileName      string `json:"file_name"`
	Status        string `json:"status"`
	TotalSize     string `json:"total_size"`
	Completed     string `json:"completed"`
	Progress      string `json:"progress"`
	DownloadSpeed string `json:"download_speed,omitempty"`
	UploadSpeed   string `json:"upload_speed,omitempty"`
	ErrorMsg      string `json:"error_msg,omitempty"`
}

// GetDownloadProgress 获取下载进度
func (s *MovieService) GetDownloadProgress(ctx context.Context, gid string) ([]DownloadProgress, error) {
	log.WithCtx(ctx).Infof("获取下载进度，GID: %s", gid)

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 获取所有活动任务
	files := client.CurrentActiveAndStopFiles()
	if len(files) == 0 {
		return []DownloadProgress{}, nil
	}

	var results []DownloadProgress
	for _, f := range files {
		progress := DownloadProgress{
			GID:       f.GID,
			FileName:  f.FileName,
			Status:    f.Status,
			TotalSize: f.Size,
			Completed: f.Completed,
			Progress:  f.Completed,
		}

		// 如果指定了GID，只返回该任务
		if gid != "" && f.GID == gid {
			return []DownloadProgress{progress}, nil
		}

		results = append(results, progress)
	}

	// 如果指定了GID但没找到
	if gid != "" {
		return nil, fmt.Errorf("未找到GID为 %s 的下载任务", gid)
	}

	return results, nil
}

// GetAllDownloadProgress 获取所有下载任务进度
func (s *MovieService) GetAllDownloadProgress(ctx context.Context) ([]DownloadProgress, error) {
	log.WithCtx(ctx).Info("获取所有下载任务进度")

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return nil, err
	}

	// 获取所有活动和已停止的任务
	files := client.CurrentActiveAndStopFiles()
	if len(files) == 0 {
		return []DownloadProgress{}, nil
	}

	var results []DownloadProgress
	for _, f := range files {
		results = append(results, DownloadProgress{
			GID:       f.GID,
			FileName:  f.FileName,
			Status:    f.Status,
			TotalSize: f.Size,
			Completed: f.Completed,
			Progress:  f.Completed,
			ErrorMsg:  f.ErrorMsg,
		})
	}

	log.WithCtx(ctx).Infof("找到 %d 个下载任务", len(results))
	return results, nil
}

// RemoveDownloadTask 删除下载任务
func (s *MovieService) RemoveDownloadTask(ctx context.Context, gid string, force bool) error {
	log.WithCtx(ctx).Infof("删除下载任务，GID: %s, 强制: %v", gid, force)

	// 获取aria2客户端
	client, err := s.getAria2Client()
	if err != nil {
		return err
	}

	// 删除任务
	if force {
		return client.ForceRemoveTask(gid)
	}
	return client.RemoveTask(gid)
}

// getVideoByID 通过ID获取视频信息
func (s *MovieService) getVideoByID(id int32) (*types.FeedVideo, error) {
	db := model.NewMovieDB()
	var video types.FeedVideo
	result := db.GetDB().Where("id = ?", id).First(&video)
	if result.Error != nil {
		return nil, fmt.Errorf("未找到ID为 %d 的视频: %w", id, result.Error)
	}
	return &video, nil
}

// 获取今日凌晨时间戳
func getTodayStartTimestamp() int64 {
	now := time.Now()
	year, month, day := now.Date()
	todayStart := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	return todayStart.Unix()
}

// PlayableTodayMovieTV 获取最近24小时内更新为可播放状态的电影或电视剧
func (s *MovieService) PlayableTodayMovieTV(ctx context.Context) ([]MovieResult, error) {
	log.WithCtx(ctx).Infof("开始检查")
	// 今日凌晨整点
	start := getTodayStartTimestamp()

	// 从数据库获取搜索结果
	videos, err := model.NewMovieDB().FetchPlayableVideosAndUpdate(ctx, true, start)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}

	// 转换结果
	var results []MovieResult
	for _, v := range videos {
		results = append(results, MovieResult{
			ID:   int32(v.ID),
			Name: v.Names,
			Type: v.Type,
		})
	}

	log.WithCtx(ctx).Infof("搜索完成，找到 %d 个结果", len(results))
	return results, nil
}

// CheckMovieIsPlayable 检查电影是否可播放
func (s *MovieService) CheckMovieIsPlayable(ctx context.Context, movieName string) (bool, error) {
	log.WithCtx(ctx).Infof("开始检查:%s", movieName)
	movies, err := s.tmdbClient.GetSearchMovies(ctx, movieName)
	if err != nil {
		return false, fmt.Errorf("搜索电影失败: %w", err)
	}
	if len(movies) == 0 {
		return false, fmt.Errorf("未找到电影: %s", movieName)
	}
	providers, err := s.tmdbClient.GetMovieWatchProviders(ctx, int(movies[0].ID))
	if err != nil {
		return false, fmt.Errorf("获取电影播放提供者失败: %w", err)
	}
	if tmdb.HasWatchProviders(providers) {
		log.WithCtx(ctx).Infof("电影 %s 可播放", movieName)
		return true, nil
	}
	return false, nil
}

// CheckTVIsPlayable 检查电视剧是否可播放
func (s *MovieService) CheckTVIsPlayable(ctx context.Context, tvName string, season int) (bool, error) {
	log.WithCtx(ctx).Infof("开始检查:%s", tvName)
	movies, err := s.tmdbClient.GetSearchTVShow(ctx, tvName)
	if err != nil {
		return false, fmt.Errorf("搜索电视剧失败: %w", err)
	}
	if len(movies) == 0 {
		return false, fmt.Errorf("未找到电视剧: %s", tvName)
	}
	providers, err := s.tmdbClient.GetTVSeasonWatchProviders(ctx, int(movies[0].ID), season, nil)
	if err != nil {
		return false, fmt.Errorf("获取电视剧播放提供者失败: %w", err)
	}

	if tmdb.HasWatchProviders(providers) {
		log.WithCtx(ctx).Infof("电视剧: %s 第%d季 可播放", tvName, season)
		return true, nil
	}

	return false, nil
}

// resolutionReg 分辨率正则表达式
var resolutionReg = regexp.MustCompile(`(2160p|2160P|1080p|1080P|720p|720P|4K)`)

// parseResolution 解析分辨率
func parseResolution(torrentName string) string {
	matches := resolutionReg.FindStringSubmatch(torrentName)
	if len(matches) > 1 {
		return matches[1]
	}
	return "未知"
}

func tokenAuth(next http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		log.WithCtx(r.Context()).Debug("token:", token)
		if token != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
