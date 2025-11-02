package download

import (
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

//nolint:gochecknoglobals
var wg sync.WaitGroup

type Download struct {
	scheduling string
	types.Resolution
}

func NewDownloader(scheduling string) *Download {
	return &Download{scheduling: scheduling}
}
func (d *Download) DownloadByName(name, resolution string) (msg string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		feedBt4g := searchspider.NewFeedBt4g(name, d.ResolutionStr2Int(resolution))
		_, err := feedBt4g.Search()
		if err != nil {
			log.Error(err)
		}
	}()
	wg.Wait()

	// 获取 磁力连接
	videos, err := model.NewMovieDB().GetFeedVideoMovieByNames([]string{name}...)
	if err != nil {
		log.Error(err)
	}

	if len(videos) == 0 {
		return "所有资源已下载过,或没有可下载资源."
	}

	// 推送 磁力连接至 aria2
	newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		log.Error(err)
	}

	for _, v := range videos {
		if v.Name == "" {
			log.Warnf("TorrentName: %v ,name is nil", v.TorrentName)
			continue
		}
		gid, err := newAria2.DownloadByWithVideo(v, v.Magnet)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 1)
		if err != nil {
			log.Error(err)
		}
	}

	return fmt.Sprintf("已将 %d 资源加入下载.", len(videos))
}
func (d *Download) Run() {
	if d.scheduling == "" {
		log.Error("Downloader: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("Downloader: Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	_, err := c.AddFunc(d.scheduling, func() {
		d.downloadTask()
	})
	if err != nil {
		log.Error("Downloader: AddFunc is null")
		os.Exit(1)
	}
	c.Start()
}

func (d *Download) downloadTask() {
	err := d.download(types.VideoTypeTV, model.NewMovieDB().GetFeedVideoTVByNames)
	if err != nil {
		log.Error(err)
	}
	err = d.download(types.VideoTypeMovie, model.NewMovieDB().GetFeedVideoMovieByNames)
	if err != nil {
		log.Error(err)
	}
}

func (d *Download) download(tvOrMovie types.VideoType, f func(names ...string) ([]*types.FeedVideo, error)) (err error) {
	log.Infow(tvOrMovie.String(), "Downloader working...")
	videos, err := model.NewMovieDB().FetchDouBanVideoByType(tvOrMovie)
	if err != nil {
		return fmt.Errorf("FetchDouBanVideoByType,err: %w", err)
	}

	//  FilterMap 暂存 电视剧名相同的视频
	var FilterMap = make(map[string][]*types.FeedVideo)

	var videoList []*types.FeedVideo
	// 归类同一个电视剧名的 feedVideo
	for douBanVideo, name := range videos {
		log.Infow(tvOrMovie.String(), "douBanVideo", douBanVideo.Names)
		videoList, err = f(name...)
		if err != nil {
			log.Warn(err)
		}
		if len(videoList) == 0 {
			continue
		}
		log.Infof("douBanVideo:%v   种子数: %#v", douBanVideo.Names, len(videoList))
		// 归类同一个电视剧名的视频
		for _, video := range videoList {
			// 添加 豆瓣ID
			video.DoubanID = douBanVideo.DoubanID
			//  将此次所有feedVideo的下载状态更新为3
			video.Download = 3
			// 如果 feedVideo 不能转化为 downloadHistory 则跳过
			downloadHistory := video.Convert2DownloadHistory()
			if downloadHistory == nil {
				log.Debugf("TorrentName: %#v 不能转化为 downloadHistory ", video.TorrentName)
				continue
			}
			FilterMap[douBanVideo.Names] = append(FilterMap[douBanVideo.Names], video)
		}
	}
	// 批量更新
	err = model.NewMovieDB().UpdateFeedVideos(videoList...)
	if err != nil {
		log.Error(err)
	}

	// 根据 清晰度 季数和集数过滤
	needDownloadFeedVideo := make([]*types.FeedVideo, 0)
	for _, v := range FilterMap {
		list := FilterByResolution(types.VideoTypeTV, v...)
		needDownloadFeedVideo = append(needDownloadFeedVideo, list...)
	}

	//  如果没有需要下载的视频 则返回
	if len(needDownloadFeedVideo) == 0 {
		log.Warn("此次没有要下载的tv.")
		return nil
	}

	// 推送 磁力连接至 aria2
	err = d.aria2Download(needDownloadFeedVideo...)
	if err != nil {
		log.Error(err)
	}

	// 更新feedVideo的下载状态，记录这一次下载的视频
	for _, video := range needDownloadFeedVideo {
		log.Infow(tvOrMovie.String(), "更新", video.TorrentName)
		UpdateFeedVideoAndDownloadHistory(video)
	}
	return nil
}

var (
	ErrVideoIsNil = errors.New("video is nil")
)

// aria2Download
//
//	@Description: 通过aria2下载
//	@receiver d
//	@param videos
//	@return err
func (d *Download) aria2Download(videos ...*types.FeedVideo) error {
	if len(videos) < 1 {
		return ErrVideoIsNil
	}
	newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		return fmt.Errorf("aria2 初始化失败,err: %w", err)
	}
	for _, v := range videos {
		if v.Name == "" {
			log.Warnf("TorrentName: %v ,name is nil", v.TorrentName)
			continue
		}
		gid, err := newAria2.DownloadByWithVideo(v, v.Magnet)
		if err != nil {
			log.Error(err)
			continue
		}

		// 如果开启了tg推送 则推送
		if config.Config.TG != nil {
			go func() {
				bus.DownloadNotifyChan <- &types.DownloadNotifyVideo{
					FeedVideo: v,
					File:      v.TorrentName,
					Gid:       gid,
				}
			}()
		}

		log.Infof(" 开始下载: %s. videoType: %s.  GID: %s.", v.TorrentName, v.Type, gid)
	}
	return nil
}
