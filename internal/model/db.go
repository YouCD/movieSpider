package model

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	log1 "log"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"os"
	"strings"
	"sync"
	"time"
)

type MovieDB struct {
	db          *gorm.DB
	feedVideoCh chan *types.FeedVideo
}

//nolint:gochecknoglobals
var (
	once         sync.Once
	db           = new(gorm.DB)
	ErrDataExist = errors.New("数据已存在")
	err          error
)

func NewMovieDB() *MovieDB {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Host, config.Config.MySQL.Port)
		//nolint:exhaustruct
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
			//nolint:exhaustruct
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				//LogLevel: logger.Info, // Log level
				Colorful: true, // 禁用彩色打印
			},
		)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Config.MySQL.User, config.Config.MySQL.Password, config.Config.MySQL.Host, config.Config.MySQL.Port, config.Config.MySQL.Database) // 连接数据库
		//nolint:exhaustruct
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		//nolint:exhaustruct
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
			feedVideo := <-m.feedVideoCh
			//  如果是空值，跳过
			if feedVideo == nil {
				continue
			}
			//  排除 低码率的视频
			if ok := tools.ExcludeVideo(feedVideo.TorrentName, config.Config.ExcludeWords); ok {
				continue
			}
			// log.Infof("%s.%s: %s 开始保存.", strings.ToUpper(feedVideo.Web), feedVideo.Type, feedVideo.Name)
			if err := NewMovieDB().CreatFeedVideo(feedVideo); err != nil {
				if errors.Is(err, ErrDataExist) {
					log.Debugf("%s.%s err: %s", strings.ToUpper(feedVideo.Web), feedVideo.Type, err)
					continue
				}
				log.Error(err)
				continue
			}
			log.Infof("%s.%s: %s 保存完毕.", strings.ToUpper(feedVideo.Web), feedVideo.Type, feedVideo.Name)
		}
	}()
}
