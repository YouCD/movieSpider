package types

import (
	"fmt"
	"testing"
)

func TestFeedVideo_FormatTvByName(t *testing.T) {
	f := &FeedVideo{
		ID:          1,
		TorrentName: "Deadliest.Catch.the.Viking.Returns.S01E04.Norwegian.Blood.iNTERNAL.1080p.WEB.h264-B2B[rartv]",
		TorrentURL:  "ccc",
		Type:        "tv",
	}

	f1 := &FeedVideo{
		TorrentName: "Spider-Man.No.Way.Home.2021.EXTENDED.1080p.WEBRip.x264-RARBG",
		TorrentURL:  "ccc",
		Type:        "movie",
	}

	fmt.Println(f.Convert2DownloadHistory())
	fmt.Printf("%+v\n", f1.Convert2DownloadHistory())
}

func TestFeedVideo_FormatName(t *testing.T) {
	// FormatName 已迁移到 TMDBVideo 上
	d := &TMDBVideo{
		Type: "tv",
	}
	fmt.Println(d.FormatName(`梦中的那片海 S01 第17-20集 国语配音+中文字幕.Where.Dreams.Begin`))
}
