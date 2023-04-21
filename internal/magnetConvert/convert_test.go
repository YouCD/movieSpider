package magnetConvert

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"testing"
)

func TestConvertFileToMagnet(t *testing.T) {
	Magnet, err := FileToMagnet("/tmp/a.torrent")
	if err != nil {
		t.Error(err)
	}
	t.Log(Magnet)
}

func TestIO2Magnet(t *testing.T) {
	url := "magnet:?xt=urn:btih:5B7650FD0AD30DF53A6F320DC21A44601824B575&dn=Moonshine.S02E03.1080p.WEBRip.x264-BAE%5Beztv.re%5D.mkv&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftorrent.gresille.org%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337"
	m, err := metainfo.ParseMagnetUri(url)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(m.InfoHash)

}
