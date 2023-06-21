package types

const (
	VideoTypeMovie = "movie"
	VideoTypeTV    = "tv"
)

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
	Video *DouBanVideo
	File  string
	Size  string
	Gid   string
}
