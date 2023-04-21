package feed

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGlodls_Crawler(t *testing.T) {
	f := &glodls{
		url: "https://glodls.to/rss.php?cat=1,41",
	}
	videos, err := f.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, v := range videos {
		bytes, _ := json.Marshal(v)
		//t.Log(v)
		fmt.Println(string(bytes))
	}
}
