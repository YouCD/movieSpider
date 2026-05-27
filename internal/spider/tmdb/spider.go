package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"movieSpider/internal/model"
	"movieSpider/internal/types"

	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
)

// TMDBSpider TMDB爬虫
type TMDBSpider struct {
	client *Client
}

// NewTMDBSpider 创建TMDB爬虫
func NewTMDBSpider(accountID int, bearerToken string) (*TMDBSpider, error) {
	client, err := NewClient(accountID, bearerToken)
	if err != nil {
		return nil, fmt.Errorf("创建TMDB客户端失败: %w", err)
	}

	return &TMDBSpider{
		client: client,
	}, nil
}

// Run 运行TMDB爬虫，每10分钟执行一次
func (s *TMDBSpider) Run(ctx context.Context) {
	// 立即执行一次
	s.crawl(ctx)

	// 创建定时任务，每10分钟执行一次
	c := cron.New()
	_, err := c.AddFunc("*/10 * * * *", func() {
		s.crawl(ctx)
	})
	if err != nil {
		log.WithCtx(ctx).Errorf("添加TMDB定时任务失败: %s", err)
		return
	}

	// 添加检查提供商的定时任务，每10分钟执行一次
	_, err = c.AddFunc("*/10 * * * *", func() {
		s.checkWatchProviders(ctx)
	})
	if err != nil {
		log.WithCtx(ctx).Errorf("添加检查提供商定时任务失败: %s", err)
		return
	}

	c.Start()

	// 阻塞等待context取消
	<-ctx.Done()
	c.Stop()
	log.WithCtx(ctx).Info("TMDB爬虫已停止")
}

// crawl 执行爬取任务
func (s *TMDBSpider) crawl(ctx context.Context) {
	log.WithCtx(ctx).Info("开始爬取TMDB Watchlist")

	// 爬取电影Watchlist
	s.crawlMovies(ctx)

	// 爬取电视剧Watchlist
	s.crawlTV(ctx)

	log.WithCtx(ctx).Info("TMDB Watchlist爬取完成")
}

// crawlMovies 爬取电影Watchlist
func (s *TMDBSpider) crawlMovies(ctx context.Context) {
	page := 1
	response, err := s.client.GetWatchlistMovies(ctx, page)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取电影Watchlist失败: %s", err)
		return
	}
	for _, movie := range response.Results {
		s.processMovie(ctx, movie)
	}
}

// crawlTV 爬取电视剧Watchlist
func (s *TMDBSpider) crawlTV(ctx context.Context) {
	page := 1
	response, err := s.client.GetWatchlistTV(ctx, page)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取电视剧Watchlist失败: %s", err)
		return
	}

	for _, tv := range response.Results {
		s.processTV(ctx, tv)
	}
}

// processMovie 处理电影数据
func (s *TMDBSpider) processMovie(ctx context.Context, movie types.WatchlistMovieItem) {
	// 获取电影详情以获取IMDB ID
	details, err := s.client.GetMovieDetails(ctx, movie.ID)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取电影详情失败, movieID: %d, err: %s", movie.ID, err)
		return
	}

	// 构建名称列表
	names := []string{movie.Title}
	if movie.OriginalTitle != "" && movie.OriginalTitle != movie.Title {
		names = append(names, movie.OriginalTitle)
	}
	namesJSON, err := json.Marshal(names)
	if err != nil {
		log.WithCtx(ctx).Errorf("序列化名称失败: %s", err)
		return
	}

	// 构建原始数据
	rowData, err := json.Marshal(details)
	if err != nil {
		log.WithCtx(ctx).Errorf("序列化原始数据失败: %s", err)
		return
	}

	// 创建TMDBVideo
	video := &types.TMDBVideo{
		Names:         string(namesJSON),
		ImdbID:        details.ImdbID,
		RowData:       string(rowData),
		Timestamp:     time.Now().Unix(),
		Type:          "movie",
		Playable:      "false",
		DatePublished: movie.ReleaseDate,
	}

	// 存入数据库
	db := model.NewMovieDB()
	if err := db.CreatTMDBVideo(ctx, video); err != nil {
		if err == model.ErrDataExist {
			log.WithCtx(ctx).Debugf("电影已存在: %s (IMDB: %s)", movie.Title, details.ImdbID)
			return
		}
		log.WithCtx(ctx).Errorf("创建电影记录失败: %s", err)
		return
	}

	log.WithCtx(ctx).Infof("成功添加电影: %s (IMDB: %s)", movie.Title, details.ImdbID)
}

