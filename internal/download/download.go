package download

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	feed2 "movieSpider/internal/feed"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	types2 "movieSpider/internal/types"
	"os"
	"strings"
)

type Download struct {
	scheduling string
	types2.Resolution
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

func (d *Download) downloadTvTask() (err error) {
	log.Info("Downloader tv working...")
	tvs, err := model.NewMovieDB().FetchDouBanVideoByType(types2.ResourceTV)
	if err != nil {
		return err
	}

	// 获取 磁力连接
	tvVides, err := model.NewMovieDB().FetchTVMagnetByName(tvs)
	if err != nil {
		return err
	}
	if len(tvVides) == 0 {
		log.Warn("Downloader: 此次没有查询到要下载的资源.")
		return
	}
	is1, is3 := d.sotByResolution(tvVides)
	// 推送 磁力连接至 aria2
	err = d.aria2Download(is1)
	if len(tvVides) == 0 {
		return err
	}
	for _, v := range is3 {
		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 3)
		if len(tvVides) == 0 {
			continue
		}
	}

	return
}

func (d *Download) downloadMovieTask() error {
	// 获取 豆瓣 数据
	log.Info("Downloader movie working...")
	names, err := model.NewMovieDB().FetchDouBanVideoByType(types2.ResourceMovie)
	if err != nil {
		return err
	}

	// 获取 磁力连接
	MovieVides, err := model.NewMovieDB().FetchMovieMagnetByName(names)
	if err != nil {
		return err
	}
	// 推送 磁力连接至 aria2
	err = d.aria2Download(MovieVides)
	if err != nil {
		return err
	}
	for _, v := range MovieVides {
		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 3)
		if err != nil {
			return err
		}
	}

	return err
}

func (d *Download) aria2Download(vides []*types2.FeedVideo) (err error) {

	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		return errors.WithMessage(err, "aria2 下载错误")
	}
	for _, v := range vides {
		gid, err := newAria2.DownloadByUrl(v.Magnet)
		if err != nil {
			return err
		}
		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 1)
		if err != nil {
			return err
		}

		// 通知
		go func() {
			bus.NotifyChan <- fmt.Sprintf("%s 开始下载. GID: %s", v.TorrentName, gid)
		}()

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

func (d *Download) DownloadByName(name, Resolution string) (msg string) {
	// 从 knaben 搜索
	feedKnaben := feed2.NewFeedKnaben(name, d.ResolutionStr2Int(Resolution))
	_, err := feedKnaben.Crawler()
	if err != nil {
		log.Error(err)
	}
	// 从 Bt4g 搜索
	feedBt4g := feed2.NewFeedBt4g(name, d.ResolutionStr2Int(Resolution))
	_, err = feedBt4g.Crawler()
	if err != nil {
		log.Error(err)
	}
	// 获取 磁力连接
	vides, err := model.NewMovieDB().FetchMovieMagnetByName([]string{name})
	if err != nil {
		log.Error(err)
	}

	if len(vides) == 0 {
		return fmt.Sprint("所有资源已下载过,或没有可下载资源.")
	}

	// 推送 磁力连接至 aria2
	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		log.Error(err)
	}
	for _, v := range vides {
		gid, err := newAria2.DownloadByUrl(v.Magnet)
		if err != nil {
			log.Error(err)
			return
		}
		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 1)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
	}

	return fmt.Sprintf("已将 %d 资源加入下载.", len(vides))
}
func (d *Download) sotByResolution(videos []*types2.FeedVideo) (downloadIs1 []*types2.FeedVideo, downloadIs3 []*types2.FeedVideo) {
	var Videos2160P []*types2.FeedVideo
	var Videos1080P []*types2.FeedVideo
	for _, v := range videos {
		switch {
		// 如果是2060p 放到 Videos2160P 列表
		case strings.Contains(v.TorrentName, "2160"):
			Videos2160P = append(Videos2160P, v)
		// 如果是1080p 放到 Videos2160P 列表
		case strings.Contains(v.TorrentName, "1080"):
			Videos1080P = append(Videos1080P, v)
		// 其他的放到 downloadIs3 列表
		default:
			downloadIs3 = append(downloadIs3, v)
		}
	}

	// 如果 Videos2160P 有 数据
	if len(Videos2160P) > 0 {
		// 如果 Videos2160P 有大于2个片源
		if len(Videos2160P) >= 2 {
			// 前两个放到 downloadIs1 列表
			downloadIs1 = append(downloadIs1, Videos2160P[0:2]...)
			// 第3个好往后放到 downloadIs3 列表
			downloadIs3 = append(downloadIs3, Videos2160P[2:]...)
			// Videos1080P 放到 downloadIs1 列表
			downloadIs3 = append(downloadIs3, Videos1080P...)
		} else {
			// 如果 Videos2160P 少于2个片源
			downloadIs1 = append(downloadIs1, Videos2160P...)
			downloadIs3 = append(downloadIs3, Videos1080P...)
		}

	} else {
		if len(Videos1080P) >= 2 {
			downloadIs1 = append(downloadIs1, Videos1080P[0:2]...)
			downloadIs3 = append(downloadIs3, Videos1080P[2:]...)
		} else {
			downloadIs1 = append(downloadIs1, Videos1080P...)
		}
	}
	return
}
