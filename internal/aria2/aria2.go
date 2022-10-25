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

func (a *aria2) DownloadByUrl(url string) (gid string, err error) {
	return a.aria2Client.AddURI([]string{url})
}
func (a *aria2) DownloadList(url string) (gid string, err error) {
	info, err := a.aria2Client.GetSessionInfo()
	if err != nil {
		return "", err
	}
	fmt.Println(info)
	return
}

func (a *aria2) CompletedFiles() (completedFiles []*types.ReportCompletedFiles) {
	sessionInfo, err := a.aria2Client.TellStopped(0, 100)
	if err != nil {
		log.Error(err)
		return nil
	}

	completedFiles1 := a.completedHandler(sessionInfo, completedFiles)
	completedFiles = append(completedFiles, completedFiles1...)
	ActiveSession, err := a.aria2Client.TellActive()
	if err != nil {
		log.Error(err)
		return nil
	}
	completedFiles2 := a.completedHandler(ActiveSession, completedFiles)
	completedFiles = append(completedFiles, completedFiles2...)
	return
}

func (a *aria2) completedHandler(sessionInfo []rpc.StatusInfo, completedFiles []*types.ReportCompletedFiles) []*types.ReportCompletedFiles {
	for _, v := range sessionInfo {
		if len(v.Files) > 0 {
			if strings.Contains(v.Files[0].Path, "[METADATA]") {
				continue
			} else {
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

				f := new(types.ReportCompletedFiles)
				f.GID = v.Gid
				f.Completed = fmt.Sprintf("%.2f%%", completed)
				f.Size = fmt.Sprintf("%.2fGB", float32(Length)/1024/1024/1024)
				_, file := path.Split(v.Files[0].Path)
				f.FileName = file
				completedFiles = append(completedFiles, f)
			}
		}

	}
	return completedFiles
}
