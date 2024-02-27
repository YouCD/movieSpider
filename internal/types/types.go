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
}
type DownloadNotifyVideo struct {
	FeedVideo   *FeedVideo
	DouBanVideo *DouBanVideo
	File        string
	Size        string
	Gid         string
}
