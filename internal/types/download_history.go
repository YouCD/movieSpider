package types

type DownloadHistory struct {
	ID          int    `gorm:"column:id;type:int(11);AUTO_INCREMENT;primary_key" json:"id"`
	Name        string `gorm:"uniqueIndex:tv_index;column:name;type:varchar(255);comment:片名;NOT NULL" json:"name"`
	Type        string `gorm:"uniqueIndex:tv_index;column:type;type:varchar(255);comment:tv或movie;NOT NULL" json:"type"`
	TorrentName string `gorm:"column:torrent_name;type:varchar(255);comment:种子名;NOT NULL" json:"torrent_name"`
	Timestamp   int64  `gorm:"column:timestamp;type:bigint(11);comment:修改创建时间;NOT NULL" json:"timestamp"`
	Resolution  int64  `gorm:"column:resolution;type:bigint(11);comment:分辨率;NOT NULL" json:"resolution"`
	Season      string `gorm:"uniqueIndex:tv_index;column:season;type:varchar(3);comment:季数;NOT NULL" json:"season"`
	Episode     string `gorm:"uniqueIndex:tv_index;column:episode;type:varchar(3);comment:集数;NOT NULL" json:"episode"`
}

func (m *DownloadHistory) TableName() string {
	return "download_history"
}
