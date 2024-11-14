package tools

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/youcd/toolkit/log"
)

// RemoveSpaceItem
//
//	@Description: 去除数组中的空格
//	@param a
//	@return ret
func RemoveSpaceItem(a []string) (ret []string) {
	aLen := len(a)
	for i := range aLen {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// ExcludeVideo
// ExcludeVideo 排除  480p 720p  dvsux  hdr 视频源
//
//	@Description:
//	@param name
//	@return bool
//
//nolint:dupword
func ExcludeVideo(name string, excludeWords []string) bool {
	lowerTorrentName := strings.ToLower(name)
	for _, word := range excludeWords {
		if strings.Contains(lowerTorrentName, word) {
			return true
		}
	}
	return false
}

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
func TorrentName2info(torrentName string) (string, string, string, error) {
	newTorrentName := ReplaceAll(torrentName)

	var name, resolution, year string
	// 先匹配 tv
	//nolint:revive
	if tvName, _, err := matchTV(newTorrentName); err != nil {
		log.Debugf("tv匹配失败: old:%s, new:%s", torrentName, newTorrentName)
		goto MatchMovie
	} else {
		name = tvName
		goto MatchResolution
	}
	//	 匹配 movie
MatchMovie:
	if movieName, nameRegexArr, err := matchMovie(newTorrentName); err == nil {
		name = movieName
	} else {
		log.Debugf("movie匹配失败: old:%s, new:%s, Arr: %s", torrentName, newTorrentName, nameRegexArr)
		return "", "", "", fmt.Errorf("%w,old:%s, new:%s", err, torrentName, newTorrentName)
	}
MatchResolution:
	// 匹配分辨率
	{
		resolutionRegex := regexp.MustCompile(`.*\.(\d{4}[p|P])\.`)
		resolutionArr := resolutionRegex.FindStringSubmatch(newTorrentName)
		if len(resolutionArr) < 2 || len(resolutionArr) == 0 {
			return "", "", "", fmt.Errorf("%w, old:%s, new:%s", ErrFeedVideoResolution, torrentName, newTorrentName)
		}
		resolution = resolutionArr[1]
	}

	// 匹配 年份
	{
		compileYearRegex := regexp.MustCompile(`.*?(\d{4}).*?\d{4}[p|P].*`)
		yearArr := compileYearRegex.FindStringSubmatch(newTorrentName)
		if len(yearArr) < 2 || len(yearArr) == 0 {
			return "", "", "", fmt.Errorf("%w, old:%s, new:%s", ErrFeedVideoYear, torrentName, newTorrentName)
		}
		year = yearArr[1]
	}

	return name, resolution, year, nil
}

func ReplaceAll(torrentName string) string {
	// 去除 -
	newTorrentName := strings.ReplaceAll(torrentName, "-", ".")
	// 去除 _
	newTorrentName = strings.ReplaceAll(newTorrentName, "_", ".")
	newTorrentName = strings.ReplaceAll(newTorrentName, ",", "")

	// 去除空格
	reg := regexp.MustCompile(`( )+|(\n)+`)
	newTorrentName = reg.ReplaceAllString(newTorrentName, "$1$2")
	newTorrentName = strings.ReplaceAll(newTorrentName, " ", ".")
	newTorrentName = strings.ReplaceAll(newTorrentName, ".", ".") //nolint:gocritic

	dotReg := regexp.MustCompile(`\.+`)
	newTorrentName = dotReg.ReplaceAllString(newTorrentName, ".")

	// 去除 []
	newTorrentName = strings.ReplaceAll(newTorrentName, "[", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, "]", "")

	// 去除 ()
	newTorrentName = strings.ReplaceAll(newTorrentName, "(", "")
	newTorrentName = strings.ReplaceAll(newTorrentName, ")", "")
	return newTorrentName
}
func matchTV(torrentName string) (string, []string, error) {
	tvReg := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9][eE][0-9][0-9])`)
	tvNameArr := tvReg.FindStringSubmatch(torrentName)
	// 如果 正则匹配过后 没有结果直接 过滤掉
	if len(tvNameArr) < 2 || len(tvNameArr) == 0 {
		tvRegA := regexp.MustCompile(`(.*)(\.[Ss][0-9][0-9]).*`)
		tvNameArrA := tvRegA.FindStringSubmatch(torrentName)
		if len(tvNameArrA) >= 2 {
			return tvNameArrA[1], nil, nil
		}
		return "", tvNameArrA, ErrNotMatchTorrentName
	}
	return tvNameArr[1], nil, nil
}

func matchMovie(torrentName string) (string, []string, error) {
	nameRegex := regexp.MustCompile(`(.*)\.\d{4}[p|P]`)
	nameRegexArr := nameRegex.FindStringSubmatch(torrentName)
	if len(nameRegexArr) >= 2 {
		name := nameRegexArr[1]
		movieA, arr, err := matchMovieA(name)
		if err == nil {
			return movieA, arr, nil
		}
		return name, nil, nil
	}
	return matchMovieA(torrentName)
}

func matchMovieA(torrentName string) (string, []string, error) {
	nameReg := regexp.MustCompile(`(.*)\.\d{4}`)
	nameSubMatch := nameReg.FindStringSubmatch(torrentName)
	if len(nameSubMatch) >= 2 {
		return nameSubMatch[1], nil, nil
	}
	return "", nameSubMatch, ErrFeedVideoMovieMatch
}
