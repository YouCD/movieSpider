package download

import (
	"context"
	"github.com/youcd/toolkit/log"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB()
}
func Test_download_Run(t *testing.T) {

	d := &Download{
		scheduling: "The.Peripheral",
	}
	d.Run()
}

func Test_download_DownloadByName(t *testing.T) {

	d := &Download{
		scheduling: "*/1 * * * *",
	}
	log.SetLogLevel("Debug")
	err := d.download(types.VideoTypeTV, model.NewMovieDB().GetFeedVideoTVByNames)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}

	//newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	//if err != nil {
	//	t.Error(err)
	//}
	//downLoadChan := make(chan *types.DownloadNotifyVideo)
	//defer close(downLoadChan)
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 1)
	//		newAria2.Subscribe(downLoadChan)
	//		select {
	//		case v, ok := <-downLoadChan:
	//			if ok {
	//				fmt.Println("subscribe", v)
	//			}
	//		}
	//	}
	//}()

}
