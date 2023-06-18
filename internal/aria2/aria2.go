package aria2

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zyxar/argo/rpc"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	aria2Client *aria2
	once        sync.Once
)

type downloadVideoComplete struct {
	Video *types.DouBanVideo
	File  string
	Size  string
}
type aria2 struct {
	aria2Client  rpc.Client
	downloadTask map[string]*types.DouBanVideo
	mtx          sync.Mutex
}

//
// NewAria2
//  @Description: 初始化aria2
//  @param label
//  @return *aria2
//  @return error
//
func NewAria2(label string) (*aria2, error) {
	var e error
	once.Do(func() {
		for _, v := range config.Aria2cList {
			if v.Label == label {
				client, err := rpc.New(context.TODO(), v.Url, v.Token, 0, nil)
				if err != nil {
					log.Error(err)
					e = err
				}
				log.Debug(config.Aria2cList)
				aria2Client = &aria2{aria2Client: client, downloadTask: make(map[string]*types.DouBanVideo)}
			}
		}
	})
	if e != nil {
		return nil, e
	}
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

func (a *aria2) DownloadByMagnet(magnet string) (gid string, err error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	//// 下载前先获得当前所有的下载任务 并记录 gid
	//beforeActiveGID, err := a.getAllActiveGID()
	//if err != nil {
	//	return "", err
	//}

	// 添加磁链
	MateGid, err := a.aria2Client.AddURI([]string{magnet})
	if err != nil {
		return "", err
	}

	files, err := a.aria2Client.GetFiles(MateGid)
	if err != nil {
		return "", err
	}

	Metafile, _ := getMaxSizeFile(files)
	Metafile = strings.ReplaceAll(Metafile, "[METADATA]", "")
	// 超时时间
	timeout := time.After(5 * time.Minute) // 设置超时时间为10秒
	complete := false                      // 标志变量，表示循环是否已完成
	//  要等种子下载完毕后 再获取活动的下载任务
	for !complete {
		select {
		case <-timeout:
			// 达到超时时间，执行相应的逻辑
			return "", errors.New("TellStatus 超时")
		default:
			time.Sleep(1 * time.Second)
			info, err := a.aria2Client.TellStatus(MateGid, "files", "status", "errorMessage", "errorCode")
			if err != nil {
				log.Error(err)
				return "", err
			}
			// active  waiting   paused   error   complete   removed
			if info.Status == "complete" {
				complete = true // 循环完成
				break
			}
			if info.Status == "error" {
				if info.ErrorCode == "12" {
					return "", errors.New("种子已经下载过了")
				}
				msg := fmt.Sprintf("code: %s, msg: %s", info.ErrorCode, info.ErrorMessage)
				log.Error(msg)
				return "", errors.New(msg)
			}
		}
	}
	afterActiveGID, err := a.getAllActiveGID()
	if err != nil {
		return "", err
	}
	// 从所有的下载任务中 找到新添加的任务  也就是磁链的任务
	if len(afterActiveGID) == 0 {
		return "", err
	}
	gid = afterActiveGID[len(afterActiveGID)-1]

	return
}

//
// getAllActiveGID
//  @Description: 获取所有正在下载的任务的gid
//  @receiver a
//  @return []string
//  @return error
//
func (a *aria2) getAllActiveGID() ([]string, error) {
	infos, err := a.List()
	if err != nil {
		return nil, err
	}
	var gid []string
	for _, info := range infos {
		gid = append(gid, info.Gid)
	}
	return gid, nil
}

func (a *aria2) DownloadByWithVideo(v *types.DouBanVideo, url string) (gid string, err error) {
	gid, err = a.DownloadByMagnet(url)
	if err != nil {
		return "", err
	}
	if v != nil {
		a.AddDownloadTask(v, gid)
	}
	return
}

//
// List
//  @Description: 获取当前正在下载的文件列表
//  @receiver a
//  @param url
//  @return gid
//  @return err
//
func (a *aria2) List() (infos []rpc.StatusInfo, err error) {
	return a.aria2Client.TellActive()
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

//
// AddDownloadTask
//  @Description: 添加下载任务
//  @receiver a
//  @param douBanVideo
//  @param gid
//
func (a *aria2) AddDownloadTask(douBanVideo *types.DouBanVideo, gid string) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.downloadTask[gid] = douBanVideo
}

func (a *aria2) Subscribe() chan *downloadVideoComplete {
	downLoadChan := make(chan *downloadVideoComplete)
	go func() {
		for gid, video := range a.downloadTask {
			time.Sleep(time.Second * 1)
			info, err := a.aria2Client.TellStatus(gid, "files", "status")
			if err != nil {
				log.Error(err)
			}
			// active  waiting   paused   error   complete   removed
			if info.Status == "complete" {
				file, size := getMaxSizeFile(info.Files)
				downLoadChan <- &downloadVideoComplete{
					Video: video,
					File:  file,
					Size:  tools.ByteCountDecimal(int64(size)),
				}
				a.mtx.Lock()
				delete(a.downloadTask, gid)
				a.mtx.Unlock()
			}

		}
		close(downLoadChan)
	}()
	return downLoadChan
}

func getMaxSizeFile(Files []rpc.FileInfo) (string, int) {
	var maxSizeFile int
	var f rpc.FileInfo
	for _, file := range Files {
		if len(file.Length) > maxSizeFile {
			maxSizeFile = len(file.Length)
			f = file
		}
	}
	filename := ""
	Length := 0
	if f.Length != "" {
		Length, _ = strconv.Atoi(f.Length)
		s := strings.Split(f.Path, "/")
		if len(s) > 0 {
			filename = s[len(s)-1]
		}
	}
	return filename, Length
}
