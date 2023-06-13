package aria2

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"path"
	"strconv"
	"strings"
	"testing"
)

func Test_aria2_DownloadList(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/go-source/tools-cmd/core/bin/core/config.yaml")
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/core/bin/core/config.yaml")

	newAria2, err := NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		t.Error(err)
	}
	//info, err := newAria2.aria2Client.GetGlobalStat()
	//if err != nil {
	//	t.Error(err)
	//}
	//sessionInfo, err := newAria2.aria2Client.TellActive()
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//for _, v := range sessionInfo {
	//	fmt.Println(v.TotalLength)
	//	fmt.Println(v.Gid)
	//	//fmt.Printf("%#v\n", v)
	//	infos, err := newAria2.aria2Client.GetServers(v.Gid)
	//
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	fmt.Println(infos)
	//	//for _, f := range v.Files {
	//	//	fmt.Printf("%#v\n", f)
	//	//
	//	//}
	//}
	sessionInfo, err := newAria2.aria2Client.TellStopped(0, 100)
	if err != nil {
		t.Error(err)
	}

	for _, v := range sessionInfo {
		//marshal, _ := json.Marshal(v.Files)
		//fmt.Println("文件", string(marshal))
		if len(v.Files) > 0 {

			if strings.Contains(v.Files[0].Path, "[METADATA]") {
				continue
			} else {
				fmt.Println("GID:", v.Gid)

				CompletedLength, err := strconv.Atoi(v.Files[0].CompletedLength)
				if err != nil {
					log.Error(err)
				}
				Length, err := strconv.Atoi(v.Files[0].Length)
				if err != nil {
					log.Error(err)
				}
				fmt.Printf("文件大小: %d\n", Length)

				fmt.Printf("文件大小: %.2fGB\n", float32(Length)/1024/1024/1024)
				fmt.Printf("文件完成度百分比: %d%%\n", CompletedLength/Length*100)
				_, file := path.Split(v.Files[0].Path)
				fmt.Println("文件:", file)
				fmt.Println("-------------")

			}
		}
		//fmt.Printf("%#v\n", v)

		//for _, f := range v.Files {
		//	fmt.Printf("%#v\n", f)
		//
		//}
	}

}

func Test_aria2_CompletedFiles(t *testing.T) {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")

	newAria2, err := NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		t.Error(err)
	}
	files := newAria2.CurrentActiveAndStopFiles()
	//var s string
	//var bs string
	//for _, file := range files {
	//	if utf8.RuneCountInString(file.FileName) > 40 {
	//		nameRune := []rune(file.FileName)
	//		bs += fmt.Sprintf("%-40s | %s\n", string(nameRune[0:40]), file.Completed)
	//	} else {
	//		bs += fmt.Sprintf("%-40s | %s\n", file.FileName, file.Completed)
	//	}
	//}

	var msg string
	for _, file := range files {
		msg += fmt.Sprintf("\nGID:%s, 大小:%s, 已完成:%s, 文件名:%s", file.GID, file.Size, file.Completed, file.FileName)
		// todo 下载完后的向TG通知
	}

	log.Infof("Report: 下载统计: %s", msg)
	//fmt.Println(bs)

}
