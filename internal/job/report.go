package job

import (
	"github.com/olekukonko/tablewriter"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
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
		//
		reportAria2DownloadRecordStatistics()
	})
	c.Start()
}

func reportAria2TaskStatistics() {
	// 下载情况统计
	aria2Client, err := aria2.NewAria2(config.Downloader.Aria2Label)
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
	downloadTableData := [][]string{}

	for _, file := range files {
		downloadTableData = append(downloadTableData, []string{file.GID, file.Size, file.Completed, file.FileName})
	}
	for _, v := range downloadTableData {
		downloadTable.Append(v)
	}

	log.Info("\n\n当前下载信息: ")
	downloadTable.Render()
}

func reportFeedVideoStatistics() {
	count, err := model.NewMovieDB().CountFeedVideo()
	if err != nil {
		log.Error("Report: err", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Web", "Count"})
	tableData := [][]string{}
	var Total int
	for _, reportCount := range count {
		Total += reportCount.Count
		tableData = append(tableData, []string{reportCount.Web, strconv.Itoa(reportCount.Count)})
	}
	for _, v := range tableData {
		table.Append(v)
	}
	table.SetFooter([]string{"总数", strconv.Itoa(Total)})
	log.Info("\n\n下载统计: ")
	table.Render()
}

// aria2 下载记录
func reportAria2DownloadRecordStatistics() {
	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		log.Error("Report: err", err)
		return
	}
	task := newAria2.GetDownloadTask()
	if len(task) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Gid", "Names", "Type"})
	tableData := [][]string{}

	for k, v := range task {
		tableData = append(tableData, []string{k, v.Names, v.Type})
	}

	log.Info("\n\n下载记录: ")
	table.Render()
}
