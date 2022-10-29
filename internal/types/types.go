package types

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"
)

type FeedVideo struct {
	ID          int32  `json:"id"`
	TorrentName string `json:"torrent_name"` // 种子名
	TorrentUrl  string `json:"torrent_url"`  // 种子引用地址
	Magnet      string `json:"magnet"`       // 磁力链接
	Year        string `json:"year"`         // 年份
	Name        string `json:"name"`         // 片名
	RowData     string `json:"row_data"`     // 原始数据
	Type        string `json:"type"`         // tv或movie
	Web         string `json:"web"`          // 站点
	Download    int    `json:"download"`     // 1:已经下载
}

func (f *FeedVideo) FormatName(name string) string {
	// 去除空格
	name = strings.ReplaceAll(name, " ", "")

	// 处理 .
	nameSlice := strings.Split(name, ".")
	ret := removeSpaceItem(nameSlice)
	name = strings.Join(ret, ".")
	// 去除 -.
	name = strings.ReplaceAll(name, ".-.", ".")
	// 去除 +.
	name = strings.ReplaceAll(name, ".+.", ".")

	return name
}

func removeSpaceItem(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

type DouBanVideo struct {
	ID       int32  `json:"id"`
	Names    string `json:"names"`
	DoubanID string `json:"douban_id"`
	ImdbID   string `json:"imdb_id"`
	RowData  string `json:"row_data"`
	Type     string `json:"type"`     // tv或movie
	Playable string `json:"playable"` //是否可以播放
}

func (d *DouBanVideo) FormatName(names string) string {
	var n []string

	split := strings.Split(names, "/")
	for _, name := range split {
		// 处理 空格
		nameSlice := strings.Split(name, " ")
		ret := removeSpaceItem(nameSlice)
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
