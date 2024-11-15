package dhtclient

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"movieSpider/internal/bus"
	"movieSpider/internal/types"
	"os"
	"time"
	"unicode"

	"github.com/youcd/toolkit/log"
)

func Boot(thread int) {
	log.Info("Start DHT Network")
	for i := 0; i < thread; i++ {
		go crawl()
	}
	//}
	//select {}
}

func crawl() {
	indexerAddrs := []string{"0.0.0.0:0"}
	interruptChan := make(chan os.Signal, 1)

	trawlingManager := NewManager([]string{
		"router.bittorrent.com:6881", "router.utorrent.com:6881",
		"dht.transmissionbt.com:6881", "dht.libtorrent.org:25401",
	}, indexerAddrs, 1, 1000)
	metadataSink := NewSink(5*time.Second, 256)

	for stopped := false; !stopped; {
		select {
		case result := <-trawlingManager.Output():
			//hash := result.InfoHash()
			//fmt.Println("", hash)
			//if !cache.InfoHashCache.Contains(hash) {
			//	cache.InfoHashCache.Add(hash)
			metadataSink.Sink(result)
			//}

		case md := <-metadataSink.Drain():
			//if db.InsertMetadata(configuration, database, md) {
			//doc.Set("Name", md.Name)
			//doc.Set("InfoHash", hex.EncodeToString(md.InfoHash))
			//doc.Set("Files", md.Files)
			//doc.Set("DiscoveredOn", md.DiscoveredOn)
			//doc.Set("TotalSize", md.TotalSize)
			// magnet:?xt=urn:btih:89c64cb5e479d4759a4a8376e024856cb1e0f7d1
			//fmt.Printf(" Added: %s    magnet:?xt=urn:btih:%s    %v \n", md.Name, hex.EncodeToString(md.InfoHash), md.DiscoveredOn)
			//db.CheckWatches(configuration, database, md, bot)
			//}

			hasHan := false
			for _, v := range md.Name {
				if unicode.Is(unicode.Han, v) {
					log.Warn("中文跳过: ", md.Name)
					hasHan = true
					break
				}
			}

			if hasHan {
				continue
			}
			bus.FeedVideoChan <- &types.FeedVideoBase{
				TorrentName: md.Name,
				TorrentURL:  "",
				Magnet:      fmt.Sprintf("magnet:?xt=urn:btih:%s", hex.EncodeToString(md.InfoHash)),
				Type:        "",
				RowData:     sql.NullString{},
				Year:        "",
				Web:         "DHT",
			}
		case <-interruptChan:
			trawlingManager.Terminate()
			stopped = true
		}
	}
}
