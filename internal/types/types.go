package types

import "strings"

type VideoType string

const (
	VideoTypeMovie   VideoType = "movie"
	VideoTypeTV      VideoType = "tv"
	VideoTypeUnknown VideoType = "unknown"
)

func Convert2VideoType(t string) VideoType {
	tt := strings.ToLower(t)
	switch tt {
	case "tv":
		return VideoTypeTV
	case "movie":
		return VideoTypeMovie
	default:
		return VideoTypeUnknown
	}
}

func (v VideoType) String() string {
	return string(v)
}

type ReportCount struct {
	Web   string `json:"web"`
	Count int    `json:"count"`
}

type ReportCompletedFiles struct {
	GID       string
	Size      string
	Completed string
	FileName  string
	Status    string // 下载状态: active, waiting, paused, error, complete, removed
	ErrorMsg  string // 错误信息（仅当状态为error时有值）
}
type DownloadNotifyVideo struct {
	FeedVideo *FeedVideo
	TMDBVideo *TMDBVideo
	File      string
	Size      string
	Gid       string
}
type LLMResult struct {
	ID         int    `json:"id"`
	TypeStr    string `json:"typeStr"`
	NewName    string `json:"newName"`
	Year       int    `json:"year,omitempty"`
	Resolution int    `json:"resolution,omitempty"`
}
