package model

import (
	"errors"
	"fmt"
	"movieSpider/internal/types"
	"strings"
	"time"

	"github.com/youcd/toolkit/log"
)

// FindLikeTVFromFeedVideo
//
//	@Description: 从feed_video表中查找电视剧
//	@receiver m
//	@param name
//	@return videos
//	@return err
func (m *MovieDB) FindLikeTVFromFeedVideo(name string) (videos []*types.FeedVideo, err error) {
	if err := m.db.Model(&types.FeedVideo{}).Select(" id,name").Where(" name like ?", fmt.Sprintf("%%%s%%", name)).Find(&videos).Error; err != nil {
		return nil, fmt.Errorf("FindLikeTVFromFeedVideo,err %w", err)
	}
	return
}

// GetFeedVideoTVByNames 通过 名称 获取 feedVideo tv
//
//	@Description:
//	@receiver m
//	@param names
//	@return videos
//	@return err
func (m *MovieDB) GetFeedVideoTVByNames(names ...string) ([]*types.FeedVideo, error) {
	log.Debugf("开始第一次查找tv数据: %s.", names)
	firstQuery := `name like ? and magnet!="" and  type="tv" and download=0;`
	argsFunc := func(name string) string {
		return fmt.Sprintf("%%%s%%", name)
	}
	videos, err := m.findTV(names, firstQuery, argsFunc)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		return videos, nil
	}

	log.Debugf("开始第二次查找tv数据: %s.", names)
	argsFunc = func(name string) string {
		return fmt.Sprintf("%s%%", name)
	}
	secondQuery := `name like ? and magnet!="" and download !=1 and type="movie"`
	videos, err = m.findTV(names, secondQuery, argsFunc)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		return videos, nil
	}

	log.Debugf("开始第三次查找tv数据: %s.", names)
	argsFunc = func(name string) string {
		return fmt.Sprintf("%s%%", name)
	}
	thirdQuery := `name like ? and magnet!="" and download !=1`
	videos, err = m.findTV(names, thirdQuery, argsFunc)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		return videos, nil
	}

	return videos, ErrTVNotFound
}

var (
	ErrMoviesNotFound = errors.New("Movies Not Found")
	ErrTVNotFound     = errors.New("TV Not Found")
)

func (m *MovieDB) findTV(names []string, query string, argsFunc func(name string) string) ([]*types.FeedVideo, error) {
	var firstVideos []*types.FeedVideo
	log.Debugf("GetFeedVideoMovieByName 开始第一次查找tv数据: %s.", names)
	//nolint:perfsprint
	for _, n := range names {
		if err := m.db.Model(&types.FeedVideo{}).Where(query, argsFunc(n)).Find(&firstVideos).Error; err != nil {
			return nil, fmt.Errorf("查找失败, err:%w", err)
		}
	}
	return firstVideos, nil
}

// UpdateFeedVideoDownloadByID
//
//	@Description: 根据id 更新下载状态
//	@receiver m
//	@param id
//	@param isDownload
//	@return err
func (m *MovieDB) UpdateFeedVideoDownloadByID(id int32, isDownload int) (err error) {
	// 定义sql
	err = m.db.Model(&types.FeedVideo{}).Where("id=?", id).Updates(types.FeedVideo{Download: isDownload}).Error
	if err != nil {
		return err
	}
	return
}

// CountFeedVideo
//
//	@Description: 统计feed_video表中的数据
//	@receiver m
//	@return counts
//	@return err
func (m *MovieDB) CountFeedVideo() (counts []*types.ReportCount, err error) {
	err = m.db.Model(&types.FeedVideo{}).Select("count(*)  as count ,web ").Group("web").Order("count").Find(&counts).Error
	if err != nil {
		return nil, fmt.Errorf("查找失败, err:%w", err)
	}
	return
}

