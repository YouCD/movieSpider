package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"time"

	"github.com/youcd/toolkit/log"
	"gorm.io/gorm"
)

// CreatTMDBVideo 创建TMDB视频
//
//	@Description: 创建TMDB视频
//	@receiver m
//	@param video
//	@return err
func (m *MovieDB) CreatTMDBVideo(video *types.TMDBVideo) (err error) {
	if video == nil {
		return ErrVideoIsNil
	}
	v, err := m.FetchOneTMDBVideoByImdbID(video.ImdbID)
	if err != nil {
		// 忽略 错误信息： sql: no rows in result set
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithCtx(context.Background()).Error("video.ImdbID : %s,err: %s", video.ImdbID, err)
		}
	}

	if v != nil {
		log.WithCtx(context.Background()).Debugf("CreatTMDBVideo已存在 %#v", v)
		// 将该记录变更为 可播放
		if err = m.UpdateTMDBVideo(video); err != nil {
			log.WithCtx(context.Background()).Error(err)
		}
		return ErrDataExist
	}

	if video.Names == "null" {
		log.WithCtx(context.Background()).Errorf("CreatTMDBVideo 数据错误. video: %#v", video)
		//nolint:nakedret
		return
	}

	err = m.db.Model(&types.TMDBVideo{}).Create(video).Error

	if err != nil {
		return fmt.Errorf("CreatTMDBVideo 数据已添加. video: %#v, err: %w", video, err)
	}
	log.WithCtx(context.Background()).Debugf("CreatTMDBVideo 数据已添加. video: %#v", video)
	//nolint:nakedret
	return
}

// RandomOneTMDBVideo 随机获取一个TMDB视频
//
//	@Description: 随机获取一个TMDB视频
//	@receiver m
//	@return video
//	@return err
func (m *MovieDB) RandomOneTMDBVideo() (video *types.TMDBVideo, err error) {
	//nolint:wastedassign
	video = new(types.TMDBVideo)
	//nolint:rowserrcheck
	rows, err := m.db.Model(&types.TMDBVideo{}).Select(" id,names,imdb_id,playable").Where("imdb_id = ''").Rows()
	if err != nil {
		return nil, fmt.Errorf("RandomOneTMDBVideo, err:%w", err)
	}
	defer rows.Close()
	var videos []*types.TMDBVideo
	for rows.Next() {
		var v types.TMDBVideo
		err = rows.Scan(&v.ID, &v.Names, &v.ImdbID, &v.Playable)
		if err != nil {
			return nil, fmt.Errorf("RandomOneTMDBVideo, err:%w", err)
		}
		videos = append(videos, &v)
	}
	if len(videos) == 0 {
		return nil, ErrVideoIsNil
	}
	rand.NewSource(time.Now().UnixNano())
	//nolint:gosec
	index := rand.Intn(len(videos))
	video = videos[index]
	log.WithCtx(context.Background()).Debugf("RandomOneTMDBVideo video: %#v", video)
	return
}

// FetchOneTMDBVideoByImdbID 根据ImdbID获取TMDB视频
//
//	@Description: 根据ImdbID获取TMDB视频
//	@receiver m
//	@param ImdbID
//	@return video
//	@return err
func (m *MovieDB) FetchOneTMDBVideoByImdbID(imdbID string) (video *types.TMDBVideo, err error) {
	err = m.db.Model(&types.TMDBVideo{}).Where("imdb_id=?", imdbID).Scan(&video).Error
	if err != nil {
		return nil, err
	}
	log.WithCtx(context.Background()).Debugf("FetchOneTMDBVideoByImdbID video: %#v", video)
	return
}

// UpdateTMDBVideo 更新TMDB视频
//
//	@Description: 更新TMDB视频
//	@receiver m
//	@param video
//	@return err
func (m *MovieDB) UpdateTMDBVideo(video *types.TMDBVideo) (err error) {
	if video == nil {
		return ErrVideoIsNil
	}
	video.Timestamp = time.Now().Unix()

	err = m.db.Model(&types.TMDBVideo{}).Where("imdb_id = ?", video.ImdbID).Updates(video).Error
	if err != nil {
		return fmt.Errorf("更新失败, video: %#v, err: %w", video, err)
	}
	return
}

// FetchTMDBVideoByType 通过类型获取TMDB视频
//
//	@Description:
//	@receiver m
//	@param typ
//	@return nameList
//	@return err
func (m *MovieDB) FetchTMDBVideoByType(typ types.VideoType) (nameList map[*types.TMDBVideo][]string, err error) {
	nameList = make(map[*types.TMDBVideo][]string)

	var videos []*types.TMDBVideo
	// 时间范围限制在近一年
	end := time.Now().Unix()
	start := time.Now().AddDate(-1, 0, 0).Unix()
	result := m.db.Model(&types.TMDBVideo{}).Where("type = ? and timestamp >= ? and timestamp <= ?", typ.String(), start, end).Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, video := range videos {
		var names []string
		err = json.Unmarshal([]byte(video.Names), &names)
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			continue
		}
		var n []string
		for _, name := range names {
			if !tools.ContainsChinese(name) {
				n = append(n, name)
			}
		}
		if len(n) != 0 {
			nameList[video] = n
		}
	}
	return
}

// FetchThisYearVideo 获取今年的视频
//
//	@Description: 获取今年的视频
//	@receiver m
//	@return []types.TMDBVideo
//	@return error
func (m *MovieDB) FetchThisYearVideo() ([]*types.TMDBVideo, error) {
	thisYear := time.Now().Format("2006")
	var videos []*types.TMDBVideo
	//nolint:perfsprint
	err := m.db.Model(&types.TMDBVideo{}).Where("date_published like  ?", fmt.Sprintf("%s%%", thisYear)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
