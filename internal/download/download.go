package download

import (
	"context"
	"errors"
	"fmt"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/searchspider"
	"movieSpider/internal/types"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
)

type Download struct {
	scheduling string
	types.Resolution
}

func NewDownloader(scheduling string) *Download {
	return &Download{scheduling: scheduling}
}

func (d *Download) DownloadByName(ctx context.Context, name, resolution string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		feedBt4g := searchspider.NewFeedBt4g(name, d.ResolutionStr2Int(resolution))
		_, err := feedBt4g.Search()
		if err != nil {
			log.WithCtx(ctx).Error(err)
		}
	}()
	wg.Wait()

	// 获取磁力连接
	videos, err := model.NewMovieDB().GetFeedVideoMovieByNames([]string{name}...)
	if err != nil {
		log.WithCtx(ctx).Error(err)
	}

	if len(videos) == 0 {
		return "所有资源已下载过,或没有可下载资源."
	}

	// 推送磁力连接至 aria2
	newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		log.WithCtx(ctx).Error(err)
		return "aria2 初始化失败"
	}

	downloadedCount := 0
	for _, v := range videos {
		if v.Name == "" {
			log.WithCtx(ctx).Warnf("TorrentName: %v ,name is nil", v.TorrentName)
			continue
		}
		gid, err := newAria2.DownloadByWithVideo(v, v.Magnet)
		if err != nil {
			log.WithCtx(ctx).Error(err)
			continue
		}
		log.WithCtx(ctx).Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
		if err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 1); err != nil {
			log.WithCtx(ctx).Error(err)
		}
		downloadedCount++
	}

	return fmt.Sprintf("已将 %d 资源加入下载.", downloadedCount)
}

func (d *Download) Run(ctx context.Context) {
	if d.scheduling == "" {
		log.WithCtx(ctx).Error("Downloader: Scheduling is null")
		os.Exit(1)
	}
	log.WithCtx(ctx).Infof("Downloader: Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	_, err := c.AddFunc(d.scheduling, func() {
		d.downloadTask(ctx)
	})
	if err != nil {
		log.WithCtx(ctx).Error("Downloader: AddFunc is null")
		os.Exit(1)
	}
	c.Start()
}

func (d *Download) downloadTask(ctx context.Context) {
	err := d.download(ctx, types.VideoTypeTV, model.NewMovieDB().GetFeedVideoTVByNames)
	if err != nil {
		log.WithCtx(ctx).Error(err)
	}
	err = d.download(ctx, types.VideoTypeMovie, model.NewMovieDB().GetFeedVideoMovieByNames)
	if err != nil {
		log.WithCtx(ctx).Error(err)
	}
}

// getVideosFunc 定义获取视频的函数类型
type getVideosFunc func(names ...string) ([]*types.FeedVideo, error)

func (d *Download) download(ctx context.Context, tvOrMovie types.VideoType, f getVideosFunc) error {
	log.WithCtx(ctx).Infof("%s Downloader working...", tvOrMovie.String())
	videos, err := model.NewMovieDB().FetchDouBanVideoByType(tvOrMovie)
	if err != nil {
		return fmt.Errorf("FetchDouBanVideoByType,err: %w", err)
	}

	// FilterMap 暂存电视剧名相同的视频
	filterMap := make(map[string][]*types.FeedVideo)

	var videoList []*types.FeedVideo
	// 归类同一个电视剧名的 feedVideo
	for douBanVideo, name := range videos {
		log.WithCtx(ctx).Infow(tvOrMovie.String(), "douBanVideo", douBanVideo.Names)
		videoList, err = f(name...)
		if err != nil {
			log.WithCtx(ctx).Warn(err)
		}
		if len(videoList) == 0 {
			continue
		}
		log.WithCtx(ctx).Infof("douBanVideo:%v   种子数: %#v", douBanVideo.Names, len(videoList))
		// 归类同一个电视剧名的视频
		for _, video := range videoList {
			// 添加豆瓣ID
			video.DoubanID = douBanVideo.DoubanID
			// 将此次所有feedVideo的下载状态更新为3
			video.Download = 3
			// 如果 feedVideo 不能转化为 downloadHistory 则跳过
			downloadHistory := video.Convert2DownloadHistory()
			if downloadHistory == nil {
				log.WithCtx(ctx).Debugf("TorrentName: %#v 不能转化为 downloadHistory ", video.TorrentName)
				continue
			}
			filterMap[douBanVideo.Names] = append(filterMap[douBanVideo.Names], video)
		}
	}
	// 批量更新
	if err = model.NewMovieDB().UpdateFeedVideos(videoList...); err != nil {
		log.WithCtx(ctx).Error(err)
	}

	// 根据清晰度、季数和集数过滤
	needDownloadFeedVideo := make([]*types.FeedVideo, 0)
	for _, v := range filterMap {
		list := FilterByResolution(types.VideoTypeTV, v...)
		needDownloadFeedVideo = append(needDownloadFeedVideo, list...)
	}

	// 如果没有需要下载的视频则返回
	if len(needDownloadFeedVideo) == 0 {
		log.WithCtx(ctx).Warn("此次没有要下载的tv.")
		return nil
	}

	// 推送磁力连接至 aria2
	if err = d.aria2Download(ctx, needDownloadFeedVideo...); err != nil {
		log.WithCtx(ctx).Error(err)
	}

	// 更新feedVideo的下载状态，记录这一次下载的视频
	for _, video := range needDownloadFeedVideo {
		log.WithCtx(ctx).Infow(tvOrMovie.String(), "更新", video.TorrentName)
		UpdateFeedVideoAndDownloadHistory(video)
	}
	return nil
}

var ErrVideoIsNil = errors.New("video is nil")

// aria2Download 通过aria2下载
func (d *Download) aria2Download(ctx context.Context, videos ...*types.FeedVideo) error {
	if len(videos) < 1 {
		return ErrVideoIsNil
	}
	newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		return fmt.Errorf("aria2 初始化失败,err: %w", err)
	}
	for _, v := range videos {
		if v.Name == "" {
			log.WithCtx(ctx).Warnf("TorrentName: %v ,name is nil", v.TorrentName)
			continue
		}
		gid, err := newAria2.DownloadByWithVideo(v, v.Magnet)
		if err != nil {
			log.WithCtx(ctx).Error(err)
			continue
		}

		// 如果开启了tg推送则推送
		if config.Config.TG != nil {
			go func(video *types.FeedVideo, gid string) {
				bus.DownloadNotifyChan <- &types.DownloadNotifyVideo{
					FeedVideo: video,
					File:      video.TorrentName,
					Gid:       gid,
				}
			}(v, gid)
		}

		log.WithCtx(ctx).Infof(" 开始下载: %s. videoType: %s.  GID: %s.", v.TorrentName, v.Type, gid)
	}
	return nil
}
