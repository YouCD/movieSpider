package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"movieSpider/internal/types"
	"time"

	"github.com/youcd/toolkit/log"
)

// CreatDouBanVideo
//
//	@Description: 创建豆瓣视频
//	@receiver m
//	@param video
//	@return err
func (m *MovieDB) CreatDouBanVideo(video *types.DouBanVideo) (err error) {
	if video == nil {
		return ErrVideoIsNil
	}
	v, err := m.FetchOneDouBanVideoByDouBanID(video.DoubanID)
	if err != nil {
		// 忽略 错误信息： sql: no rows in result set
		if !errors.Is(sql.ErrNoRows, err) {
			log.Error("video.DoubanID : %s,err: %s", video.DoubanID, err)
		}
	}

	if v != nil {
		log.Debugf("CreatDouBanVideo已存在 %#v", v)
		// 将该记录变更为 可播放
		err = m.UpdateDouBanVideo(video)
		if err != nil {
			log.Error(err)
		}
		return ErrDataExist
	}

	if video.Names == "null" {
		log.Errorf("CreatDouBanVideo 数据错误. video: %#v", video)
		//nolint:nakedret
		return
	}

	err = m.db.Model(&types.DouBanVideo{}).Create(video).Error

	if err != nil {
		return fmt.Errorf("CreatDouBanVideo 数据已添加. video: %#v, err: %w", video, err)
	}
	log.Debugf("CreatDouBanVideo 数据已添加. video: %#v", video)
	//nolint:nakedret
	return
}

// RandomOneDouBanVideo
//
//	@Description: 随机获取一个豆瓣视频
//	@receiver m
//	@return video
//	@return err
func (m *MovieDB) RandomOneDouBanVideo() (video *types.DouBanVideo, err error) {
	//nolint:wastedassign
	video = new(types.DouBanVideo)
	//nolint:rowserrcheck
	rows, err := m.db.Model(&types.DouBanVideo{}).Select(" id,names,douban_id,playable").Where("imdb_id = ''").Rows()
	if err != nil {
		return nil, fmt.Errorf("RandomOneDouBanVideo, err:%w", err)
	}
	defer rows.Close()
	var videos []*types.DouBanVideo
	for rows.Next() {
		var v types.DouBanVideo
		err = rows.Scan(&v.ID, &v.Names, &v.DoubanID, &v.Playable)
		if err != nil {
			return nil, fmt.Errorf("RandomOneDouBanVideo, err:%w", err)
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
	log.Debugf("RandomOneDouBanVideo video: %#v", video)
	return
}

// FetchOneDouBanVideoByDouBanID
//
//	@Description: 根据豆瓣ID获取豆瓣视频
//	@receiver m
//	@param DouBanID
//	@return video
//	@return err
func (m *MovieDB) FetchOneDouBanVideoByDouBanID(douBanID string) (video *types.DouBanVideo, err error) {
	err = m.db.Model(&types.DouBanVideo{}).Where("douban_id=?", douBanID).Scan(&video).Error
	if err != nil {
		return nil, err
	}
	log.Debugf("FetchOneDouBanVideoByDouBanID video: %#v", video)
	return
}

// UpdateDouBanVideo
//
//	@Description: 更新豆瓣视频
//	@receiver m
//	@param video
//	@return err
func (m *MovieDB) UpdateDouBanVideo(video *types.DouBanVideo) (err error) {
	if video == nil {
		return ErrVideoIsNil
	}
	video.Timestamp = time.Now().Unix()

	err = m.db.Model(&types.DouBanVideo{}).Where("douban_id = ?", video.DoubanID).Updates(video).Error
	if err != nil {
		return fmt.Errorf("更新失败, video: %#v, err: %w", video, err)
	}
	return
}

// FetchDouBanVideoByType 通过类型获取豆瓣视频
//
//	@Description:
//	@receiver m
//	@param typ
//	@return nameList
//	@return err
func (m *MovieDB) FetchDouBanVideoByType(typ types.VideoType) (nameList map[*types.DouBanVideo][]string, err error) {
	nameList = make(map[*types.DouBanVideo][]string)

	var videos []*types.DouBanVideo
	if err = m.db.Model(&types.DouBanVideo{}).Where("type = ?", typ.String()).Find(&videos).Error; err != nil {
		return nil, err
	}
	for _, video := range videos {
		var names []string
		if err = json.Unmarshal([]byte(video.Names), &names); err != nil {
			log.Error(err)
			continue
		}
		nameList[video] = names
	}
	return
}

// FetchThisYearVideo
//
//	@Description: 获取今年的视频
//	@receiver m
//	@return []types.DouBanVideo
//	@return error
func (m *MovieDB) FetchThisYearVideo() ([]*types.DouBanVideo, error) {
	thisYear := time.Now().Format("2006")
	var videos []*types.DouBanVideo
	//nolint:perfsprint
	err := m.db.Model(&types.DouBanVideo{}).Where("date_published like  ?", fmt.Sprintf("%s%%", thisYear)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
