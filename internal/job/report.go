package job

import (
	"context"
	"os"
	"strconv"

	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/olekukonko/tablewriter"
	"github.com/robfig/cron/v3"
	"github.com/youcd/toolkit/log"
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
func (r *Report) Run(ctx context.Context) {
	if r.scheduling == "" {
		log.WithCtx(context.Background()).Error("Report: Scheduling is null")
		return
	}
	log.WithCtx(context.Background()).Infof("Report: Scheduling is: [%s]", r.scheduling)
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
	<-ctx.Done()
	c.Stop()
	log.WithCtx(context.Background()).Info("Report: Stopped")
}

func reportAria2TaskStatistics() {
	// 下载情况统计
	aria2Client, err := aria2.NewAria2(config.Config().Downloader.Aria2Label)
	if err != nil {
		log.WithCtx(context.Background()).Error("Report: err", err)
		return
	}
	files := aria2Client.CurrentActiveAndStopFiles()
	if len(files) == 0 {
		return
	}

	downloadTable := tablewriter.NewWriter(os.Stdout)
	downloadTable.Header([]string{"GID", "大小", "已完成", "文件名"})
	for _, file := range files {
		downloadTable.Append([]string{file.GID, file.Size, file.Completed, file.FileName})
	}

	log.WithCtx(context.Background()).Info("\n\n当前下载信息: ")
	downloadTable.Render()
}

func reportFeedVideoStatistics() {
	count, err := model.NewMovieDB().CountFeedVideo()
	if err != nil {
		log.WithCtx(context.Background()).Error("Report: err", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Web", "Count"})
	var Total int
	for _, reportCount := range count {
		Total += reportCount.Count
		table.Append([]string{reportCount.Web, strconv.Itoa(reportCount.Count)})
	}
	if Total == 0 {
		log.WithCtx(context.Background()).Info("无下载资源")
		return
	}
	table.Footer([]string{"总数", strconv.Itoa(Total)})
	log.WithCtx(context.Background()).Info("\n\n下载统计: ")
	table.Render()

	count, err = model.NewMovieDB().CountFeedVideoByToday()
	if err != nil {
		log.WithCtx(context.Background()).Error("Report: err", err)
		return
	}

	slice.SortBy(count, func(a, b *types.ReportCount) bool {
		if a.Count < b.Count {
			return true
		}
		return false
	})

	table = tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Web", "Count"})
	// 今日统计使用独立变量, 避免累加到总库量 Total 上
	var todayTotal int
	for _, reportCount := range count {
		todayTotal += reportCount.Count
		table.Append([]string{reportCount.Web, strconv.Itoa(reportCount.Count)})
	}
	if todayTotal == 0 {
		log.WithCtx(context.Background()).Info("今日无新增资源")
		return
	}
	table.Footer([]string{"总数", strconv.Itoa(todayTotal)})
	log.WithCtx(context.Background()).Info("\n\n今日下载统计: ")
	table.Render()
}

// aria2 下载队列
func reportAria2DownloadQueue() {
	newAria2, err := aria2.NewAria2(config.Config().Downloader.Aria2Label)
	if err != nil {
		log.WithCtx(context.Background()).Error("Report: err", err)
		return
	}
	task := newAria2.GetDownloadTask()
	if len(task) == 0 {
		log.WithCtx(context.Background()).Info("Report: aria2 下载队列为空")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Gid", "Names", "Type"})
	for k, v := range task {
		table.Append([]string{k, v.Name, v.Type})
	}
	log.WithCtx(context.Background()).Info("\n\n下载队列: ")
	table.Render()
}
