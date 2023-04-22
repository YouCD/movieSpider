package searchSpider

import (
	"fmt"
	"sync"
	"testing"
)

func Test_crawlerRarbg(t *testing.T) {
	//rarbg, err := crawlerRarbg("tt11198330")
	rarbgs, err := crawlerRarbg("tt7631058")
	if err != nil {
		t.Error(err)
		return
	}

	var wg sync.WaitGroup
	var Videos []*RarbgVideo
	for _, v := range rarbgs {
		//fmt.Println("cc", v)
		var v1 = *v
		go func(v2 RarbgVideo) {
			wg.Add(1)
			magnet, err := crawlerRarbgMagnet(v2.TorrentUrl)
			if err != nil {
				t.Error(err)
				return
			}
			v2.Magnet = magnet
			//fmt.Println("11", v2)
			Videos = append(Videos, &v2)
			wg.Done()
		}(v1)
	}
	wg.Wait()
	for _, v := range Videos {
		fmt.Println("vvvv", v)
	}

}

func Test_crawlerRarbgMagnet(t *testing.T) {
	gotVideos, err := crawlerRarbgMagnet(rarbgTorrentUrlPrefix + "/torrent/fzy8xp4")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(gotVideos)
}
