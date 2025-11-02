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
			item := <-m.feedVideoCh
			// 检查 item 是否为 nil
			if item == nil {
				log.Debug("Received nil item from feedVideoCh, skipping")
				continue
			}
			feedVideo, err := FilterVideo(item)
			if err != nil {
				log.Debugf("web:%s,err:%s", item.Web, err)
				continue
			}

			// 检查 feedVideo 是否为 nil
			if feedVideo == nil {
				log.Debug("Filtered feedVideo is nil, skipping")
				continue
			}

			err = NewMovieDB().CreatFeedVideo(feedVideo)
			if err != nil {
				if errors.Is(err, ErrDataExist) {
					log.Debugf("%s.%s err: %s", feedVideo.Web, feedVideo.Type, err)
					continue
				}
				log.Error(err)
				continue
			}
			msg := fmt.Sprintf("%s.%s: %s 保存完毕.", feedVideo.Web, feedVideo.Type, feedVideo.Name)
			if feedVideo.Type == "" {
				msg = fmt.Sprintf("%s: %s 保存完毕.", feedVideo.Web, feedVideo.Name)
			}
			log.Info(msg)
		}
	}()
}
func (m *MovieDB) GetDB() *gorm.DB {
	return m.db
}

func FilterVideo(feedVideoBase *types.FeedVideoBase) (*types.FeedVideo, error) {
	//  如果是空值，跳过
	if feedVideoBase == nil {
		return nil, ErrFeedVideoIsNil
	}
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
		log.Error(err)
	}
	if video != nil {
		return nil, fmt.Errorf("torrent_name:%s, err:%w", feedVideo.TorrentName, ErrFeedVideoExist)
	}

	// 使用模型解析种子名
	typeStr, newName, year, resolution, err := nameparser.ModelHandler(context.Background(), feedVideo.TorrentName)
	if err != nil {
		log.Warnf("TorrentName: %#v,err: %s", feedVideo.TorrentName, err)
		return nil, fmt.Errorf("FilterVideo err: %w", err)
	}
	if resolution == "" {
		//nolint:err113
		return nil, fmt.Errorf("feedVideo.TorrentName: %#v,resolution is empty", feedVideo.TorrentName)
	}
	if len([]rune(newName)) > len([]rune(feedVideo.TorrentName)) {
		log.Error("nameParser err: %s", feedVideo.TorrentName)
		return nil, ErrNameParser
	}

	feedVideo.Name = newName
	feedVideo.Year = year
	feedVideo.Type = typeStr
	log.Infow("nameParser", "input", feedVideo.TorrentName, "type", typeStr, "name", newName, "year", year, "resolution", resolution)
	return feedVideo, nil
}
