package feedspider

import (
	"testing"
)

func TestNewBtbt(t *testing.T) {
	feeder := NewBtbt()
	videos, err := feeder.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		t.Log(video)
	}

}
