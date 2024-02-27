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
	db          *gorm.DB
	feedVideoCh chan *types.FeedVideoBase
	cache       *cache.Cache
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
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Host, config.Config.MySQL.Port)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			os.Exit(1)
		}
		sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s  CHARACTER SET utf8mb4 ", config.Config.MySQL.Database)
		// 创建数据库
		err = db.Exec(sql).Error
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			os.Exit(1)
		}

		setLogLevel := func() logger.LogLevel {
			if config.Config.Global.LogLevel == "debug" {
				return logger.Info
			}
			return logger.Silent
		}
		newLogger := logger.New(
			log1.New(os.Stdout, "\r\n", log1.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      setLogLevel(), // Log level
				//LogLevel: logger.Info, // Log level
				Colorful: true, // 禁用彩色打印
			},
		)

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Host, config.Config.MySQL.Port, config.Config.MySQL.Database) // 连接数据库

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})

		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			os.Exit(1)
		}

		err = db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&types.FeedVideo{}, &types.DownloadHistory{}, &types.DouBanVideo{})
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			os.Exit(1)
		}
	})
	return &MovieDB{
		db,
		bus.FeedVideoChan,
		cache.New(24*time.Hour, 24*time.Hour),
	}
}

// SaveFeedVideoFromChan
//
//	@Description: 从通道中获取 feedVideo 并保存
//	@receiver m
func (m *MovieDB) SaveFeedVideoFromChan(ctx context.Context) {
	go func() {
		clearTicker := time.NewTicker(time.Hour) // 每小时检查一次
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.WithCtx(ctx).Infof("cache size: %d", m.cache.ItemCount())
			case <-clearTicker.C:
				m.cache.DeleteExpired()
				log.WithCtx(ctx).Infof("cache size after cleanup: %d", m.cache.ItemCount())
			}
		}
	}()
	go func() {
		buffer := make([]*types.FeedVideo, 0, 30)
		for {
			select {
			case item := <-m.feedVideoCh:
				// 检查 item 是否为 nil
				if item == nil {
					log.WithCtx(ctx).Debug("Received nil item from feedVideoCh, skipping")
					continue
				}
				// 应用过滤逻辑
				feedVideo, err := FilterVideo(item)
				if err != nil {
					// 记录但不中断其他项目的处理
					log.WithCtx(context.Background()).Debugf("Filtering failed for %s: %v", item.TorrentName, err)
					continue
				}
				if feedVideo == nil {
					continue
				}

				// 检查缓存中是否存在该 torrent name
				_, found := m.cache.Get(item.TorrentName)
				if found {
					log.WithCtx(ctx).Infof("Item %s already processed, skipping", item.TorrentName)
					continue
				}

				// 添加到缓冲区
				buffer = append(buffer, feedVideo)
				log.WithCtx(ctx).Infof("Received item %s from feedVideoCh", item.TorrentName)

				// 当缓冲区达到30个项目时进行处理
				if len(buffer) >= 30 {
					m.processFeedVideos(ctx, buffer...)
					// 重置缓冲区
					buffer = make([]*types.FeedVideo, 0, 30)
				}
			}
		}
	}()
}

func (m *MovieDB) processFeedVideos(ctx context.Context, items ...*types.FeedVideo) {
	for _, item := range items {
		// 先将所有项目加入缓存，防止重复处理
		m.cache.Set(item.TorrentName, true, 24*time.Hour)
	}

	// 对于成功过滤的项目，调用nameparser.ModelHandler进行处理
	torrentNamesToParse := make([]string, 0, len(items))
	feedVideoMap := make(map[string]*types.FeedVideo)

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
		err = m.CreatFeedVideo(feedVideo)
		if err != nil {
			if errors.Is(err, ErrDataExist) {
				log.WithCtx(ctx).Debugf("%s.%s err: %s", feedVideo.Web, feedVideo.Type, err)
				continue
			}
			log.WithCtx(ctx).Error(err)
			continue
		}

		msg := fmt.Sprintf("%s.%s: %s 保存完毕.", feedVideo.Web, feedVideo.Type, feedVideo.Name)
		if feedVideo.Type == "" {
			msg = fmt.Sprintf("%s: %s 保存完毕.", feedVideo.Web, feedVideo.Name)
		}
		log.WithCtx(ctx).Info(msg)
	}
}

func (m *MovieDB) GetDB() *gorm.DB {
	return m.db
}

func FilterVideo(feedVideoBase *types.FeedVideoBase) (*types.FeedVideo, error) {
	//  排除 低码率的视频
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
