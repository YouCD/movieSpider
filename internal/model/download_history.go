package model

import (
	"context"
	"errors"
	"fmt"
	"movieSpider/internal/types"
	"strings"
	"time"

	"github.com/youcd/toolkit/log"
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

// FindFeedVideoInDownloadHistory
//
//	@Description: 查找已经下载过的视频
//	@receiver m
//	@param v
//	@return *types.FeedVideo
//	@return error
func (m *MovieDB) FindFeedVideoInDownloadHistory(v *types.FeedVideo) (*types.FeedVideo, error) {
	if v == nil {
		return nil, ErrVideoIsNil
	}

	//  将 FeedVideo 转换为 download_history
	downloadHistory := v.Convert2DownloadHistory()
	if downloadHistory.Resolution == 0 {
		return nil, fmt.Errorf("种子名: %s,分辨率: %d, err:%w", v.TorrentName, downloadHistory.Resolution, ErrFeedVideoResolutionTooLow)
	}

	if downloadHistory == nil {
		//nolint:err113
		return nil, fmt.Errorf("不能将种子: %#v 转换为 downloadHistory", v.TorrentName)
	}

	// 查找
	var d *types.DownloadHistory
	err := m.db.Model(&types.DownloadHistory{}).Where("name=? and season=? and episode=?", downloadHistory.Name, downloadHistory.Season, downloadHistory.Episode).Scan(&d).Error
	if err != nil {
		// log.Error(downloadHistory.TorrentName, err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return v, nil
		}
		return nil, err
	}
	//  如果没有找到就直接保存
	if d == nil {
		// 这里 不管有没有错误，都直接返回 d
		// log.Errorf("%#v", downloadHistory)
		err = m.AddDownloadHistory(downloadHistory)
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
		}
		return v, nil
	}
	// log.Errorf("downloadHistory.Name: %s, d.Name: %s, downloadHistory.Resolution: %d, d.Resolution: %d", downloadHistory.Name, d.Name, downloadHistory.Resolution, d.Resolution)
	// 如果 查找的 video 的分辨率小于 download_history 的分辨率，就不用下载，返回 nil
	if downloadHistory.Resolution <= d.Resolution {
		//nolint:err113
		return nil, fmt.Errorf("种子名: %s, 分辨率: %d ,已经下载过相同分辨率，或分辨率小于已经下载的种子", downloadHistory.Name, downloadHistory.Resolution)
	}

	return v, nil
}