// GetFeedVideoMovieByName 通过 名称 获得 feedVideo movie
//
//	@Description:
//	@receiver m
//	@param names
//	@return videos
//	@return err
func (m *MovieDB) GetFeedVideoMovieByName(names ...string) (videos []*types.FeedVideo, err error) {
	var videos1 []*types.FeedVideo
	log.Debugf("GetFeedVideoMovieByName 开始第一次查找Movie数据: %s.", names)
	for _, n := range names {
		// var likeName string
		// likeName = fmt.Sprintf("%%%s%%", n)
		// 只查找 没有下载过 && 类型为movie数据   and download=0
		//nolint:rowserrcheck
		rows, err := m.db.Model(&types.FeedVideo{}).Where(`name = ? and magnet!="" and  type="movie" `, n).Rows()
		if err != nil {
			return nil, fmt.Errorf("查找失败, err:%w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, fmt.Errorf("GetFeedVideoMovieByName, err:%w", err)
			}
			// 将feedVideo的名称设置为搜索的名称
			video.Name = n
			videos1 = append(videos1, &video)
		}
	}
	if len(videos1) > 0 {
		return videos1, nil
	}

	for _, n := range names {
		// 查找 没有下载过 && 类型不等于TV的数据

		/*
			var likeName string
			if strings.Contains(n, ".") {
				likeName = fmt.Sprintf("%%.%s.%%", n)
			} else {
				likeName = fmt.Sprintf("%%%s%%", n)
			}
		*/
		//nolint:rowserrcheck
		rows, err := m.db.Model(&types.FeedVideo{}).Where(`name like ? and magnet!="" and download=0 and type="movie"`, n).Rows()
		// rows, err := m.db.Model(&types.FeedVideo{}).Where(`name like ? and magnet!="" and download=0 and type!="tv"`, n).Rows()
		if err != nil {
			return nil, fmt.Errorf("查找失败, err:%w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, fmt.Errorf("GetFeedVideoMovieByName, err:%w", err)
			}
			video.Name = n
			videos = append(videos, &video)
		}
	}
	//nolint:nakedret
	return
}

func (m *MovieDB) GetFeedVideoMovieByNames(names ...string) ([]*types.FeedVideo, error) {
	log.Debugf("开始第一次查找Movie数据: %s.", names)
	firstQuery := `name = ? and magnet!="" and type="movie"`
	videos, err := m.findMovie(names, firstQuery)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		log.Debugf("%#v 种子数:%d", names, len(videos))
		return videos, nil
	}
	log.Debugf("开始第二次查找Movie数据: %s.", names)

	secondQuery := `name = ? and magnet!="" and download!=1 and type="movie"`
	videos, err = m.findMovie(names, secondQuery)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		log.Debugf("%#v 种子数:%d", names, len(videos))
		return videos, nil
	}

	log.Debugf("开始第三次查找Movie数据: %s.", names)
	thirdQuery := `name = ? and magnet!=""`
	videos, err = m.findMovie(names, thirdQuery)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		log.Debugf("%#v 种子数:%d", names, len(videos))
		return videos, nil
	}
	return videos, ErrMoviesNotFound
}

func (m *MovieDB) findMovie(names []string, query string) ([]*types.FeedVideo, error) {
	var movies []*types.FeedVideo
	for _, n := range names {
		//  只查找 没有下载过 && 类型为movie数据   and download=0
		if err := m.db.Model(&types.FeedVideo{}).Where(query, n).Find(&movies).Error; err != nil {
			return nil, fmt.Errorf("查找失败, err:%w", err)
		}
	}
	if len(movies) > 0 {
		return movies, nil
	}
	return movies, ErrMoviesNotFound
}

// CreatFeedVideo
//
//	@Description: 创建feed视频
//	@receiver m
//	@param video
//	@return err
func (m *MovieDB) CreatFeedVideo(video *types.FeedVideo) (err error) {
	if video.Magnet == "" {
		//nolint:goerr113
		return fmt.Errorf("CreatFeedVideo Magnet is nill : %#v", video)
	}
	video.Timestamp = time.Now().Unix()
	video.RowData.Valid = true

	err = m.db.Model(types.FeedVideo{}).Create(video).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062") {
			log.Debugf("CreatFeedVideo 数据已存在 video: %#v", video)
			return fmt.Errorf("name: %s type: %s. err:%w", video.Name, video.Type, ErrDataExist)
		}
		return fmt.Errorf("%s err:%w", video.Name, err)
	}
	log.Debugf("CreatFeedVideo 数据已添加 video: %#v", video)
	return
}

// UpdateFeedVideo
//
//	@Description: 更新feed视频 所有字段
//	@receiver m
//	@param video
//	@return err
func (m *MovieDB) UpdateFeedVideo(video *types.FeedVideo) (err error) {
	err = m.db.Model(&types.FeedVideo{}).Where("id=?", video.ID).Updates(video).Error
	if err != nil {
		return err
	}
	return
}
