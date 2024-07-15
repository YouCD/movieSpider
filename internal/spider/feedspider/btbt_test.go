package feedspider

import (
	"testing"
)

func TestNewBtbt(t *testing.T) {
	btbt := NewBtbt("*/5 * * * *", "https://www.1lou.me/forum-1.htm")
	videos, err := btbt.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, video := range videos {
		t.Log(video)
	}

}
