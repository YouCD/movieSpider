package types

import (
	"fmt"
	"testing"
)

func TestFeedVideo_FormatTvByName(t *testing.T) {
	f := &FeedVideo{
		ID:          1,
		TorrentName: "Deadliest.Catch.the.Viking.Returns.S01E04.Norwegian.Blood.iNTERNAL.1080p.WEB.h264-B2B[rartv]",
		TorrentUrl:  "ccc",
		Type:        "tv",
	}

	f1 := &FeedVideo{
		TorrentName: "Spider-Man.No.Way.Home.2021.EXTENDED.1080p.WEBRip.x264-RARBG",
		TorrentUrl:  "ccc",
		Type:        "movie",
	}

	fmt.Println(f.Convert2DownloadHistory())
	fmt.Printf("%+v\n", f1.Convert2DownloadHistory())
}

func TestFeedVideo_FormatName(t *testing.T) {
	f := &FeedVideo{
		Name: "Black.Adam.(2022).1080p.H265.ita.eng.AC3.5.1.sub.ita.eng.Licdom.mkv",
	}
	fmt.Println(f.FormatName(f.Name))
}
