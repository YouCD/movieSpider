package types

import (
	"encoding/json"
	"movieSpider/internal/tools"
	"regexp"
	"strings"
	"unicode"
)

type DouBanVideo struct {
	ID        int    `gorm:"column:id;type:int(11);AUTO_INCREMENT;primary_key" json:"id"`
	Names     string `gorm:"uniqueIndex;column:names;type:varchar(255);comment:片名列表;NOT NULL" json:"names"`
	DoubanID  string `gorm:"column:douban_id;type:varchar(255);comment:豆瓣ID;NOT NULL" json:"douban_id"`
	ImdbID    string `gorm:"column:imdb_id;type:varchar(255);comment:imdbID;NOT NULL" json:"imdb_id"`
	RowData   string `gorm:"column:row_data;type:longtext;comment:原始数据;NOT NULL" json:"row_data"`
	Timestamp int64  `gorm:"column:timestamp;type:bigint(11);comment:修改创建时间;NOT NULL" json:"timestamp"`
	Type      string `gorm:"column:type;type:varchar(255);comment:类型;NOT NULL" json:"type"`
	Playable  string `gorm:"column:playable;type:varchar(255);comment:是否可以播放;NOT NULL" json:"playable"`
}

func (d *DouBanVideo) TableName() string {
	return "douban_video"
}

func (d *DouBanVideo) FormatName(names string) string {
	var n []string

	split := strings.Split(names, "/")
	for _, name := range split {
		// 处理 空格
		nameSlice := strings.Split(name, " ")
		ret := tools.RemoveSpaceItem(nameSlice)
		name = strings.Join(ret, ".")
		name = strings.ReplaceAll(name, ":.", ":")
		name = strings.ReplaceAll(name, "..", ".")

		if d.Type == "tv" {
			ok := d.isChineseChar(name)
			if ok {
				compileRegex := regexp.MustCompile("(.*)\\.第.季")
				matchArr := compileRegex.FindStringSubmatch(name)
				if len(matchArr) > 1 {
					n = append(n, matchArr[1])
				}
				continue
			}
			eRegex := regexp.MustCompile("(.*)\\.Season\\.\\d+")
			EmatchArr := eRegex.FindStringSubmatch(name)
			if len(EmatchArr) > 1 {
				n = append(n, EmatchArr[1])
				continue
			}
			n = append(n, name)
		}
		if d.Type == "movie" {
			n = append(n, name)
		}
	}

	marshal, err := json.Marshal(n)
	if err != nil {
		return ""
	}

	return string(marshal)
}
func (d *DouBanVideo) FormatType(typ string) string {
	if strings.ToLower(typ) == "tvseries" {
		return "tv"
	}
	return "movie"
}
func (d *DouBanVideo) isChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}