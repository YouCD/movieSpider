package download

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"os"
)

type Download struct {
	scheduling string
	types.Resolution
}

func NewDownloader(scheduling string) *Download {
	return &Download{scheduling: scheduling}
}

func (d *Download) downloadTask() {
	err := d.downloadMovieTask()
	if err != nil {
		log.Error(err)
	}
	err = d.downloadTvTask()
	if err != nil {
		log.Error(err)
	}
}

//
// downloadTvTask
//  @Description: 下载电视剧
//  @receiver d
//  @return err
//
func (d *Download) downloadTvTask() (err error) {
	log.Info("Downloader tv working...")
	tvs, err := model.NewMovieDB().FetchDouBanTvVideo()
	if err != nil {
		return err
	}

	//  FilterMap 暂存 电视剧名相同的视频
	var FilterMap = make(map[string][]*types.FeedVideo)
	log.Info("查找需要下载的tv.")

	//归类同一个电视剧名的 feedVideo
	for name, tvName := range tvs {
		// 获取 tv
		tvVideos, err := model.NewMovieDB().GetFeedVideoTVByName(tvName...)
		if err != nil {
			log.Warn(err)
		}
		if len(tvVideos) == 0 {
			log.Warnf("name: %s 已全部下载完毕，或该影片没有更新.", name)
			continue
		}
		//log.Errorf("%s             %d      %s", name, len(tvVideos), tvName)

		// 归类同一个电视剧名的视频
		for _, video := range tvVideos {
			//  将此次所有feedVideo的下载状态更新为3
			err = model.NewMovieDB().UpdateFeedVideoDownloadByID(video.ID, 3)
			if err != nil {
				log.Error(err)
			}

			downloadHistory := video.Convert2DownloadHistory()
			if downloadHistory == nil {
				log.Debugf("TorrentName: %#v 不能转化为 downloadHistory ", video.TorrentName)
				continue
			}

			FilterMap[name] = append(FilterMap[name], video)
		}
	}

	// 根据 清晰度 季数和集数过滤
	needDownloadFeedVideo := make([]*types.FeedVideo, 0)
	for _, v := range FilterMap {
		list := filterByResolution(v...)
		needDownloadFeedVideo = append(needDownloadFeedVideo, list...)
	}
	//  如果没有需要下载的视频 则返回
	if len(needDownloadFeedVideo) == 0 {
		log.Warn("此次没有要下载的tv.")
		return
	}

	//推送 磁力连接至 aria2
	err = d.aria2Download(needDownloadFeedVideo...)
	if err != nil {
		log.Error(err)
	}

	// 更新feedVideo的下载状态，记录这一次下载的视频
	for _, video := range needDownloadFeedVideo {
		UpdateFeedVideoAndDownloadHistory(video)
	}

	return
}

//
// downloadMovieTask
//  @Description: 下载电影
//  @receiver d
//  @return error
//
func (d *Download) downloadMovieTask() error {
	// 获取 豆瓣 数据
	log.Info("Downloader movie working...")
	names, err := model.NewMovieDB().FetchDouBanVideoByType(types.ResourceMovie)
	if err != nil {
		return err
	}

	// 获取 feedVideo movie
	log.Info("查找需要下载的movie.")
	MovieVideos, err := model.NewMovieDB().GetFeedVideoMovieByName(names...)
	if err != nil {
		return err
	}

	//  将此次所有feedVideo movie的下载状态更新为3
	for _, v := range MovieVideos {
		if err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 3); err != nil {
			log.Error(err)
			continue
		}
	}

	//log.Error("MovieVideos: ", len(MovieVideos))
	// 通过清晰度过滤已经下载过的视频
	needDownloadMovieList := filterByResolution(MovieVideos...)

	//  如果没有需要下载的视频 则返回
	if len(needDownloadMovieList) == 0 {
		log.Warn("此次没有要下载的movie.")
		return nil
	}

	// 推送 磁力连接至 aria2
	err = d.aria2Download(MovieVideos...)
	if err != nil {
		log.Warn(err)
	}

	for _, video := range needDownloadMovieList {
		UpdateFeedVideoAndDownloadHistory(video)
	}

	log.Error("needDownloadMovieList: ", len(needDownloadMovieList))

	return err
}

//
// aria2Download
//  @Description: 通过aria2下载
//  @receiver d
//  @param videos
//  @return err
//
func (d *Download) aria2Download(videos ...*types.FeedVideo) (err error) {
	if len(videos) < 1 {
		return errors.New("没有需要下载的视频")
	}
	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		return errors.WithMessage(err, "aria2 初始化失败")
	}
	for _, v := range videos {
		gid, err := newAria2.DownloadByUrl(v.Magnet)
		if err != nil {
			return err
		}

		// 如果开启了tg推送 则推送
		if config.TG.Enable {
			go func() {
				bus.NotifyChan <- fmt.Sprintf("%s 开始下载. GID: %s", v.TorrentName, gid)
			}()
		}

		log.Infof("Downloader: %s 开始下载. GID: %s", v.TorrentName, gid)
	}
	return nil
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

//func (d *Download) DownloadByName(name, Resolution string) (msg string) {
//	// todo 从 knaben 搜索 todo
//	feedKnaben := feedSpider.NewFeedKnaben(name, d.ResolutionStr2Int(Resolution))
//	_, err := feedKnaben.Crawler()
//	if err != nil {
//		log.Error(err)
//	}
//	//todo 从 Bt4g 搜索
//	feedBt4g := feedSpider.NewFeedBt4g(name, d.ResolutionStr2Int(Resolution))
//	_, err = feedBt4g.Crawler()
//	if err != nil {
//		log.Error(err)
//	}
//
//	// 获取 磁力连接
//	videos, err := model.NewMovieDB().GetFeedVideoMovieByName([]string{name}...)
//	if err != nil {
//		log.Error(err)
//	}
//
//	if len(videos) == 0 {
//		return fmt.Sprint("所有资源已下载过,或没有可下载资源.")
//	}
//
//	// 推送 磁力连接至 aria2
//	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
//	if err != nil {
//		log.Error(err)
//	}
//	for _, v := range videos {
//		gid, err := newAria2.DownloadByUrl(v.Magnet)
//		if err != nil {
//			log.Error(err)
//			return
//		}
//		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 1)
//		if err != nil {
//			log.Error(err)
//		}
//		log.Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
//	}
//
//	return fmt.Sprintf("已将 %d 资源加入下载.", len(videos))
//}