// processTV 处理电视剧数据
func (s *TMDBSpider) processTV(ctx context.Context, tv types.WatchlistTVItem) {
	// 获取电视剧详情
	details, err := s.client.GetTVDetails(ctx, tv.ID)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取电视剧详情失败, tvID: %d, err: %s", tv.ID, err)
		return
	}

	// 获取外部ID以获取IMDB ID
	externalIDs, err := s.client.GetTVExternalIDs(ctx, tv.ID)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取电视剧外部ID失败, tvID: %d, err: %s", tv.ID, err)
		return
	}

	// 构建名称列表
	names := []string{tv.Name}
	if tv.OriginalName != "" && tv.OriginalName != tv.Name {
		names = append(names, tv.OriginalName)
	}
	namesJSON, err := json.Marshal(names)
	if err != nil {
		log.WithCtx(ctx).Errorf("序列化名称失败: %s", err)
		return
	}

	// 构建原始数据
	rowData, err := json.Marshal(details)
	if err != nil {
		log.WithCtx(ctx).Errorf("序列化原始数据失败: %s", err)
		return
	}

	// 解析出电视剧的所有季信息
	var seasons []types.SeasonInfo
	for _, season := range details.Seasons {
		for i := range season.EpisodeCount {
			seasons = append(seasons, types.SeasonInfo{
				S:    season.SeasonNumber,
				E:    i + 1,
				Date: season.AirDate,
			})
		}
	}
	seasonsInfo, _ := json.Marshal(seasons)
	s.client.GetTVSeasonDetails(ctx, tv.ID, 1)

	// 创建TMDBVideo
	video := &types.TMDBVideo{
		Names:         string(namesJSON),
		ImdbID:        externalIDs.ImdbID,
		RowData:       string(rowData),
		Timestamp:     time.Now().Unix(),
		Type:          "tv",
		Playable:      "false",
		DatePublished: tv.FirstAirDate,
		SeasonInfo:    string(seasonsInfo),
	}

	// 存入数据库
	db := model.NewMovieDB()
	if err := db.CreatTMDBVideo(ctx, video); err != nil {
		if err == model.ErrDataExist {
			log.WithCtx(ctx).Debugf("电视剧已存在: %s (IMDB: %s)", tv.Name, externalIDs.ImdbID)
			return
		}
		log.WithCtx(ctx).Errorf("创建电视剧记录失败: %s", err)
		return
	}

	log.WithCtx(ctx).Infof("成功添加电视剧: %s (IMDB: %s)", tv.Name, externalIDs.ImdbID)
}

// checkWatchProviders 检查视频的观看提供商
func (s *TMDBSpider) checkWatchProviders(ctx context.Context) {
	log.WithCtx(ctx).Info("开始检查视频的观看提供商")

	db := model.NewMovieDB()
	videos, err := db.FetchPlayableVideos(ctx, false)
	if err != nil {
		log.WithCtx(ctx).Errorf("获取不可播放视频失败: %s", err)
		return
	}

	log.WithCtx(ctx).Infof("找到 %d 个不可播放的视频", len(videos))

	providers := s.checkVideoWatchProviders(ctx, videos...)

	for _, provider := range providers {
		log.WithCtx(ctx).Infof("更新 %s 的播放状态为可播放", provider.Names)
	}
	log.WithCtx(ctx).Info("检查检查视频的观看提供商完成")
}

// checkVideoWatchProviders 检查单个视频的观看提供商
func (s *TMDBSpider) checkVideoWatchProviders(ctx context.Context, videos ...*types.TMDBVideo) []*types.TMDBVideo {
	var result []*types.TMDBVideo
	for _, video := range videos {
		// 解析 row_data 获取 TMDB ID
		tmdbID, err := s.parseTmdbIDFromRowData(video.RowData, video.Type)
		if err != nil {
			log.WithCtx(ctx).Errorf("解析TMDB ID失败, imdbID: %s, err: %s", video.ImdbID, err)
			continue
		}

		// 根据类型获取观看提供商
		var providers *WatchProvidersResponse
		if video.Type == "movie" {
			providers, err = s.client.GetMovieWatchProviders(ctx, tmdbID)
		} else {
			providers, err = s.client.GetTVWatchProviders(ctx, tmdbID)
		}
		if err != nil {
			log.WithCtx(ctx).Errorf("获取观看提供商失败, tmdbID: %d, type: %s, err: %s", tmdbID, video.Type, err)
			continue
		}
		// 检查是否有提供商
		if HasWatchProviders(providers) {
			log.WithCtx(ctx).Infof("发现可观看资源: %s (IMDB: %s, TMDB ID: %d, Type: %s)", video.Names, video.ImdbID, tmdbID, video.Type)

			// 更新播放状态
			db := model.NewMovieDB()
			if err := db.UpdatePlayableStatus(ctx, video.ImdbID); err != nil {
				log.WithCtx(ctx).Errorf("更新播放状态失败, imdbID: %s, err: %s", video.ImdbID, err)
			}
			result = append(result, video)
		}
	}

	return result
}

// parseTmdbIDFromRowData 从 row_data JSON 中解析 TMDB ID
func (s *TMDBSpider) parseTmdbIDFromRowData(rowData, videoType string) (int, error) {
	if videoType == "movie" {
		var detail types.TmDBMovieDetailData
		if err := json.Unmarshal([]byte(rowData), &detail); err != nil {
			return 0, fmt.Errorf("解析电影详情JSON失败: %w", err)
		}
		return detail.ID, nil
	}

	var detail types.TmDBTVDetailData
	if err := json.Unmarshal([]byte(rowData), &detail); err != nil {
		return 0, fmt.Errorf("解析电视剧详情JSON失败: %w", err)
	}
	return detail.ID, nil
}
