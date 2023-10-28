package download

import (
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"strings"
)

// SotByResolution
//
//	@Description:  根据分辨率排序
//	@receiver d
//	@param videos
//	@return Videos2160P
//	@return Videos1080P
func SotByResolution(videos []*types.FeedVideo) (videos2160P []*types.FeedVideo, videos1080P []*types.FeedVideo) {
	if len(videos) < 1 {
		return
	}

	for _, v := range videos {
		switch {
		// 如果是2060p 放到 videos2160P 列表
		case strings.Contains(v.TorrentName, "2160"):
			videos2160P = append(videos2160P, v)
		// 如果是1080p 放到 videos1080P 列表
		case strings.Contains(v.TorrentName, "1080"):
			videos1080P = append(videos1080P, v)
		}
	}
	return videos2160P, videos1080P
}

// FilterByResolution
//
//	@Description: 根据 清晰度 过滤
//	@param videos
//	@return list
func FilterByResolution(movieOrTV types.VideoType, videos ...*types.FeedVideo) (list []*types.FeedVideo) {
	// 1. 在下载历史表中查看是否有此视频的下载记录
	inDownloadHistory := FilterByResolutionInDownloadHistory(videos...)

	// 2. 如果有多个视频 还需要在这一次的视频中过滤出清晰度最高的2个
	needDownloadFeedVideo, _ := FilterVideosByResolution(movieOrTV, inDownloadHistory...)
	return needDownloadFeedVideo
}

// FilterByResolutionInDownloadHistory
//
//	@Description: 根据 清晰度 在下载历史表中  过滤
//	@param videos
//	@return list
func FilterByResolutionInDownloadHistory(videos ...*types.FeedVideo) (list []*types.FeedVideo) {
	for _, video := range videos {
		// 通过清晰度过滤
		v, err := model.NewMovieDB().FindFeedVideoInDownloadHistory(video)
		if err != nil {
			log.Warn(err)
			continue
		}
		list = append(list, v)
	}

	return
}

// UpdateFeedVideoAndDownloadHistory
//
//	@Description: 更新 feedVideo 的下载状态为1 记录这一次下载的视频
//	@param video
func UpdateFeedVideoAndDownloadHistory(video *types.FeedVideo) {
	//  更新 feedVideo 的下载状态为1
	if err := model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 1); err != nil {
		log.Error(err)
	}

	//  记录这一次下载的视频
	if err := model.NewMovieDB().UpdateOrAddDownloadHistory(video.Convert2DownloadHistory()); err != nil {
		log.Error(video.TorrentName, video.Name, err)
	}
}

