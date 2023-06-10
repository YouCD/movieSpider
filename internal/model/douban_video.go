package model

import (
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"math/rand"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"time"
)

//
// CreatDouBanVideo
//  @Description: 创建豆瓣视频
//  @receiver m
//  @param video
//  @return err
//
func (m *movieDB) CreatDouBanVideo(video *types.DouBanVideo) (err error) {
	if video == nil {
		return errors.New("video 不能为nil")
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
		if v.Playable != video.Playable {
			v.Playable = video.Playable
			log.Debugf("FetchOneDouBanVideoByDouBanID %#v", v)
			err = m.UpdateDouBanVideo(video)
			return errors.WithMessagef(err, "UpdateDouBanVideo %s", v.Names)
		}
		return nil
	}

	if video.Names == "null" {
		log.Errorf("CreatDouBanVideo 数据错误. video: %#v", video)
		return
	}

	err = m.db.Model(&types.DouBanVideo{}).Create(video).Error

	if err != nil {
		return errors.WithMessage(err, video.Names)
	}
	log.Debugf("CreatDouBanVideo 数据已添加. video: %#v", video)
	return
}

//
// RandomOneDouBanVideo
//  @Description: 随机获取一个豆瓣视频
//  @receiver m
//  @return video
//  @return err
//
func (m *movieDB) RandomOneDouBanVideo() (video *types.DouBanVideo, err error) {
	video = new(types.DouBanVideo)
	rows, err := m.db.Model(&types.DouBanVideo{}).Select(" id,names,douban_id,playable").Where("imdb_id = ''").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []*types.DouBanVideo
	for rows.Next() {
		var v types.DouBanVideo
		err = rows.Scan(&v.ID, &v.Names, &v.DoubanID, &v.Playable)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &v)
	}
	if len(videos) == 0 {
		return nil, errors.New("RandomOneDouBanVideo data is null")
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(videos))
	video = videos[index]
	log.Debugf("RandomOneDouBanVideo video: %#v", video)
	return
}

//
// FetchOneDouBanVideoByDouBanID
//  @Description: 根据豆瓣ID获取豆瓣视频
//  @receiver m
//  @param DouBanID
//  @return video
//  @return err
//
func (m *movieDB) FetchOneDouBanVideoByDouBanID(DouBanID string) (video *types.DouBanVideo, err error) {
	err = m.db.Model(&types.DouBanVideo{}).Where("douban_id=?", DouBanID).Scan(&video).Error
	if err != nil {
		return nil, err
	}
	log.Debugf("FetchOneDouBanVideoByDouBanID video: %#v", video)
	return
}

//
// UpdateDouBanVideo
//  @Description: 更新豆瓣视频
//  @receiver m
//  @param video
//  @return err
//
func (m *movieDB) UpdateDouBanVideo(video *types.DouBanVideo) (err error) {
	if video == nil {
		return errors.New("空数据")
	}

	err = m.db.Model(&types.DouBanVideo{}).Where("douban_id = ?", video.DoubanID).Updates(&types.DouBanVideo{
		Names:     video.Names,
		ImdbID:    video.ImdbID,
		RowData:   video.RowData,
		Timestamp: time.Now().Unix(),
		Type:      video.Type,
		Playable:  video.Playable,
	}).Error
	if err != nil {
		return errors.WithMessage(err, video.Names)
	}
	return
}

//
// FetchDouBanVideoByType 通过类型获取豆瓣视频
//  @Description:
//  @receiver m
//  @param typ
//  @return nameList
//  @return err
//
func (m *movieDB) FetchDouBanVideoByType(typ types.Resource) (nameList map[types.DouBanVideo][]string, err error) {
	nameList = make(map[types.DouBanVideo][]string)
	rows, err := m.db.Model(&types.DouBanVideo{}).Select("douban_id,names").Where("type = ?", typ.Typ()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tv types.DouBanVideo

		if err = rows.Scan(&tv.DoubanID, &tv.Names); err != nil {
			continue
		}
		var names []string
		if err = json.Unmarshal([]byte(tv.Names), &names); err != nil {
			log.Error(err)
			continue
		}
		nameList[tv] = names

	}

	return
}
