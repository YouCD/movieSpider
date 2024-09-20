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
		name, _, year, err := torrentName2info(feedVideoBase.TorrentName)
		if err != nil {
			return nil, fmt.Errorf("web:%s, err:%w", feedVideoBase.Web, err)
		}
		feedVideo.Name = name
		feedVideo.Year = year
	}
	return feedVideo, nil
}

func torrentName2info(torrentName string) (string, string, string, error) {
	// 去除 -
	newTorrentName := strings.ReplaceAll(torrentName, "-", ".")
	// 去除 _
	newTorrentName = strings.ReplaceAll(newTorrentName, "_", ".")
	newTorrentName = strings.ReplaceAll(newTorrentName, ",", "")

	// 去除空格
	reg := regexp.MustCompile(`( )+|(\n)+`)
	newTorrentName = reg.ReplaceAllString(newTorrentName, "$1$2")
	newTorrentName = strings.ReplaceAll(newTorrentName, " ", ".")
	newTorrentName = strings.ReplaceAll(newTorrentName, ".", ".") //nolint:gocritic

	dotReg := regexp.MustCompile(`\.+`)
	newTorrentName = dotReg.ReplaceAllString(newTorrentName, ".")

	// 去除 []
	newTorrentName = strings.ReplaceAll(newTorrentName, "[", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, "]", "")

	// 去除 ()
	newTorrentName = strings.ReplaceAll(newTorrentName, "(", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, ")", "")

	var name, resolution, year string
	// 先匹配 tv
	//nolint:revive
	if tvName, tvNameArr, err := matchTV(newTorrentName); err != nil {
		log.Debugf("tv匹配失败: old:%s, new:%s , tvNameArr: %s", torrentName, newTorrentName, tvNameArr)
		goto MatchMovie
	} else {
		name = tvName
		goto MatchResolution
	}
	//	 匹配 movie
MatchMovie:
	if movieName, nameRegexArr, err := matchMovie(newTorrentName); err == nil {
		name = movieName
	} else {
		log.Debugf("movie匹配失败: old:%s, new:%s, Arr: %s", torrentName, newTorrentName, nameRegexArr)
		return "", "", "", fmt.Errorf("%w,old:%s, new:%s,  Arr: %s", err, torrentName, newTorrentName, nameRegexArr)
	}
MatchResolution:
	// 匹配分辨率
	{
		resolutionRegex := regexp.MustCompile(`.*\.(\d{4}[p|P])\.`)
		resolutionArr := resolutionRegex.FindStringSubmatch(newTorrentName)
		if len(resolutionArr) < 2 || len(resolutionArr) == 0 {
			return "", "", "", fmt.Errorf("%w, old:%s, new:%s, Arr: %s", ErrFeedVideoResolution, torrentName, newTorrentName, resolutionArr)
		}
		resolution = resolutionArr[1]
	}

	// 匹配 年份
	{
		compileYearRegex := regexp.MustCompile(`.*?(\d{4}).*?\d{4}[p|P].*`)
		yearArr := compileYearRegex.FindStringSubmatch(newTorrentName)
		if len(yearArr) < 2 || len(yearArr) == 0 {
			return "", "", "", fmt.Errorf("%w, old:%s, new:%s, Arr: %s", ErrFeedVideoYear, torrentName, newTorrentName, yearArr)
		}
		year = yearArr[1]
	}

	return name, resolution, year, nil
}

func matchMovie(torrentName string) (string, []string, error) {
	nameRegex := regexp.MustCompile(`(.*)\.\d{4}[p|P]`)
	nameRegexArr := nameRegex.FindStringSubmatch(torrentName)
	if len(nameRegexArr) >= 2 {
		name := nameRegexArr[1]
		movieA, arr, err := matchMovieA(name)
		if err == nil {
			return movieA, arr, nil
		}
		return name, nil, nil
	}
	return matchMovieA(torrentName)
}

func matchMovieA(torrentName string) (string, []string, error) {
	nameReg := regexp.MustCompile(`(.*)\.\d{4}`)
	nameSubMatch := nameReg.FindStringSubmatch(torrentName)
	if len(nameSubMatch) >= 2 {
		return nameSubMatch[1], nil, nil
	}
	return "", nameSubMatch, ErrFeedVideoMovieMatch
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
