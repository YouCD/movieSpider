package report

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"os"
)

type Report struct {
	scheduling string
}

func NewReport(scheduling string) *Report {
	return &Report{scheduling: scheduling}
}

func (r *Report) Run() {
	if r.scheduling == "" {
		log.Error("Report: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("Report: Scheduling is: [%s]", r.scheduling)
	c := cron.New()
	c.AddFunc(r.scheduling, func() {

		// 资源统计
		count, err := model.NewMovieDB().CountFeedVideo()
		if err != nil {
			log.Error("Report: err", err)
		}
		var s string
		var Total int
		for _, reportCount := range count {
			Total += reportCount.Count
			s += fmt.Sprintf("%s: %d ", reportCount.Web, reportCount.Count)
		}
		log.Infof("Report: feed_video 数据统计: Total: %d  %s", Total, s)

		// 下载情况统计
		aria2Client, err := aria2.NewAria2(config.Downloader.Aria2Label)
		if err != nil {
			log.Error("Report: err", err)
			return
		}
		files := aria2Client.CompletedFiles()
		var msg string
		for _, file := range files {
			s += fmt.Sprintf("\nGID:%s, 大小:%s, 已完成:%s, 文件名:%s", file.GID, file.Size, file.Completed, file.FileName)
			// todo 下载完后的向TG通知
		}
		log.Infof("Report: 下载统计: %s", msg)

	})
	c.Start()
}
