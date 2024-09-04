package feedspider

import (
	"movieSpider/internal/model"
	"testing"

	"github.com/youcd/toolkit/log"
)

func TestThePirateBay_Crawler(t1 *testing.T) {
	thePirateBay := NewThePirateBay()
	gotVideos, err := thePirateBay.Crawler()
	if err != nil {
		t1.Errorf("Crawler() error = %v", err)
		return
	}
	for _, video := range gotVideos {
		filterVideo, err := model.FilterVideo(video)
		if err != nil {
			log.Errorf("err: %s    %#v", err, video)
			continue
		}
		log.Infof("%#v", filterVideo)
	}
}
