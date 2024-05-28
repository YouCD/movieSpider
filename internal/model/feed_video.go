package model

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/youcd/toolkit/log"
	"movieSpider/internal/types"
	"strings"
	"time"
)

// FindLikeTVFromFeedVideo
//
//	@Description: 从feed_video表中查找电视剧
//	@receiver m
//	@param name
//	@return videos
//	@return err
func (m *MovieDB) FindLikeTVFromFeedVideo(name string) (videos []*types.FeedVideo, err error) {
	//nolint:sqlclosecheck, rowserrcheck
	rows, err := m.db.Model(&types.FeedVideo{}).Select(" id,name").Where(" name like ?", fmt.Sprintf("%%%s%%", name)).Rows()
	if err != nil {
		return nil, errors.WithMessage(err, "FindLikeTVFromFeedVideo")
	}

	for rows.Next() {
		v := new(types.FeedVideo)
		err = rows.Scan(&v.ID, &v.Name)
		if err != nil {
			return nil, errors.WithMessage(err, "FindLikeTVFromFeedVideo")
		}
		videos = append(videos, v)
	}
	return
}

// UpdateFeedVideoNameByID
//
//	@Description: 根据id 更新 Feed 电影名
//	@receiver m
//	@param id
//	@param name
//	@param resource
//	@return err
func (m *MovieDB) UpdateFeedVideoNameByID(id int32, name string, resource types.VideoType) (err error) {
	err = m.db.Model(&types.FeedVideo{}).Where("id=?", id).Updates(types.FeedVideo{Name: name, Type: resource.String()}).Error
	if err != nil {
		return err
	}
	return
}

// GetFeedVideoTVByName 通过 名称 获取 feedVideo tv
//
//	@Description:
//	@receiver m
//	@param names
//	@return videos
//	@return err
func (m *MovieDB) GetFeedVideoTVByName(doubanID string, names ...string) (videos []*types.FeedVideo, err error) {
	var videos1 []*types.FeedVideo
	for _, n := range names {
		log.Debugf("GetFeedVideoMovieByName 开始第一次查找tv数据: %s.", n)
		// var likeName string
		// likeName = fmt.Sprintf("%s%%", n)
		//nolint:rowserrcheck,perfsprint
		rows, err := m.db.Model(&types.FeedVideo{}).Where(`name like ? and magnet!="" and  type="tv" and download=0;`, fmt.Sprintf("%s%%", n)).Rows()
		if err != nil {
			return nil, errors.WithMessage(err, "查找失败")
		}
		defer rows.Close()
		// 只查找 没有下载过 && 类型为tv数据
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, errors.WithMessage(err, "GetFeedVideoMovieByName 扫描rows失败")
			}
			// 添加豆瓣id
			video.DoubanID = doubanID

			videos1 = append(videos1, &video)
		}
	}
	if len(videos1) > 0 {
		return videos1, nil
	}

	for _, n := range names {
		log.Debugf("GetFeedVideoMovieByName 开始第二次查找tv数据: %s.", n)
		// 查找 没有下载过 && 类型不等于TV的数据
		/*
			var likeName string
			if strings.Contains(n, ".") {
				likeName = fmt.Sprintf("%%.%s.%%", n)
			} else {
				likeName = fmt.Sprintf("%%%s%%", n)
			}
		*/
		//nolint:rowserrcheck, perfsprint
		rows, err := m.db.Model(&types.FeedVideo{}).Where(`name like ? and magnet!="" and download !=1 and type="movie" `, fmt.Sprintf("%s%%", n)).Rows()
		// rows, err := m.db.Model(&types.FeedVideo{}).Where(`name = ? and magnet!="" and download !=1 and  type !="movie" `, n).Rows()
		if err != nil {
			return nil, errors.WithMessage(err, "GetFeedVideoMovieByName 第二次查找tv数据失败")
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, errors.WithMessage(err, "解析为结构体失败")
			}
			// 添加豆瓣id
			video.DoubanID = doubanID
			videos = append(videos, &video)
		}
	}
	//nolint:nakedret
	return
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
	//nolint:rowserrcheck, sqlclosecheck
	rows, err := m.db.Model(&types.FeedVideo{}).Select("count(*)  as count ,web ").Group("web").Order("count").Rows()
	if err != nil {
		return nil, errors.WithMessage(err, "查找失败")
	}

	for rows.Next() {
		c := new(types.ReportCount)
		err = rows.Scan(&c.Count, &c.Web)
		if err != nil {
			return nil, errors.WithMessage(err, "CountFeedVideo")
		}
		counts = append(counts, c)
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
			return nil, errors.WithMessage(err, "查找失败")
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, errors.WithMessage(err, "GetFeedVideoMovieByName")
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
			return nil, errors.WithMessage(err, "查找失败")
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, errors.WithMessage(err, "GetFeedVideoMovieByName")
			}
			video.Name = n
			videos = append(videos, &video)
		}
	}
	//nolint:nakedret
	return
}

func (m *MovieDB) GetFeedVideoMovieByNameAndDoubanID(doubanID string, names ...string) (videos []*types.FeedVideo, err error) {
	var videos1 []*types.FeedVideo
	log.Debugf("GetFeedVideoMovieByName 开始第一次查找Movie数据: %s.", names)
	for _, n := range names {
		//  只查找 没有下载过 && 类型为movie数据   and download=0
		//nolint:rowserrcheck
		rows, err := m.db.Model(&types.FeedVideo{}).Where(`name = ? and magnet!="" and  type="movie" `, n).Rows()
		if err != nil {
			return nil, errors.WithMessage(err, "查找失败")
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, errors.WithMessage(err, "GetFeedVideoMovieByName")
			}
			video.DoubanID = doubanID
			videos1 = append(videos1, &video)
		}
	}
	if len(videos1) > 0 {
		return videos1, nil
	}

	log.Debugf("GetFeedVideoMovieByName 开始第二次查找Movie数据: %s.", names)
	for _, n := range names {
		// var likeName string
		// if strings.Contains(n, ".") {
		// 	likeName = fmt.Sprintf("%%.%s.%%", n)
		// } else {
		// 	likeName = fmt.Sprintf("%%%s%%", n)
		// }
		//nolint:rowserrcheck
		rows, err := m.db.Model(&types.FeedVideo{}).Where(`name  = ? and magnet!="" and download!=1  and type="movie"`, n).Rows()
		if err != nil {
			return nil, errors.WithMessage(err, "查找失败")
		}
		defer rows.Close()
		for rows.Next() {
			var video types.FeedVideo
			err = m.db.ScanRows(rows, &video)
			if err != nil {
				return nil, errors.WithMessage(err, "解析为结构体失败")
			}
			video.DoubanID = doubanID
			videos = append(videos, &video)
		}
	}
	//nolint:nakedret
	return
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
		if strings.Contains(err.Error(), "Duplicate entry") {
			log.Debugf("CreatFeedVideo 数据已存在 video: %#v", video)
			//nolint:revive
			return errors.WithMessagef(ErrDataExist, "name: %s type: %s.", video.Name, video.Type)
		}
		return errors.WithMessage(err, video.Name)
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
	video.Name = video.FormatName(video.Name)

	err = m.db.Model(&types.FeedVideo{}).Where("id=?", video.ID).Updates(video).Error
	if err != nil {
		return err
	}
	return
}
