package feedspider

import (
	"testing"
)

func TestNewGlodls(t *testing.T) {
	feeder := NewGlodls()
	videos, err := feeder.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		t.Log(video)
	}

}
