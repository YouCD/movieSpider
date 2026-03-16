package model

import (
	"errors"
	"fmt"
	"movieSpider/internal/types"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrVideoIsNil   = errors.New("video is nil")
	ErrHistoryIsNil = errors.New("history is nil")
)

// AddDownloadHistory
//
//	@Description: 添加下载历史
//	@receiver m
//	@param history
//	@return err
func (m *MovieDB) AddDownloadHistory(history *types.DownloadHistory) (err error) {
	if history == nil {
		return ErrHistoryIsNil
	}
	history.Timestamp = time.Now().Unix()

	// 插入数据
	err = m.db.Create(&history).Error
	if err != nil {
		//  如果重复插入，就返回nil
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}
		return err
	}
	return nil
}

// UpdateOrAddDownloadHistory
//
//	@Description: 更新或者添加下载历史
//	@receiver m
//	@param history
//	@return err
func (m *MovieDB) UpdateOrAddDownloadHistory(history *types.DownloadHistory) (err error) {
	if history == nil {
		return ErrHistoryIsNil
	}
	id, exist := m.checkDownloadHistory(history)
	// 如果存在，就更新
	if exist {
		history.Timestamp = time.Now().Unix()

		return m.db.Model(&types.DownloadHistory{}).Where("id=?", id).Updates(history).Error
	}
	// 如果不存在，就插入
	return m.AddDownloadHistory(history)
}

// checkDownloadHistory
//
//	@Description: 检查 是否已存在记录
//	@receiver m
//	@param history
//	@return flag
func (m *MovieDB) checkDownloadHistory(history *types.DownloadHistory) (id int, exist bool) {
	if history == nil {
		return 0, false
	}

	m.db.Model(&types.DownloadHistory{}).Select("id").Where("name=? and type=?  and season=? and episode=?", history.Name, history.Type, history.Season, history.Episode).Scan(&id)
	// 扫描
	if id == 0 {
		return 0, false
	}
	return id, true
}

// ShouldDownloadResult 表示过滤结果
type ShouldDownloadResult struct {
	ShouldDownload bool
	Reason         string // 过滤原因，用于日志记录
}

// FindFeedVideoInDownloadHistory
//
//	@Description: 查找已经下载过的视频（保留原有函数以兼容现有调用）
//	@receiver m
//	@param v
//	@return *types.FeedVideo
//	@return error
func (m *MovieDB) FindFeedVideoInDownloadHistory(v *types.FeedVideo) (*types.FeedVideo, error) {
	result := m.ShouldDownload(v)
	if result.ShouldDownload {
		return v, nil
	}
	return nil, fmt.Errorf("%s", result.Reason)
}

// ShouldDownload
//
//	@Description: 判断视频是否应该下载（纯查询，无副作用）
//	@receiver m
//	@param v
//	@return ShouldDownloadResult
func (m *MovieDB) ShouldDownload(v *types.FeedVideo) ShouldDownloadResult {
	if v == nil {
		return ShouldDownloadResult{ShouldDownload: false, Reason: "video is nil"}
	}

	// 将 FeedVideo 转换为 download_history
	downloadHistory := v.Convert2DownloadHistory()
	if downloadHistory == nil {
		return ShouldDownloadResult{
			ShouldDownload: false,
			Reason:         fmt.Sprintf("不能将种子转换为 downloadHistory: %s", v.TorrentName),
		}
	}

	if downloadHistory.Resolution == 0 {
		return ShouldDownloadResult{
			ShouldDownload: false,
			Reason:         fmt.Sprintf("种子名: %s, 分辨率为0, err: %w", v.TorrentName, ErrFeedVideoResolutionTooLow),
		}
	}

	// 查询历史记录
	var d types.DownloadHistory
	err := m.db.Model(&types.DownloadHistory{}).
		Where("name=? and season=? and episode=?", downloadHistory.Name, downloadHistory.Season, downloadHistory.Episode).
		First(&d).Error
	if err != nil {
		// 没有找到历史记录，应该下载
		if strings.Contains(err.Error(), "no rows in result set") || errors.Is(err, gorm.ErrRecordNotFound) {
			return ShouldDownloadResult{ShouldDownload: true, Reason: "新视频，未在下载历史中找到"}
		}
		return ShouldDownloadResult{ShouldDownload: false, Reason: fmt.Sprintf("查询数据库失败: %v", err)}
	}

	// 如果当前视频分辨率更高，应该下载
	if downloadHistory.Resolution > d.Resolution {
		return ShouldDownloadResult{
			ShouldDownload: true,
			Reason:         fmt.Sprintf("找到更高分辨率版本: %d -> %d", d.Resolution, downloadHistory.Resolution),
		}
	}

	// 已下载过相同或更高分辨率版本
	return ShouldDownloadResult{
		ShouldDownload: false,
		Reason:         fmt.Sprintf("已下载过相同或更高分辨率: 当前%d, 历史%d, 种子名: %s", downloadHistory.Resolution, d.Resolution, v.TorrentName),
	}
}
