package model

import (
	"errors"
	"fmt"
	log1 "log"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/youcd/toolkit/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MovieDB struct {
	db          *gorm.DB
	feedVideoCh chan *types.FeedVideoBase
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
			log.Error(err)
			os.Exit(1)
		}
		sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s  CHARACTER SET utf8mb4 ", config.Config.MySQL.Database)
		// 创建数据库
		err = db.Exec(sql).Error
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		newLogger := logger.New(
			log1.New(os.Stdout, "\r\n", log1.LstdFlags), // io writer

			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				//LogLevel: logger.Info, // Log level
				Colorful: true, // 禁用彩色打印
			},
		)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Host, config.Config.MySQL.Port, config.Config.MySQL.Database) // 连接数据库

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&types.FeedVideo{}, &types.DownloadHistory{}, &types.DouBanVideo{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})
	return &MovieDB{
		db,
		bus.FeedVideoChan,
	}
}

// SaveFeedVideoFromChan
//
//	@Description: 从通道中获取 feedVideo 并保存
//	@receiver m
func (m *MovieDB) SaveFeedVideoFromChan() {
	go func() {
		for {
			feedVideo, err := FilterVideo(<-m.feedVideoCh)
			if err != nil {
				log.Warn(err)
				continue
			}
			if feedVideo.Name == "" {
				log.Warnf("feedVideo.Name is empty: %v", feedVideo)
				continue
			}
			if err := NewMovieDB().CreatFeedVideo(feedVideo); err != nil {
				if errors.Is(err, ErrDataExist) {
					log.Debugf("%s.%s err: %s", strings.ToUpper(feedVideo.Web), feedVideo.Type, err)
					continue
				}
				log.Error(err)
				continue
			}
			msg := fmt.Sprintf("%s.%s: %s 保存完毕.", strings.ToUpper(feedVideo.Web), feedVideo.Type, feedVideo.Name)
			if feedVideo.Type == "" {
				msg = fmt.Sprintf("%s: %s 保存完毕.", strings.ToUpper(feedVideo.Web), feedVideo.Name)
			}
			log.Info(msg)
		}
	}()
}

func FilterVideo(feedVideoBase *types.FeedVideoBase) (*types.FeedVideo, error) {
	//  如果是空值，跳过
	if feedVideoBase == nil {
		return nil, ErrFeedVideoIsNil
	}
	//  排除 低码率的视频
	if ok := tools.ExcludeVideo(feedVideoBase.TorrentName, config.Config.ExcludeWords); ok {
		return nil, fmt.Errorf("err:%w,TorrentName:%v", ErrFeedVideoExclude, feedVideoBase.TorrentName)
	}

	feedVideo := &types.FeedVideo{
		FeedVideoBase: *feedVideoBase,
	}
	switch feedVideoBase.Web {
	case "btbt":
		feedVideo.Name = feedVideoBase.TorrentName
	default:
		// 片名 resolution year
		name, _, year, err := tools.TorrentName2info(feedVideoBase.TorrentName)
		if err != nil {
			return nil, fmt.Errorf("web:%s, err:%w", feedVideoBase.Web, err)
		}
		feedVideo.Name = name
		feedVideo.Year = year
	}
	return feedVideo, nil
}
