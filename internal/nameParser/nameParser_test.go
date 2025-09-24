package nameParser

import (
	"movieSpider/internal/config"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
}
func TestNameParserModelHandler(t *testing.T) {
	handler, s, s2, s3, err := NameParserModelHandler("The Morning Show S04E02 The Revolution Will Be Televised XviD-AFG")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(handler, s, s2, s3)
}
