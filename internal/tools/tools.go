package tools

import (
	"fmt"
	"strings"
)

//
// RemoveSpaceItem
//  @Description: 去除数组中的空格
//  @param a
//  @return ret
//
func RemoveSpaceItem(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// ExcludeVideo
// ExcludeVideo 排除  480p 720p  dvsux  hdr 视频源
//  @Description:
//  @param name
//  @return bool
func ExcludeVideo(name string) bool {
	lowerTorrentName := strings.ToLower(name)
	if strings.Contains(lowerTorrentName, "720p") || strings.Contains(lowerTorrentName, "dvsux") || strings.Contains(lowerTorrentName, "480p") || strings.Contains(lowerTorrentName, "hdr") || strings.Contains(lowerTorrentName, ".DV.") {
		return true
	}
	return false
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
