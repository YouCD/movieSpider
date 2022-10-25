package feed

import (
	"github.com/pkg/errors"
	"movieSpider/internal/convert"
	httpClient2 "movieSpider/internal/httpClient"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"net/http"
	"strings"
)

func fetchMagnet(url string) (magnet string, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.WithMessage(err, "TGx: 磁链获取错误")
	}
	client := httpClient2.NewHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		return "", errors.WithMessage(err, "TGx: 磁链获取错误")
	}
	if resp == nil {
		return "", errors.New("TGx: response is nil")
	}
	defer resp.Body.Close()

	magnet, err = convert.IO2Magnet(resp.Body)
	if err != nil {
		return "", errors.New("TGx: 磁链转换错误")
	}

	return magnet, nil
}

func proxySaveVideo2DB(videos []*types.FeedVideo) {
	if videos == nil || len(videos) == 0 {
		log.Warn("没有数据")
		return
	}

	for _, v := range videos {
		go func(video *types.FeedVideo) {
			err := model.NewMovieDB().CreatFeedVideo(video)
			if err != nil {
				if errors.Is(err, model.ErrorDataExist) {
					log.Warnf("%s.%s err: %s", strings.ToUpper(video.Web), video.Type, err)
					return
				}
				log.Error(err)
				return
			}
			log.Infof("%s.%s: %s 保存完毕.", strings.ToUpper(video.Web), video.Type, video.Name)
		}(v)
	}
}

// excludeVideo 排除  480p 720p  dvsux  hdr 视频源
func excludeVideo(name string) bool {
	lowerTorrentName := strings.ToLower(name)
	if strings.Contains(lowerTorrentName, "720p") || strings.Contains(lowerTorrentName, "dvsux") || strings.Contains(lowerTorrentName, "480p") || strings.Contains(lowerTorrentName, "hdr") {
		return true
	} else {
		return false
	}
}
