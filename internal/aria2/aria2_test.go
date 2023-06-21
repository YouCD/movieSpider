package aria2

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"testing"
)

var newAria2 = &aria2{}

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	newAria2, _ = NewAria2(config.Downloader.Aria2Label)

}
func downloadCompleteNotify() {

	for {
		subscribeCh := newAria2.Subscribe()
		select {
		case completedDownload, ok := <-subscribeCh:
			if ok {
				fmt.Println("ccc", newAria2)
				// 处理已完成的下载任务
				fmt.Println("Received completed video:", completedDownload)
			} else {
				continue
			}

		}
	}
}

func Test_aria2_DownloadList(t *testing.T) {
	newAria2, err := NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		t.Error(err)
	}
	//v := &types.DouBanVideo{
	//	ID:        99119,
	//	Names:     `["阿凡达3：带种者111","Avatar:The.Seed.Bearer11"]`,
	//	DoubanID:  "878",
	//	ImdbID:    "444444444",
	//	RowData:   "444444444",
	//	Timestamp: 0,
	//	Type:      "444444444",
	//	Playable:  "三21问问是岁",
	//}
	//
	//newAria2.AddDownloadTask(v, "b26721bb7071c8c6")
	////newAria2.AddDownloadTask(v, "89e9fc2f333d2830")
	//subscribeCh := newAria2.Subscribe()
	//for complete := range subscribeCh {
	//	fmt.Println(complete)
	//}
	//go downloadCompleteNotify()
	//
	//tick := time.Tick(time.Second * 10)
	//
	//for {
	//	//subscribeCh := newAria2.Subscribe()
	//
	//	select {
	//	case <-tick:
	//		fmt.Println("添加", "c611bb47ac6673eb")
	//		newAria2.AddDownloadTask(v, "c611bb47ac6673eb")
	//		//case v := <-subscribeCh:
	//		//	fmt.Println("subscribe", v)
	//	}
	//}
	//url := `magnet:?xt=urn:btih:0ceaa977f733050a60c0164488f70fdad14ac4d9&dn=Extraction.2.2023.2160p.NF.WEB-DL.DDP5.1.Atmos.DV.H.265-FLUX&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2740%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.pirateparty.gr%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentor.org%3A2710%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.si%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.tracker.cl%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2740%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2790%2Fannounce`
	urls := []string{
		`magnet:?xt=urn:btih:BF6266404D800D2ACFFB143708DD3A9E93C1D938&dn=The.Flash.2014.S09E11.1080p.HEVC.x265-MeGusta%5Beztv.re%5D.mkv&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftorrent.gresille.org%3A80%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2710%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337`}
	for _, url := range urls {

		gid, err := newAria2.DownloadByMagnet(url)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("gidgidgidgid             ", gid)

	}

}

func Test_aria2_CompletedFiles(t *testing.T) {

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
	}

	log.Infof("Report: 下载统计: %s", msg)
	//fmt.Println(bs)

}

func Test_aria2_getAllActiveGID(t *testing.T) {
	// 03b5a009333b8158 种子
	// 90d7c6958629d5a0 磁力
	info, err := newAria2.aria2Client.TellStatus("187a69be1bfeec9e", "files", "gid", "status", "errorMessage", "belongsTo", "following", "followedBy")
	if err != nil {
		t.Error(err)

	}
	fmt.Println(info.Files)
	fmt.Println(info.Status)

	fmt.Println(info.FollowedBy)
}

func Test_aria2_DownloadByMagnet(t *testing.T) {
	urls := []string{
		`magnet:?xt=urn:btih:AEBCD2368F6E5C992B1332B9164AE53B6AF85553&dn=The.Flash.2014.S09E01.1080p.WEB.H264-GGWP%5BTGx%5D&tr=udp://open.stealth.si:80/announce&tr=udp://tracker.tiny-vps.com:6969/announce&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://tracker.torrent.eu.org:451/announce&tr=udp://explodie.org:6969/announce&tr=udp://tracker.cyberia.is:6969/announce&tr=udp://ipv4.tracker.harry.lu:80/announce&tr=udp://p4p.arenabg.com:1337/announce&tr=udp://tracker.birkenwald.de:6969/announce&tr=udp://tracker.moeking.me:6969/announce&tr=udp://opentor.org:2710/announce&tr=udp://tracker.dler.org:6969/announce&tr=udp://9.rarbg.me:2970/announce&tr=https://tracker.foreverpirates.co:443/announce&tr=http://vps02.net.orel.ru:80/announce`,
	}
	for _, url := range urls {

		gid, err := newAria2.DownloadByMagnet(url)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("gidgidgidgid             ", gid)

	}
}
