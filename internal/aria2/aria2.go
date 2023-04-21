package aria2

import (
	"context"
	"fmt"
	"github.com/zyxar/argo/rpc"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/types"
	"path"
	"strconv"
	"strings"
	"sync"
)

var (
	aria2Client *aria2
	once        sync.Once
)

type aria2 struct {
	aria2Client rpc.Client
}

//
// NewAria2
//  @Description: 初始化aria2
//  @param label
//  @return *aria2
//  @return error
//
func NewAria2(label string) (*aria2, error) {
	once.Do(func() {
		for _, v := range config.Aria2cList {
			if v.Label == label {
				client, err := rpc.New(context.TODO(), v.Url, v.Token, 0, nil)
				if err != nil {
					log.Error(err)
				}
				log.Debug(config.Aria2cList)
				aria2Client = &aria2{client}
			}
		}
	})
	return aria2Client, nil
}

//
// DownloadByUrl
//  @Description: 通过url下载
//  @receiver a
//  @param url
//  @return gid
//  @return err
//
func (a *aria2) DownloadByUrl(url string) (gid string, err error) {
	return a.aria2Client.AddURI([]string{url})
}

//
// DownloadList
//  @Description: 下载列表
//  @receiver a
//  @param url
//  @return gid
//  @return err
//
func (a *aria2) DownloadList(url string) (gid string, err error) {
	info, err := a.aria2Client.GetSessionInfo()
	if err != nil {
		return "", err
	}
	log.Infof("DownloadList: %s", info)
	return
}

//
// CurrentActiveAndStopFiles
//  @Description: 获取当前正在下载以及停止下载的文件
//  @receiver a
//  @return completedFiles
//
func (a *aria2) CurrentActiveAndStopFiles() (completedFiles []*types.ReportCompletedFiles) {
	// 获取已停止下载的文件
	sessionInfo, err := a.aria2Client.TellStopped(0, 100)
	if err != nil {
		log.Error(err)
		return nil
	}

	completedFiles1 := a.completedHandler(sessionInfo, completedFiles...)
	completedFiles = append(completedFiles, completedFiles1...)
	// 获取正在下载的文件
	ActiveSession, err := a.aria2Client.TellActive()
	if err != nil {
		log.Error(err)
		return nil
	}
	completedFiles2 := a.completedHandler(ActiveSession, completedFiles...)
	completedFiles = append(completedFiles, completedFiles2...)
	return
}

//
// completedHandler
//  @Description: 处理已完成的文件
//  @receiver a
//  @param sessionInfo
//  @param completedFiles
//  @return []*types.ReportCompletedFiles
//
func (a *aria2) completedHandler(sessionInfo []rpc.StatusInfo, completedFiles ...*types.ReportCompletedFiles) []*types.ReportCompletedFiles {
	for _, v := range sessionInfo {
		if len(v.Files) > 0 {
			//  过滤掉元数据文件
			if strings.Contains(v.Files[0].Path, "[METADATA]") {
				continue
			}

			// 下载了多少
			CompletedLength, err := strconv.Atoi(v.CompletedLength)
			if err != nil {
				log.Error(err)
			}
			// 文件大小
			Length, err := strconv.Atoi(v.TotalLength)
			if err != nil {
				log.Error(err)
			}

			//文件完成度百分比
			completed := float32(CompletedLength) / float32(Length) * 100

			_, file := path.Split(v.Files[0].Path)

			completedFiles = append(completedFiles, &types.ReportCompletedFiles{
				GID:       v.Gid,
				Size:      fmt.Sprintf("%.2fGB", float32(Length)/1024/1024/1024),
				Completed: fmt.Sprintf("%.2f%%", completed),
				FileName:  file,
			})
		}

	}
	return completedFiles
}
