package model

import (
	"context"
	"errors"
	"fmt"
	log1 "log"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/nameparser"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"os"
	"sync"
	"time"

	// 引入 MySQL 驱动以初始化数据库连接
	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"github.com/youcd/toolkit/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MovieDB struct {
	db    *gorm.DB
	cache *cache.Cache
}

//nolint:gochecknoglobals
var (
	once         sync.Once
	db           = new(gorm.DB)
	ErrDataExist = errors.New("数据已存在")
)

func NewMovieDB() *MovieDB {
	var err error
	once.Do(func() {
		if err = initDatabase(); err != nil {
			log.WithCtx(context.Background()).Error(err)
			os.Exit(1)
		}
	})
	return &MovieDB{
		db,
		cache.New(24*time.Hour, 24*time.Hour),
	}
}

// initDatabase 初始化数据库连接
func initDatabase() error {
	cfg := config.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 创建数据库
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s  CHARACTER SET utf8mb4 ", cfg.Database)
	if err = gormDB.Exec(sql).Error; err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}

	newLogger := createLogger()

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	if err = gormDB.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&types.FeedVideo{}, &types.DownloadHistory{}, &types.TMDBVideo{}); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	// 赋值给全局变量 db
	db = gormDB

	return nil
}

// createLogger 创建 GORM 日志记录器
func createLogger() logger.Interface {
	logLevel := logger.Silent
	if config.Config.Global.LogLevel == "debug" {
		logLevel = logger.Info
	}

	return logger.New(
		log1.New(os.Stdout, "\r\n", log1.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)
}

// SaveFeedVideoFromChan 从通道中获取 feedVideo 并保存
func (m *MovieDB) SaveFeedVideoFromChan(ctx context.Context) {
	// 启动缓存清理协程
	go m.runCacheCleanup(ctx)

	// 启动视频处理协程
	go m.processFeedVideoChannel(ctx)
}

// runCacheCleanup 定期清理缓存
func (m *MovieDB) runCacheCleanup(ctx context.Context) {
	clearTicker := time.NewTicker(time.Hour) // 每小时检查一次
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	defer clearTicker.Stop()

	for {
		select {
		case <-ticker.C:
			log.WithCtx(ctx).Infof("cache size: %d", m.cache.ItemCount())
		case <-clearTicker.C:
			m.cache.DeleteExpired()
			log.WithCtx(ctx).Infof("cache size after cleanup: %d", m.cache.ItemCount())
		case <-ctx.Done():
			return
		}
	}
}

// processFeedVideoChannel 处理 feedVideo 通道
func (m *MovieDB) processFeedVideoChannel(ctx context.Context) {
	const bufferSize = 30
	buffer := make([]*types.FeedVideo, 0, bufferSize)

	for {
		select {
		case item := <-bus.FeedVideoChan:
			if item == nil {
				log.WithCtx(ctx).Debug("Received nil item from feedVideoCh, skipping")
				continue
			}

			feedVideo, err := FilterVideo(item)
			if err != nil {
				log.WithCtx(ctx).Debugw("FilterVideo", "web", item.Web, "TorrentName", item.TorrentName, "err", err)
				continue
			}
			if feedVideo == nil {
				continue
			}

			// 检查缓存中是否存在该 torrent name
			if _, found := m.cache.Get(item.TorrentName); found {
				log.WithCtx(ctx).Infow("Skipping", "web", item.Web, "TorrentName", item.TorrentName)
				continue
			}

			// 添加到缓冲区
			buffer = append(buffer, feedVideo)
			log.WithCtx(ctx).Infow("Received", "web", item.Web, "TorrentName", item.TorrentName)

			// 当缓冲区达到指定大小时进行处理
			if len(buffer) >= bufferSize {
				m.processFeedVideos(ctx, buffer...)
				buffer = make([]*types.FeedVideo, 0, bufferSize)
			}

		case <-ctx.Done():
			// 处理剩余的缓冲区数据
			if len(buffer) > 0 {
				m.processFeedVideos(ctx, buffer...)
			}
			return
		}
	}
}

func (m *MovieDB) processFeedVideos(ctx context.Context, items ...*types.FeedVideo) {
	for _, item := range items {
		// 先将所有项目加入缓存，防止重复处理
		m.cache.Set(item.TorrentName, true, 24*time.Hour)
	}

	// 对于成功过滤的项目，调用nameparser.ModelHandler进行处理
	torrentNamesToParse := make([]string, 0, len(items))
	feedVideoMap := make(map[string]*types.FeedVideo, len(items))

	for _, feedVideo := range items {
		torrentNamesToParse = append(torrentNamesToParse, feedVideo.TorrentName)
		feedVideoMap[feedVideo.TorrentName] = feedVideo
	}

	// 使用模型解析种子名
	results, err := nameparser.ModelHandler(ctx, torrentNamesToParse...)
	if err != nil {
		if errors.Is(err, nameparser.ErrNamesIsEmpty) {
			return
		}
		log.WithCtx(ctx).Errorf("ModelHandler failed: %v", err)
		return
	}

	// 处理解析结果并保存到数据库
	for torrentName, result := range results {
		feedVideo := feedVideoMap[torrentName]
		if feedVideo == nil {
			continue
		}

		// 更新feedVideo信息
		feedVideo.Name = result.NewName
		feedVideo.Year = cast.ToString(result.Year)
		feedVideo.Type = result.TypeStr

		// 保存到数据库
		if err = m.CreatFeedVideo(feedVideo); err != nil {
			if errors.Is(err, ErrDataExist) {
				log.WithCtx(ctx).Debugf("%s.%s err: %s", feedVideo.Web, feedVideo.Type, err)
				continue
			}
			log.WithCtx(ctx).Error(err)
			continue
		}

		msg := formatSaveMessage(feedVideo)
		log.WithCtx(ctx).Info(msg)
	}
}

// formatSaveMessage 格式化保存消息
func formatSaveMessage(feedVideo *types.FeedVideo) string {
	if feedVideo.Type == "" {
		return fmt.Sprintf("%s: %s 保存完毕.", feedVideo.Web, feedVideo.Name)
	}
	return fmt.Sprintf("%s.%s: %s 保存完毕.", feedVideo.Web, feedVideo.Type, feedVideo.Name)
}

func (m *MovieDB) GetDB() *gorm.DB {
	return m.db
}

func FilterVideo(feedVideoBase *types.FeedVideoBase) (*types.FeedVideo, error) {
	// 排除低码率的视频
	if ok := tools.ExcludeVideo(feedVideoBase.TorrentName, config.Config.ExcludeWords); ok {
		//nolint:err113
		return nil, fmt.Errorf("excludeWords, web:%s,TorrentName:%v", feedVideoBase.Web, feedVideoBase.TorrentName)
	}

	feedVideo := &types.FeedVideo{
		FeedVideoBase: *feedVideoBase,
	}
	if feedVideo.TorrentName == "" {
		//nolint:err113
		return nil, fmt.Errorf("feedVideo.TorrentName is empty: %#v", feedVideo)
	}

	// 解析前先查库
	video, err := NewMovieDB().GetFeedVideoByName(feedVideo.TorrentName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithCtx(context.Background()).Error(err)
	}
	if video != nil {
		return nil, fmt.Errorf("torrent_name:%s, err:%w", feedVideo.TorrentName, ErrFeedVideoExist)
	}

	return feedVideo, nil
}
