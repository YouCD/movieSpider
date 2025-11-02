package model

import (
	"errors"
	"movieSpider/internal/types"
	"strings"

	"github.com/youcd/toolkit/log"
)

// ProxySaveVideo2DB
//
//	@Description: 代理保存视频到数据库
//	@param videos
//

func ProxySaveVideo2DB(videos ...*types.FeedVideo) {
	if len(videos) == 0 {
		log.Warn("没有数据")
		return
	}

	for _, v := range videos {
		go func(video *types.FeedVideo) {
			err := NewMovieDB().CreatFeedVideo(video)
			if err != nil {
				if errors.Is(err, ErrDataExist) {
					log.Debugf("%s.%s err: %s", strings.ToUpper(video.Web), video.Type, err)
					return
				}
				log.Error(err)
				return
			}
			log.Infof("%s.%s: %s 保存完毕.", strings.ToUpper(video.Web), video.Type, video.Name)
		}(v)
	}
}
