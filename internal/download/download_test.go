package download

import (
	"context"
	"testing"

	"github.com/youcd/toolkit/log"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
)

func init() {
	config.InitConfig("../../config.local.yaml")
	model.NewMovieDB()
}

func Test_download_Run(t *testing.T) {
	d := &Download{
		scheduling: "The.Peripheral",
	}
	d.Run(context.Background())
}

func Test_download_DownloadByName(t *testing.T) {
	d := &Download{
		scheduling: "*/1 * * * *",
	}
	log.SetLogLevel("Debug")
	err := d.download(context.Background(), types.VideoTypeTV, model.NewMovieDB().GetFeedVideoTVByNames)
	if err != nil {
		log.WithCtx(context.Background()).Error(err)
	}

	//newAria2, err := aria2.NewAria2(config.Config().Downloader.Aria2Label)
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
