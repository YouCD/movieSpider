package model

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	log1 "log"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"os"
	"sync"
	"time"
)

type movieDB struct {
	db *gorm.DB
}

var (
	once           sync.Once
	db             = new(gorm.DB)
	ErrorDataExist = errors.New("数据已存在")
	err            error
)

func NewMovieDB() *movieDB {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local", config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s  CHARACTER SET utf8mb4 ", config.MySQL.Database)
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
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.Database) // 连接数据库

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
	return &movieDB{
		db,
	}

}