// FilterVideosByResolution
//
//	@Description: 根据分辨率排序
//	@receiver d
//	@param videos
//	@return needDownloadFeedVideo  需要下载的资源
//	@return downloadIs3 需要记录的资源
//
//nolint:nakedret
func FilterVideosByResolution(movieOrTV types.VideoType, videos ...*types.FeedVideo) (needDownloadFeedVideo []*types.FeedVideo, needRecordFeedVideo []*types.FeedVideo) {
	// 如果类型是电影
	if movieOrTV == types.VideoTypeMovie {
		return HandlerMovie(videos...)
	}

	// 如果是电视剧
	if movieOrTV == types.VideoTypeTV {
		// 1. 先根据分辨归类
		Videos2160P, Videos1080P := SotByResolution(videos)
		// 2. 创建一个 map 用来存放需要下载的 feedVideo  key的格式：Name + Season + Episode
		needDownloadFeedVideoMap := make(map[string][]*types.FeedVideo)
		// 3. 一个用来存放 2160P 的桶
		needDownloadFeedVideo2160PMap := make(map[string][]*types.FeedVideo)
		// 4. 一个用来存放 1080P 的桶
		needDownloadFeedVideo1080PMap := make(map[string][]*types.FeedVideo)
		// 5. 把 2160P 的视频放到 needDownloadFeedVideo2160PMap 中
		if len(Videos2160P) > 0 {
			needDownloadFeedVideo2160PMap = HandlerTv(Videos2160P...)
		}
		//  6. 把 1080P 的视频放到 needDownloadFeedVideo1080PMap 中
		if len(Videos1080P) > 0 {
			needDownloadFeedVideo1080PMap = HandlerTv(Videos1080P...)
		}
		// 7. 把 2160P 和 1080P 的视频放到相同的桶中
		for s, feedVideos := range needDownloadFeedVideo2160PMap {
			needDownloadFeedVideoMap[s] = append(needDownloadFeedVideoMap[s], feedVideos...)
		}
		for s, feedVideos := range needDownloadFeedVideo1080PMap {
			needDownloadFeedVideoMap[s] = append(needDownloadFeedVideoMap[s], feedVideos...)
		}

		// 8. 遍历 needDownloadFeedVideoMap
		for _, feedVideos := range needDownloadFeedVideoMap {
			// 9. 如果这一集tv 有多个视频
			if (len(feedVideos)) >= 2 {
				// 10. 利用 handlerTv 处理这一集tv
				need, Record := HandlerMovie(feedVideos...)
				needDownloadFeedVideo = append(needDownloadFeedVideo, need...)
				needRecordFeedVideo = append(needRecordFeedVideo, Record...)
			} else {
				needDownloadFeedVideo = append(needDownloadFeedVideo, feedVideos...)
			}
		}
	}

	return
}

// HandlerMovie
//
//	@Description: 处理电影类型的视频
//	@param videos
//	@return needDownloadFeedVideo
//	@return needRecordFeedVideo
func HandlerMovie(videos ...*types.FeedVideo) (needDownloadFeedVideo []*types.FeedVideo, needRecordFeedVideo []*types.FeedVideo) {
	Videos2160P, Videos1080P := SotByResolution(videos)
	// 如果 Videos2160P 有 数据
	if len(Videos2160P) > 0 {
		// 如果 Videos2160P 有大于2个片源
		if len(Videos2160P) >= 2 {
			// 前两个放到 needDownloadFeedVideo 列表
			needDownloadFeedVideo = append(needDownloadFeedVideo, Videos2160P[0:2]...)
			// 第3个往后放到 needRecordFeedVideo 列表
			needRecordFeedVideo = append(needRecordFeedVideo, Videos2160P[2:]...)
			// Videos1080P 放到 needDownloadFeedVideo 列表
			needRecordFeedVideo = append(needRecordFeedVideo, Videos1080P...)
		} else {
			// 如果 Videos2160P 少于2个片源
			needDownloadFeedVideo = append(needDownloadFeedVideo, Videos2160P...)
			needRecordFeedVideo = append(needRecordFeedVideo, Videos1080P...)
		}
	}
	if len(Videos1080P) >= 2 {
		needDownloadFeedVideo = append(needDownloadFeedVideo, Videos1080P[0:2]...)
		needRecordFeedVideo = append(needRecordFeedVideo, Videos1080P[2:]...)
	} else {
		needDownloadFeedVideo = append(needDownloadFeedVideo, Videos1080P...)
	}
	return
}

// HandlerTv
//
//	@Description: 处理电视剧类型的视频
//	@param videos
//	@return map[string][]*types.FeedVideo
func HandlerTv(videos ...*types.FeedVideo) map[string][]*types.FeedVideo {
	if len(videos) < 1 {
		return nil
	}
	needDownloadFeedVideoMap := make(map[string][]*types.FeedVideo)
	for _, video := range videos {
		historyObj := video.Convert2DownloadHistory()
		if historyObj == nil {
			continue
		}

		key := historyObj.Name + historyObj.Season + historyObj.Episode
		needDownloadFeedVideoMap[key] = append(needDownloadFeedVideoMap[key], video)
	}
	return needDownloadFeedVideoMap
}
