package aria2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zyxar/argo/rpc"
	"io"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"
	"strconv"
	"strings"
	"sync"
	"time"
)

//nolint:gochecknoglobals
var (
	aria2Client *Aria2
	once        sync.Once
)

type Aria2 struct {
	aria2Client  rpc.Client
	downloadTask map[string]*types.FeedVideo
	mtx          sync.Mutex
}

// NewAria2
//
//	@Description: 初始化aria2
//	@param label
//	@return *Aria2
//	@return error
//
//nolint:exhaustruct
func NewAria2(label string) (*Aria2, error) {
	var e error
	once.Do(func() {
		for _, v := range config.Config.Aria2cList {
			var URL string
			if strings.HasSuffix(v.URL, "jsonrpc") {
				URL = v.URL
			} else {
				URL = v.URL + "/jsonrpc"
			}

			if v.Label == label {
				client, err := rpc.New(context.TODO(), URL, v.Token, 0, nil)
				if err != nil {
					log.Error(err)
					e = err
					return
				}
				marshal, _ := json.Marshal(config.Config.Aria2cList)
				log.Debug(string(marshal))
				aria2Client = &Aria2{aria2Client: client, downloadTask: make(map[string]*types.FeedVideo)}
			}
		}
	})
	if e != nil {
		return nil, e
	}
	return aria2Client, nil
}

// DownloadByURL
//
//	@Description: 通过url下载
//	@receiver a
//	@param url
//	@return gid
//	@return err
//
//nolint:wrapcheck
func (a *Aria2) DownloadByURL(url string) (gid string, err error) {
	return a.aria2Client.AddURI([]string{url})
}

func (a *Aria2) DownloadByMagnet(magnet string) (gid string, err error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	// 添加磁链
	MateGid, err := a.aria2Client.AddURI([]string{magnet})
	if err != nil {
		return "", errors.WithMessage(err, "AddURI")
	}

	// 超时时间
	timeout := time.After(5 * time.Minute) // 设置超时时间为10秒
	//  要等种子下载完毕后 再获取活动的下载任务
	for {
		select {
		case <-timeout:
			// 达到超时时间，执行相应的逻辑
			return "", errors.New("TellStatus 超时")
		default:
			time.Sleep(1 * time.Second)
			info, err := a.aria2Client.TellStatus(MateGid, "files", "gid", "status", "errorMessage", "errorCode", "followedBy")
			if err != nil {
				log.Error(err)
				return "", errors.WithMessage(err, "TellStatus")
			}
			// active  waiting   paused   error   complete   removed
			if info.Status == "complete" {
				// 如果有 followedBy 说明是磁链下载的种子
				if len(info.FollowedBy) > 0 {
					return info.FollowedBy[0], nil
				}
				return info.Gid, nil
			}
			if info.Status == "error" {
				if info.ErrorCode == "12" {
					return info.Gid, nil
				}
				msg := fmt.Sprintf("code: %s, msg: %s", info.ErrorCode, info.ErrorMessage)
				log.Error(msg)
				return "", errors.New(msg)
			}
		}
	}
}

/*
// getAllActiveGID
//
//	@Description: 获取所有正在下载的任务的gid
//	@receiver a
//	@return []string
//	@return error
//
//nolint:prealloc

	func (a *Aria2) getAllActiveGID() ([]string, error) {
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
*/

func (a *Aria2) DownloadByWithVideo(v *types.FeedVideo, url string) (gid string, err error) {
	gid, err = a.DownloadByMagnet(url)
	if err != nil {
		return "", err
	}
	if v != nil {
		a.AddDownloadTask(v, gid)
	}
	return
}

// List
//
//	@Description: 获取当前正在下载的文件列表
//	@receiver a
//	@param url
//	@return gid
//	@return err
//
//nolint:wrapcheck
func (a *Aria2) List() (infos []rpc.StatusInfo, err error) {
	return a.aria2Client.TellActive()
}

// CurrentActiveAndStopFiles
//
//	@Description: 获取当前正在下载以及停止下载的文件
//	@receiver a
//	@return completedFiles
func (a *Aria2) CurrentActiveAndStopFiles() (completedFiles []*types.ReportCompletedFiles) {
	// 获取已停止下载的文件
	sessionInfo, err := a.aria2Client.TellStopped(0, 100)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Error(err)
		return nil
	}

	completedFiles1 := a.completedHandler(sessionInfo, completedFiles...)
	completedFiles = append(completedFiles, completedFiles1...)
	// 获取正在下载的文件
	ActiveSession, err := a.aria2Client.TellActive()
	if err != nil && !errors.Is(err, io.EOF) {
		log.Error(err)
		return nil
	}
	completedFiles2 := a.completedHandler(ActiveSession, completedFiles...)
	completedFiles = append(completedFiles, completedFiles2...)
	return
}

// completedHandler
//
//	@Description: 处理已完成的文件
//	@receiver a
//	@param sessionInfo
//	@param completedFiles
//	@return []*types.ReportCompletedFiles
func (a *Aria2) completedHandler(sessionInfo []rpc.StatusInfo, completedFiles ...*types.ReportCompletedFiles) []*types.ReportCompletedFiles {
	for _, v := range sessionInfo {
		file, size := getMaxSizeFile(v.Files)
		if file == "" {
			continue
		}
		if strings.Contains(file, "[METADATA]") {
			continue
		}
		// 下载了多少
		CompletedLength, err := strconv.Atoi(v.CompletedLength)
		if err != nil {
			log.Error(err)
		}
		// 文件完成度百分比
		completed := float32(CompletedLength) / float32(size) * 100
		completedFiles = append(completedFiles, &types.ReportCompletedFiles{
			GID:       v.Gid,
			Size:      tools.ByteCountBinary(int64(size)),
			Completed: fmt.Sprintf("%.2f%%", completed),
			FileName:  file,
		})
	}
	return completedFiles
}

// AddDownloadTask
//
//	@Description: 添加下载任务
//	@receiver a
//	@param douBanVideo
//	@param gid
func (a *Aria2) AddDownloadTask(feedVideo *types.FeedVideo, gid string) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.downloadTask[gid] = feedVideo
}
func (a *Aria2) GetDownloadTask() map[string]*types.FeedVideo {
	return a.downloadTask
}

func (a *Aria2) Subscribe(downLoadChan chan *types.DownloadNotifyVideo) {
	go func() {
		a.mtx.Lock()
		for gid, feedVideo := range a.downloadTask {
			info, err := a.aria2Client.TellStatus(gid, "files", "status")
			if err != nil {
				log.Error(err)
				continue
			}
			// active  waiting   paused   error   complete   removed
			if info.Status == "complete" {
				file, size := getMaxSizeFile(info.Files)
				downLoadChan <- &types.DownloadNotifyVideo{
					FeedVideo: feedVideo,
					File:      file,
					Size:      tools.ByteCountBinary(int64(size)),
					Gid:       gid,
				}
				delete(a.downloadTask, gid)
			}
		}
		a.mtx.Unlock()
	}()
}

func getMaxSizeFile(files []rpc.FileInfo) (string, int) {
	var maxSizeFile int
	var f rpc.FileInfo
	for _, file := range files {
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
