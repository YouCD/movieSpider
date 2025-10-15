package nameParser

import (
	"context"
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func TestNameParserModelHandler(t *testing.T) {
	typeStr, newName, year, resolution, err := NameParserModelHandler(context.Background(), "www.UIndex.org    -    Novocaine 2025 1080p BluRay x265 10bit Atmos TrueHD7 1-WiKi")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v,%#v,%#v,%#v", typeStr, newName, year, resolution)
}
