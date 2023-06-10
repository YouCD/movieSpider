package download

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/model"
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
		scheduling: "tt.fields.scheduling",
	}
	d.DownloadByName("House.Of.The.Dragon", "1080")
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
	select {}
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
