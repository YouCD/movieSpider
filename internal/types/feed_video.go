package types

import (
	"database/sql"
	"fmt"
	"movieSpider/internal/tools"
	"regexp"
	"strings"
)

type FeedVideo struct {
	ID          int32          `gorm:"column:id;type:int(11);AUTO_INCREMENT;primary_key" json:"id"`
	Name        string         `gorm:"uniqueIndex:nt;column:name;type:varchar(255);comment:片名;NOT NULL" json:"name"`
	TorrentName string         `gorm:"uniqueIndex:nt;column:torrent_name;type:varchar(255);comment:种子名;NOT NULL" json:"torrent_name"`
	TorrentUrl  string         `gorm:"column:torrent_url;type:varchar(255);comment:种子引用地址;NOT NULL" json:"torrent_url"`
	Magnet      string         `gorm:"column:magnet;type:longtext;comment:磁力链接;NOT NULL" json:"magnet"`
	Year        string         `gorm:"column:year;type:varchar(255);comment:年份;NOT NULL" json:"year"`
	Type        string         `gorm:"column:type;type:varchar(255);comment:tv或movie;NOT NULL" json:"type"`
	RowData     sql.NullString `gorm:"column:row_data;type:longtext;comment:原始数据" json:"row_data"`
	Web         string         `gorm:"column:web;type:varchar(255);comment:站点;NOT NULL" json:"web"`
	Download    int            `gorm:"column:download;type:int(11);comment:1:已经下载;NOT NULL" json:"download"`
	Timestamp   int64          `gorm:"column:timestamp;type:bigint(11);comment:修改创建时间;NOT NULL" json:"timestamp"`
}

func (f *FeedVideo) TableName() string {
	return "feed_video"
}

func (f *FeedVideo) FormatName(name string) string {
	// 去除空格
	name = strings.ReplaceAll(name, " ", "")

	// 处理 .
	nameSlice := strings.Split(name, ".")
	ret := tools.RemoveSpaceItem(nameSlice)
	name = strings.Join(ret, ".")
	// 去除 -.
	name = strings.ReplaceAll(name, ".-.", ".")
	// 去除 +.
	name = strings.ReplaceAll(name, ".+.", ".")

	return name
}

var (
	tvReg         = regexp.MustCompile("[Ss]([0-9][0-9])[eE]([0-9][0-9])")
	resolutionReg = regexp.MustCompile("(2160p|2160P|1080p|1080P)")
)

func (f *FeedVideo) Convert2DownloadHistory() *DownloadHistory {
	var downloadHistory DownloadHistory
	downloadHistory.Name = f.Name
	downloadHistory.TorrentName = f.TorrentName
	downloadHistory.Type = f.Type
	switch f.Type {
	case VideoTypeTV:
		TVNameArr := tvReg.FindStringSubmatch(f.TorrentName)
		if len(TVNameArr) <= 2 {
			return nil
		}
		downloadHistory.Resolution = parseResolution(f.TorrentName)
		downloadHistory.Episode = TVNameArr[2]
		downloadHistory.Season = TVNameArr[1]
	case VideoTypeMovie:
		downloadHistory.Resolution = parseResolution(f.TorrentName)
	}
	//log.Errorf("%#v", downloadHistory)

	return &downloadHistory
}

//
// parseResolution
//  @Description: 解析分辨率
//  @param torrentName
//  @return resolution
//
func parseResolution(torrentName string) (resolution int64) {

	resolutionArr := resolutionReg.FindStringSubmatch(torrentName)
	if len(resolutionArr) < 1 {
		return 0
	}
	//  转换为int
	_, err := fmt.Sscanf(resolutionArr[1], "%d", &resolution)
	if err != nil {
		return 0
	}
	return
}
