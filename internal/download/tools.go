package download

import (
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"strings"
)

//
// filterByResolution
//  @Description: 根据 清晰度 过滤
//  @param videos
//  @return list
//
func filterByResolution(videos ...*types.FeedVideo) (list []*types.FeedVideo) {
	// 1. 在下载历史表中查看是否有此视频的下载记录
	inDownloadHistory := filterByResolutionInDownloadHistory(videos...)

	for i, video := range inDownloadHistory {
		log.Info(i, "  ", video.TorrentName, "  ", video.Name, "  ", video.ID)
	}

	// 2. 如果有多个视频 还需要在这一次的视频中过滤出清晰度最高的2个
	needDownloadFeedVideo, _ := filterVideosByResolution(inDownloadHistory...)
	return needDownloadFeedVideo
}

//
// filterByResolutionInDownloadHistory
//  @Description: 根据 清晰度 在下载历史表中  过滤
//  @param videos
//  @return list
//
func filterByResolutionInDownloadHistory(videos ...*types.FeedVideo) (list []*types.FeedVideo) {
	for _, video := range videos {
		// 通过清晰度过滤
		v, err := model.NewMovieDB().FindFeedVideoInDownloadHistory(video)
		if err != nil {
			//log.Warn(err)
			continue
		}
		list = append(list, v)
		log.Errorf("%s  %s   %d   ", v.TorrentName, v.Name, v.ID)
	}

	return
}

//
// UpdateFeedVideoAndDownloadHistory
//  @Description: 更新 feedVideo 的下载状态为1 记录这一次下载的视频
//  @param video
//
func UpdateFeedVideoAndDownloadHistory(video *types.FeedVideo) {
	//  更新 feedVideo 的下载状态为1
	if err := model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 1); err != nil {
		log.Error(err)
	}
	//log.Errorf("%s       %s       %d", video.TorrentName, video.Name, video.Convert2DownloadHistory().Resolution)
	//  记录这一次下载的视频
	if err := model.NewMovieDB().UpdateOrAddDownloadHistory(video.Convert2DownloadHistory()); err != nil {
		log.Error(video.TorrentName, video.Name, err)
	}
}

//
// filterVideosByResolution
//  @Description: 根据分辨率排序
//  @receiver d
//  @param videos
//  @return needDownloadFeedVideo  需要下载的资源
//  @return downloadIs3 需要记录的资源
//
func filterVideosByResolution(videos ...*types.FeedVideo) (needDownloadFeedVideo []*types.FeedVideo, needRecordFeedVideo []*types.FeedVideo) {
	//  将视频按照分辨率排序，当前仅排序 2160P 和 1080P
	Videos2160P, Videos1080P := sotByResolution(videos)

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

	} else {
		if len(Videos1080P) >= 2 {
			needDownloadFeedVideo = append(needDownloadFeedVideo, Videos1080P[0:2]...)
			needRecordFeedVideo = append(needRecordFeedVideo, Videos1080P[2:]...)
		} else {
			needDownloadFeedVideo = append(needDownloadFeedVideo, Videos1080P...)
		}
	}
	return
}

//
// sotByResolution
//  @Description:  根据分辨率排序
//  @receiver d
//  @param videos
//  @return Videos2160P
//  @return Videos1080P
//
func sotByResolution(videos []*types.FeedVideo) (Videos2160P []*types.FeedVideo, Videos1080P []*types.FeedVideo) {
	if len(videos) < 1 {
		return
	}

	for _, v := range videos {
		switch {
		// 如果是2060p 放到 Videos2160P 列表
		case strings.Contains(v.TorrentName, "2160"):
			Videos2160P = append(Videos2160P, v)
		// 如果是1080p 放到 Videos1080P 列表
		case strings.Contains(v.TorrentName, "1080"):
			Videos1080P = append(Videos1080P, v)
		}
	}
	return Videos2160P, Videos1080P
}
