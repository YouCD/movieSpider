package aria2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/tools"
	"movieSpider/internal/types"

	"github.com/spf13/cast"
	"github.com/youcd/toolkit/log"
	"github.com/zyxar/argo/rpc"
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
func NewAria2(label string) (*Aria2, error) {
	var e error
	once.Do(func() {
		for _, v := range config.Config.Aria2cList {
			URL := v.URL + "/jsonrpc"
			if strings.HasSuffix(v.URL, "jsonrpc") {
				URL = v.URL
			}

			if v.Label == label {
				client, err := rpc.New(context.Background(), URL, v.Token, 0, nil)
				if err != nil {
					log.WithCtx(context.Background()).Error(err)
					e = err
					return
				}
				marshal, _ := json.Marshal(config.Config.Aria2cList)
				log.WithCtx(context.Background()).Debug(string(marshal))
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

var ErrTellStatus = errors.New("TellStatus 超时")

func (a *Aria2) DownloadByMagnet(magnet string) (gid string, err error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	// 添加磁链
	MateGid, err := a.aria2Client.AddURI([]string{magnet})
	if err != nil {
		return "", fmt.Errorf("AddURI err:%w", err)
	}

	// 超时时间
	timeout := time.After(5 * time.Minute) // 设置超时时间为10秒
	//  要等种子下载完毕后 再获取活动的下载任务
	for {
		select {
		case <-timeout:
			// 达到超时时间，执行相应的逻辑
			return "", ErrTellStatus
		default:
			time.Sleep(1 * time.Second)
			info, err := a.aria2Client.TellStatus(MateGid, "files", "gid", "status", "errorMessage", "errorCode", "followedBy")
			if err != nil {
				log.WithCtx(context.Background()).Error(err)
				return "", fmt.Errorf("TellStatus err:%w", err)
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
				log.WithCtx(context.Background()).Error(msg)
				//nolint:err113
				return "", errors.New(msg)
			}
		}
	}
}

func (a *Aria2) DownloadByWithVideo(ctx context.Context, v *types.FeedVideo, url string) (gid string, err error) {
	gid, err = a.DownloadByMagnet(url)
	if err != nil {
		return "", err
	}
	if v == nil {
		return gid, nil
	}
	a.AddDownloadTask(v, gid)
	//  将下载信息存储到数据库中
	err = a.updateSeasonInfo(ctx, v)
	if err != nil {
		log.WithCtx(ctx).Error(err)
	}
	return
}

func (a *Aria2) updateSeasonInfo(ctx context.Context, v *types.FeedVideo) error {
	db := model.NewMovieDB()
	video, err := db.FetchOneTMDBVideoByImdbID(ctx, v.ImdbID)
	if err != nil {
		return err
	}
	var info []*types.SeasonInfo
	if err = json.Unmarshal([]byte(video.SeasonInfo), &info); err != nil {
		return err
	}

	for _, seasonInfo := range info {
		if seasonInfo.S == v.Season && seasonInfo.E == v.Episode {
			seasonInfo.Status = "download"
			break
		}
	}
	marshal, _ := json.Marshal(info)
	video.SeasonInfo = string(marshal)
	return db.UpdateTMDBVideo(ctx, video)
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

// RemoveTask 删除下载任务（支持活动和已完成的任务）
func (a *Aria2) RemoveTask(gid string) error {
	// 先尝试删除活动任务
	_, err := a.aria2Client.Remove(gid)
	if err != nil {
		// 如果删除活动任务失败，尝试删除已完成任务的记录
		_, err = a.aria2Client.RemoveDownloadResult(gid)
		return err
	}
	return nil
}

// ForceRemoveTask 强制删除下载任务（支持活动和已完成的任务）
func (a *Aria2) ForceRemoveTask(gid string) error {
	// 先尝试强制删除活动任务
	_, err := a.aria2Client.ForceRemove(gid)
	if err != nil {
		// 如果删除活动任务失败，尝试删除已完成任务的记录
		_, err = a.aria2Client.RemoveDownloadResult(gid)
		return err
	}
	return nil
}

// CurrentActiveAndStopFiles
//
//	@Description: 获取当前正在下载以及停止下载的文件
//	@receiver a
//	@return completedFiles
func (a *Aria2) CurrentActiveAndStopFiles() []*types.ReportCompletedFiles {
	// 用于去重
	seenGIDs := make(map[string]bool)
	var completedFiles []*types.ReportCompletedFiles

	// 获取已停止下载的文件
	sessionInfo, err := a.aria2Client.TellStopped(0, 100)
	if err != nil && !errors.Is(err, io.EOF) {
		log.WithCtx(context.Background()).Error(err)
		return nil
	}

	completedFiles = a.completedHandler(sessionInfo)
	for _, f := range completedFiles {
		seenGIDs[f.GID] = true
	}

	// 获取正在下载的文件
	ActiveSession, err := a.aria2Client.TellActive()
	if err != nil && !errors.Is(err, io.EOF) {
		log.WithCtx(context.Background()).Error(err)
		return nil
	}
	activeFiles := a.completedHandler(ActiveSession)

	// 添加活动文件，跳过已存在的GID
	for _, f := range activeFiles {
		if !seenGIDs[f.GID] {
			completedFiles = append(completedFiles, f)
		}
	}
	return completedFiles
}

func getFile(info rpc.StatusInfo) string {
	var filename string
	for _, f := range info.Files {
		// 处理 METADATA 类型的文件（磁力链接下载的种子）
		if strings.HasPrefix(f.Path, "[METADATA]") {
			// 提取 METADATA 后面的内容作为文件名
			filename = strings.TrimPrefix(f.Path, "[METADATA]")
			continue
		}
		s := strings.Split(f.Path, "/")

		if len(s) >= 3 {
			filename = s[2]
		}
	}
	return filename
}

// AddDownloadTask
//
//	@Description: 添加下载任务
//	@receiver a
//	@param feedVideo
//	@param gid
func (a *Aria2) AddDownloadTask(feedVideo *types.FeedVideo, gid string) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.downloadTask[gid] = feedVideo
}

func (a *Aria2) GetDownloadTask() map[string]*types.FeedVideo {
	return a.downloadTask
}

// Subscribe
//
//	@Description: 检查下载任务中已完成的任务，返回完成的任务列表
//	@receiver a
//	@return []*types.DownloadNotifyVideo 已完成的下载通知列表
func (a *Aria2) Subscribe() []*types.DownloadNotifyVideo {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	var completedVideos []*types.DownloadNotifyVideo
	for gid, feedVideo := range a.downloadTask {
		info, err := a.aria2Client.TellStatus(gid, "files", "status")
		if err != nil {
			log.WithCtx(context.Background()).Error(err)
			continue
		}
		// active  waiting   paused   error   complete   removed
		if info.Status == "complete" {
			file, size := getMaxSizeFile(info.Files)
			completedVideos = append(completedVideos, &types.DownloadNotifyVideo{
				FeedVideo: feedVideo,
				File:      file,
				Size:      tools.ByteCountBinary(int64(size)),
				Gid:       gid,
			})
			delete(a.downloadTask, gid)
		}
	}
	return completedVideos
}

// completedHandler
//
//	@Description: 处理已完成的文件
//	@receiver a
//	@param sessionInfo
//	@return []*types.ReportCompletedFiles
func (a *Aria2) completedHandler(sessionInfo []rpc.StatusInfo) []*types.ReportCompletedFiles {
	var completedFiles []*types.ReportCompletedFiles
	for _, v := range sessionInfo {
		file := getFile(v)
		if file == "" {
			continue
		}

		// 获取总大小（优先使用 TotalLength，如果为空则计算文件总大小）
		var totalSize int64
		if v.TotalLength != "" {
			totalSize = cast.ToInt64(v.TotalLength)
		} else {
			for _, f := range v.Files {
				totalSize += cast.ToInt64(f.Length)
			}
		}

		// 计算完成度百分比
		var completedPercent float32
		completedLength := cast.ToInt64(v.CompletedLength)
		if totalSize > 0 {
			completedPercent = float32(completedLength) / float32(totalSize) * 100
		}

		// 根据状态设置完成度显示
		completedStr := fmt.Sprintf("%.2f%%", completedPercent)
		status := v.Status
		var errorMsg string

		switch v.Status {
		case "complete":
			completedStr = "100%"
		case "error":
			if v.ErrorMessage != "" {
				errorMsg = v.ErrorMessage
				completedStr = fmt.Sprintf("错误: %s", v.ErrorMessage)
			} else {
				errorMsg = fmt.Sprintf("code: %s", v.ErrorCode)
				completedStr = fmt.Sprintf("错误(code: %s)", v.ErrorCode)
			}
		case "paused":
			completedStr = fmt.Sprintf("已暂停 %s", completedStr)
		case "waiting":
			completedStr = fmt.Sprintf("等待中 %s", completedStr)
		case "removed":
			completedStr = "已移除"
		}

		completedFiles = append(completedFiles, &types.ReportCompletedFiles{
			GID:       v.Gid,
			Size:      tools.ByteCountBinary(totalSize),
			Completed: completedStr,
			FileName:  file,
			Status:    status,
			ErrorMsg:  errorMsg,
		})
	}
	return completedFiles
}

func getMaxSizeFile(files []rpc.FileInfo) (string, int) {
	var maxSizeFile int
	var f rpc.FileInfo
	for _, file := range files {
		if strings.HasPrefix(file.Path, "[METADATA]") {
			continue
		}
		if cast.ToInt(file.Length) > maxSizeFile {
			maxSizeFile = cast.ToInt(file.Length)
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
