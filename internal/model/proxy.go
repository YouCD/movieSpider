package model

import (
	"context"
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

func ProxySaveVideo2DB(ctx context.Context, videos ...*types.FeedVideo) {
	if len(videos) == 0 {
		log.WithCtx(ctx).Warn("没有数据")
		return
	}

	for _, v := range videos {
		go func(video *types.FeedVideo) {
			err := NewMovieDB().CreatFeedVideo(ctx, video)
			if err != nil {
				if errors.Is(err, ErrDataExist) {
					log.WithCtx(ctx).Debugf("%s.%s err: %s", strings.ToUpper(video.Web), video.Type, err)
					return
				}
				log.WithCtx(ctx).Error(err)
				return
			}
			log.WithCtx(ctx).Infof("%s.%s: %s 保存完毕.", strings.ToUpper(video.Web), video.Type, video.Name)
		}(v)
	}
}
