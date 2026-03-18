package model

import (
	"context"
	"errors"
	"fmt"
	"movieSpider/internal/tools"
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
	result := m.db.Model(&types.FeedVideo{}).Select("id,name").Where(" name like ?", fmt.Sprintf("%%%s%%", name)).Find(&videos)
	if result.Error != nil {
		return nil, fmt.Errorf("FindLikeTVFromFeedVideo,err %w", err)
	}
	return
}
func (m *MovieDB) GetFeedVideoByName(name string) (*types.FeedVideo, error) {
	var video *types.FeedVideo
	err := m.db.Model(&types.FeedVideo{}).Where("torrent_name = ?", name).First(&video).Error
	if err != nil {
		return nil, err
	}
	return video, err
}

// GetFeedVideoTVByNames 通过 名称 获取 feedVideo tv
//
//	@Description:
//	@receiver m
//	@param names
//	@return videos
//	@return err
func (m *MovieDB) GetFeedVideoTVByNames(ctx context.Context, names ...string) ([]*types.FeedVideo, error) {
	log.WithCtx(ctx).Debugf("开始第一次查找tv数据: %s.", names)
	firstQuery := `name like ? and magnet!="" and  type="tv" and download=0;`
	argsFunc := func(name string) string {
		return fmt.Sprintf("%%%s%%", name)
	}
	videos, err := m.findTV(ctx, names, firstQuery, argsFunc)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		return videos, nil
	}

	log.WithCtx(ctx).Debugf("开始第二次查找tv数据: %s.", names)
	argsFunc = func(name string) string {
		return name + "%"
	}
	secondQuery := `name like ? and magnet!="" and download !=1 and type="movie"`
	videos, err = m.findTV(ctx, names, secondQuery, argsFunc)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		return videos, nil
	}

	log.WithCtx(ctx).Debugf("开始第三次查找tv数据: %s.", names)
	argsFunc = func(name string) string {
		return name + "%"
	}
	thirdQuery := `name like ? and magnet!="" and download !=1`
	videos, err = m.findTV(ctx, names, thirdQuery, argsFunc)
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

func (m *MovieDB) findTV(ctx context.Context, names []string, query string, argsFunc func(name string) string) ([]*types.FeedVideo, error) {
	var firstVideos []*types.FeedVideo
	log.WithCtx(ctx).Debugf("GetFeedVideoMovieByName 开始第一次查找tv数据: %s.", names)
	for _, n := range names {
		if tools.ContainsChinese(n) {
			continue
		}
		log.WithCtx(ctx).Debugf("findTV 获取数据: %s.", n)
		result := m.db.Model(&types.FeedVideo{}).WithContext(ctx).Where(query, argsFunc(n)).Find(&firstVideos)
		if result.Error != nil {
			return nil, fmt.Errorf("查找失败, err:%w", result.Error)
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

func (m *MovieDB) GetFeedVideoMovieByNames(ctx context.Context, names ...string) ([]*types.FeedVideo, error) {
	log.WithCtx(ctx).Debugf("开始第一次查找Movie数据: %s.", names)
	firstQuery := `name = ? and magnet!="" and type="movie"`
	videos, err := m.findMovie(ctx, names, firstQuery)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		log.WithCtx(ctx).Debugf("%#v 种子数:%d", names, len(videos))
		return videos, nil
	}
	log.WithCtx(ctx).Debugf("开始第二次查找Movie数据: %s.", names)

	secondQuery := `name = ? and magnet!="" and download!=1 and type="movie"`
	videos, err = m.findMovie(ctx, names, secondQuery)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		log.WithCtx(ctx).Debugf("%#v 种子数:%d", names, len(videos))
		return videos, nil
	}

	log.WithCtx(ctx).Debugf("开始第三次查找Movie数据: %s.", names)
	thirdQuery := `name = ? and magnet!=""`
	videos, err = m.findMovie(ctx, names, thirdQuery)
	if err != nil {
		return videos, err
	}
	if len(videos) > 0 {
		log.WithCtx(ctx).Debugf("%#v 种子数:%d", names, len(videos))
		return videos, nil
	}
	return videos, ErrMoviesNotFound
}

func (m *MovieDB) findMovie(ctx context.Context, names []string, query string) ([]*types.FeedVideo, error) {
	var movies []*types.FeedVideo
	for _, n := range names {
		//  只查找 没有下载过 && 类型为movie数据   and download=0
		result := m.db.Model(&types.FeedVideo{}).WithContext(ctx).Where(query, n).Find(&movies)
		if result.Error != nil {
			return nil, fmt.Errorf("查找失败, err:%w", result.Error)
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
func (m *MovieDB) CreatFeedVideo(ctx context.Context, video *types.FeedVideo) (err error) {
	if video.Magnet == "" {
		//nolint:err113
		return fmt.Errorf("CreatFeedVideo Magnet is nill : %#v", video)
	}
	video.Timestamp = time.Now().Unix()
	video.RowData.Valid = true

	err = m.db.Model(types.FeedVideo{}).WithContext(ctx).Create(video).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062") {
			log.WithCtx(ctx).Debugf("CreatFeedVideo 数据已存在 video: %#v", video)
			return fmt.Errorf("name: %s type: %s. err:%w", video.Name, video.Type, ErrDataExist)
		}
		return fmt.Errorf("%s err:%w", video.Name, err)
	}
	log.WithCtx(ctx).Debugf("CreatFeedVideo 数据已添加 video.TorrentName: %s", video.TorrentName)
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

// UpdateFeedVideos
//
//	@Description: 批量更新feed视频 所有字段
//	@receiver m
//	@param videos
//	@return err
func (m *MovieDB) UpdateFeedVideos(ctx context.Context, videos ...*types.FeedVideo) (err error) {
	if len(videos) == 0 {
		return nil
	}
	log.WithCtx(ctx).Debugw("MovieDB", "数据量", len(videos))

	// 使用事务确保数据一致性
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	for _, video := range videos {
		result := tx.Model(&types.FeedVideo{}).Where("id=?", video.ID).Updates(video)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	return tx.Commit().Error
}
func (m *MovieDB) CountFeedVideoByToday() (counts []*types.ReportCount, err error) {
	// select count(*) as cnt, web from feed_video  where timestamp>=1767801600 group by web  ;
	// 获取当前时间
	now := time.Now()
	// 获取今天的年月日
	year, month, day := now.Date()
	// 创建今天00:00:00的时间对象
	todayStart := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	// 获取Unix时间戳
	timestamp := todayStart.Unix()
	err = m.db.Model(&types.FeedVideo{}).Select("count(*) as count, web").Where("timestamp>=?", timestamp).Group("web").Find(&counts).Error
	if err != nil {
		return nil, fmt.Errorf("查找失败, err:%w", err)
	}
	return
}
