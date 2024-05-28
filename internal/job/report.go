package job

import (
	"github.com/olekukonko/tablewriter"
	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"os"
	"strconv"
)

type Report struct {
	scheduling string
}

// NewReport
//
//	@Description: 新建Report
//	@param scheduling
//	@return *Report
func NewReport(scheduling string) *Report {
	return &Report{scheduling: scheduling}
}

// Run
//
//	@Description: 运行
//	@receiver r
func (r *Report) Run() {
	if r.scheduling == "" {
		log.Error("Report: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("Report: Scheduling is: [%s]", r.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(r.scheduling, func() {
		// FeedVideo 资源统计
		reportFeedVideoStatistics()
		// aria下载列表统计
		reportAria2TaskStatistics()
		// aria2 下载队列
		reportAria2DownloadQueue()
	})
	c.Start()
}

func reportAria2TaskStatistics() {
	// 下载情况统计
	aria2Client, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		log.Error("Report: err", err)
		return
	}
	files := aria2Client.CurrentActiveAndStopFiles()
	if len(files) == 0 {
		return
	}

	downloadTable := tablewriter.NewWriter(os.Stdout)
	downloadTable.SetHeader([]string{"GID", "大小", "已完成", "文件名"})
	for _, file := range files {
		downloadTable.Append([]string{file.GID, file.Size, file.Completed, file.FileName})
	}

	log.Info("\n\n当前下载信息: ")
	downloadTable.Render()
}

func reportFeedVideoStatistics() {
	count, err := model.NewMovieDB().CountFeedVideo()
	if err != nil {
		log.Error("Report: err", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Web", "Count"})
	var Total int
	for _, reportCount := range count {
		Total += reportCount.Count
		table.Append([]string{reportCount.Web, strconv.Itoa(reportCount.Count)})
	}

	table.SetFooter([]string{"总数", strconv.Itoa(Total)})
	log.Info("\n\n下载统计: ")
	table.Render()
}

// aria2 下载队列
func reportAria2DownloadQueue() {
	newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		log.Error("Report: err", err)
		return
	}
	task := newAria2.GetDownloadTask()
	if len(task) == 0 {
		log.Info("Report: aria2 下载队列为空")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Gid", "Names", "Type"})
	for k, v := range task {
		table.Append([]string{k, v.Name, v.Type})
	}
	log.Info("\n\n下载队列: ")
	table.Render()
}
