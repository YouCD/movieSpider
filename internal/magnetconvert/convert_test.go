package magnetconvert

import (
	"fmt"
	"testing"

	"github.com/anacrolix/torrent/metainfo"
)

func TestConvertFileToMagnet(t *testing.T) {
	Magnet, err := FileToMagnet("/home/ycd/Downloads/a.torrent")
	if err != nil {
		t.Error(err)
	}
	t.Log(Magnet)
}

func TestIO2Magnet(t *testing.T) {
	url := "magnet:?xt=urn:btih:5B7650FD0AD30DF53A6F320DC21A44601824B575"
	m, err := metainfo.ParseMagnetUri(url)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(m.InfoHash)
	fmt.Println(m.Trackers)

}
