package tools

import (
	"testing"
)

func TestTorrentName2info(t *testing.T) {
	info, s, s2, err := TorrentName2info(`[MagicStar].Yukionna.to.Kani.wo.Kuu.[WEBDL].[1080p].[AMZN]`)
	if err != nil {
		t.Error(err)
	}
	t.Log(info, s, s2)
}
