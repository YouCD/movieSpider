package nameParser

import (
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func TestNameParserModelHandler(t *testing.T) {
	handler, s, s2, s3, err := NameParserModelHandler("American.Psycho.2000.Remastered.1080p.BluRay.X264.AC3.Wi")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(handler, s, s2, s3)
}
