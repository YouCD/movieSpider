package download

import (
	"fmt"
	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
	"time"
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

	newAria2, err := aria2.NewAria2(config.Config.Downloader.Aria2Label)
	if err != nil {
		t.Error(err)
	}
	downLoadChan := make(chan *types.DownloadNotifyVideo)
	go func() {
		for {
			time.Sleep(time.Second * 1)
			newAria2.Subscribe(downLoadChan)
			select {
			case v, ok := <-downLoadChan:
				if ok {
					fmt.Println("subscribe", v)
				}
			}
		}
	}()
	d.downloadTask()
	select {}
}

func Test_download_DownloadByName1(t *testing.T) {
	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	d.DownloadByName("Werewolf.by.Night", "2160")
}

func Test_download_downloadTvTask(t *testing.T) {

	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	err := d.downloadTvTask()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("done")
	//select {}
}

func TestDownload_downloadMovieTask(t *testing.T) {
	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	err := d.downloadMovieTask()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("done")
	select {}
}
