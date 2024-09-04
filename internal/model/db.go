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
	"regexp"
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
	err          error
)
var (
	ErrNotMatchTorrentName = errors.New("torrent name not match")
	ErrFeedVideoIsNil      = errors.New("feedVideo is nil")
	ErrFeedVideoExclude    = errors.New("feedVideo exclude")
	ErrFeedVideoResolution = errors.New("feedVideo resolution match")
	ErrFeedVideoYear       = errors.New("feedVideo year match")
)

func NewMovieDB() *MovieDB {
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
				log.Error(err)
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

func FilterVideo(feedVideoBase *types.FeedVideoBase) (*types.FeedVideo, error) {
	//  如果是空值，跳过
	if feedVideoBase == nil {
		return nil, ErrFeedVideoIsNil
	}
	//  排除 低码率的视频
	if ok := tools.ExcludeVideo(feedVideoBase.TorrentName, config.Config.ExcludeWords); ok {
		return nil, ErrFeedVideoExclude
	}

	feedVideo := &types.FeedVideo{
		FeedVideoBase: *feedVideoBase,
	}
	switch feedVideoBase.Web {
	case "btbt":
		feedVideo.Name = feedVideoBase.TorrentName
	default:
		// 片名 resolution year
		name, _, year, err := torrentName2info(feedVideoBase.TorrentName)
		if err != nil {
			return nil, err
		}
		feedVideo.Name = name
		feedVideo.Year = year
	}

	// 处理 电影 名字
	switch feedVideoBase.Type {
	case types.VideoTypeMovie.String():

	case types.VideoTypeTV.String():

	default:
		log.Warn("未知类型:", feedVideoBase)
	}

	//// log.Infof("%s.%s: %s 开始保存.", strings.ToUpper(feedVideoBase.Web), feedVideoBase.Type, feedVideoBase.Name)
	// if err := NewMovieDB().CreatFeedVideo(feedVideoBase); err != nil {
	// 	if errors.Is(err, ErrDataExist) {
	// 		log.Debugf("%s.%s err: %s", strings.ToUpper(feedVideoBase.Web), feedVideoBase.Type, err)
	// 		continue
	// 	}
	// 	log.Error(err)
	// 	continue
	// }
	// log.Infof("%s.%s: %s 保存完毕.", strings.ToUpper(feedVideoBase.Web), feedVideoBase.Type, feedVideoBase.Name)
	return feedVideo, nil
}

//nolint:ineffassign,wastedassign
func torrentName2info(torrentName string) (string, string, string, error) {
	// 去除 -
	newTorrentName := strings.ReplaceAll(torrentName, "-", ".")
	// 去除 _
	newTorrentName = strings.ReplaceAll(newTorrentName, "_", ".")

	// 去除空格
	reg := regexp.MustCompile(`( )+|(\n)+`)
	newTorrentName = reg.ReplaceAllString(newTorrentName, "$1$2")
	newTorrentName = strings.ReplaceAll(newTorrentName, " ", ".")
	newTorrentName = strings.ReplaceAll(newTorrentName, ".", ".") //nolint:gocritic

	// 去除 []
	newTorrentName = strings.ReplaceAll(newTorrentName, "[", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, "]", "")

	// 去除 ()
	newTorrentName = strings.ReplaceAll(newTorrentName, "(", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, ")", "")

	var name, resolution, year string
	// 先匹配 tv
	{
		if tvName, tvNameArr, err := matchTV(newTorrentName); err != nil {
			log.Debugf("tv匹配失败: old:%s, new:%s , tvNameArr: %s", torrentName, newTorrentName, tvNameArr)
		} else {
			name = tvName
			// 匹配 片名
			if movieName, nameRegexArr, err := matchMovie(newTorrentName); err == nil {
				name = movieName
			} else {
				log.Debugf("movie匹配失败: old:%s, new:%s , movieNameArr: %s", torrentName, newTorrentName, nameRegexArr)
				return "", "", "", fmt.Errorf("old:%s, new:%s movieNam匹配失败:%w, movieNameArr: %s", torrentName, newTorrentName, err, nameRegexArr)
			}
		}
	}

	// 匹配分辨率
	{
		resolutionRegex := regexp.MustCompile(`.*\.(1080[p|P]|2106[p|P])\.`)
		resolutionArr := resolutionRegex.FindStringSubmatch(newTorrentName)
		if len(resolutionArr) < 2 || len(resolutionArr) == 0 {
			return "", "", "", fmt.Errorf("old:%s, new:%s resolution匹配失败:%w, resolutionArr: %s", torrentName, newTorrentName, ErrFeedVideoResolution, resolutionArr)
		}
		resolution = resolutionArr[1]
	}

	// 匹配 年份
	{
		compileYearRegex := regexp.MustCompile(`.*?(\d{4}).*?1080[p|P]|2106[p|P].*`)
		yearArr := compileYearRegex.FindStringSubmatch(newTorrentName)
		if len(yearArr) < 2 || len(yearArr) == 0 {
			return "", "", "", fmt.Errorf("old:%s, new:%s yearArr匹配失败:%w, yearArr: %s", torrentName, newTorrentName, ErrFeedVideoYear, yearArr)
		}
		year = yearArr[1]
	}

	return name, resolution, year, nil
}

func matchMovie(torrentName string) (string, []string, error) {
	nameRegex := regexp.MustCompile(`(.*)\.1080[p|P]|2106[p|P]\.`)
	nameRegexArr := nameRegex.FindStringSubmatch(torrentName)
	if len(nameRegexArr) >= 2 {
		name := nameRegexArr[1]
		// Longlegs.2024
		nameReg := regexp.MustCompile(`(.*)\.\d{4}`)
		nameSubMatch := nameReg.FindStringSubmatch(name)
		if len(nameSubMatch) >= 2 {
			name = nameSubMatch[1]
			return name, nil, nil
		}
		return name, nil, nil
	}
	return "", nameRegexArr, ErrNotMatchTorrentName
}

func matchTV(torrentName string) (string, []string, error) {
	tvReg := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9][eE][0-9][0-9])`)
	tvNameArr := tvReg.FindStringSubmatch(torrentName)
	// 如果 正则匹配过后 没有结果直接 过滤掉
	if len(tvNameArr) < 2 || len(tvNameArr) == 0 {
		tvRegA := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9]).*`)
		tvNameArrA := tvRegA.FindStringSubmatch(torrentName)
		if len(tvNameArrA) >= 2 {
			return tvNameArrA[1], nil, nil
		}
		return "", tvNameArrA, ErrNotMatchTorrentName
	}
	return tvNameArr[1], nil, nil
}
