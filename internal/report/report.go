package report

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/ipProxy"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"os"
)

type Report struct {
	scheduling string
}

//
// NewReport
//  @Description: 新建Report
//  @param scheduling
//  @return *Report
//
func NewReport(scheduling string) *Report {
	return &Report{scheduling: scheduling}
}

//
// Run
//  @Description: 运行
//  @receiver r
//
func (r *Report) Run() {
	if r.scheduling == "" {
		log.Error("Report: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("Report: Scheduling is: [%s]", r.scheduling)
	c := cron.New()
	c.AddFunc(r.scheduling, func() {

		// FeedVideo 资源统计
		reportFeedVideoStatistics()
		// aria下载列表统计
		reportAria2TaskStatistics()
		// 代理统计
		reportIpProxyStatistics()
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
		// todo 下载完后的向TG通知
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
		tableData = append(tableData, []string{reportCount.Web, fmt.Sprintf("%d", reportCount.Count)})
	}
	for _, v := range tableData {
		table.Append(v)
	}
	table.SetFooter([]string{"总数", fmt.Sprintf("%d", Total)})
	log.Info("\n\n下载统计: ")
	table.Render()

}

func reportIpProxyStatistics() {

	c := ipProxy.FetchProxyTypeCount()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Http", "Https", "tcp", "Other"})
	tableData := [][]string{}

	tableData = append(tableData, []string{fmt.Sprintf("%d", c.Http), fmt.Sprintf("%d", c.Https), fmt.Sprintf("%d", c.Tcp), fmt.Sprintf("%d", c.Other)})
	for _, v := range tableData {
		table.Append(v)
	}

	log.Info("\n\n代理统计: ")
	table.Render()
}
