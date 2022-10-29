package download

import (
	"movieSpider/internal/config"
	"movieSpider/internal/model"
	"testing"
)

func Test_download_Run(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()
	d := &Download{
		scheduling: "The.Peripheral",
	}
	d.Run()
}

func Test_download_DownloadByName(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()
	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	d.DownloadByName("House.Of.The.Dragon", "1080")
}

func Test_download_DownloadByName1(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()
	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	d.DownloadByName("Werewolf.by.Night", "2160")
}

func Test_download_downloadTvTask(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	model.NewMovieDB()

	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	err := d.downloadTvTask()
	if err != nil {
		t.Error(err)
	}

}

func TestDownload_downloadMovieTask(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpiderCore/bin/movieSpiderCore/config.yaml")
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")

	model.NewMovieDB()

	d := &Download{
		scheduling: "tt.fields.scheduling",
	}
	err := d.downloadMovieTask()
	if err != nil {
		t.Error(err)
	}

}
