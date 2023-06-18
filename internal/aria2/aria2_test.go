package aria2

import (
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"testing"
)

func downloadCompleteNotify() {
	aria2Server, err := NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		subscribeCh := aria2Server.Subscribe()
		select {
		case completedDownload, ok := <-subscribeCh:
			if ok {
				fmt.Println("ccc", aria2Server)
				// 处理已完成的下载任务
				fmt.Println("Received completed video:", completedDownload)
			} else {
				continue
			}

		}
	}
}

func Test_aria2_DownloadList(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/go-source/tools-cmd/core/bin/core/config.yaml")
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")

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
	url := `magnet:?xt=urn:btih:94dc34d89787110d0224936f63d449114becb1c9&dn=%E3%80%90%E9%A6%96%E5%8F%91%E4%BA%8E%E9%AB%98%E6%B8%85%E5%BD%B1%E8%A7%86%E4%B9%8B%E5%AE%B6+www.BBQDDQ.com%E3%80%91%E6%83%8A%E5%A4%A9%E8%90%A5%E6%95%91%5B%E7%AE%80%E7%B9%81%E8%8B%B1%E5%AD%97%E5%B9%95%5D.Extraction.2020.2160p.NF.WEB-DL.DDP.5.1.Atmos.H.265-DreamHD&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=wss%3A%2F%2Ftracker.openwebtorrent.com%3A443%2Fannounce&tr=http%3A%2F%2F107.189.10.20.sslip.io%3A7777%2Fannounce&tr=http%3A%2F%2F1337.abcvg.info%3A80%2Fannounce&tr=http%3A%2F%2Fbt.endpot.com%3A80%2Fannounce&tr=http%3A%2F%2Fi-p-v-6.tk%3A6969%2Fannounce&tr=http%3A%2F%2Fipv6.1337.cx%3A6969%2Fannounce&tr=http%3A%2F%2Fipv6.govt.hu%3A6969%2Fannounce&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=http%3A%2F%2Fopen-v6.demonoid.ch%3A6969%2Fannounce&tr=http%3A%2F%2Fopen.acgnxtracker.com%3A80%2Fannounce&tr=http%3A%2F%2Fopen.tracker.ink%3A6969%2Fannounce&tr=http%3A%2F%2Fp2p.0g.cx%3A6969%2Fannounce&tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce&tr=http%3A%2F%2Ft.nyaatracker.com%3A80%2Fannounce&tr=http%3A%2F%2Ftorrentsmd.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker.bt4g.com%3A2095%2Fannounce&tr=http%3A%2F%2Ftracker.files.fm%3A6969%2Fannounce&tr=http%3A%2F%2Ftracker.gbitt.info%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.ipv6tracker.ru%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.k.vu%3A6969%2Fannounce&tr=http%3A%2F%2Ftrackme.theom.nz%3A80%2Fannounce&tr=http%3A%2F%2Fv6-tracker.0g.cx%3A6969%2Fannounce&tr=http%3A%2F%2Fwww.all4nothin.net%3A80%2Fannounce.php&tr=http%3A%2F%2Fwww.wareztorrent.com%3A80%2Fannounce&tr=https%3A%2F%2F1337.abcvg.info%3A443%2Fannounce&tr=https%3A%2F%2Fopentracker.i2p.rocks%3A443%2Fannounce&tr=https%3A%2F%2Ft1.hloli.org%3A443%2Fannounce&tr=https%3A%2F%2Ftr.abiir.top%3A443%2Fannounce&tr=https%3A%2F%2Ftr.abir.ga%3A443%2Fannounce&tr=https%3A%2F%2Ftr.burnabyhighstar.com%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.foreverpirates.co%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.gbitt.info%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.imgoingto.icu%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.kuroy.me%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.lilithraws.cf%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.lilithraws.org%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.loligirl.cn%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.mlsub.net%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.nanoha.org%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.tamersunion.org%3A443%2Fannounce&tr=https%3A%2F%2Ftracker1.520.jp%3A443%2Fannounce&tr=https%3A%2F%2Ftrackme.theom.nz%3A443%2Fannounce&tr=udp%3A%2F%2F6ahddutb1ucc3cp.ru%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.com%3A2810%2Fannounce&tr=udp%3A%2F%2Faarsen.me%3A6969%2Fannounce&tr=udp%3A%2F%2Facxx.de%3A6969%2Fannounce&tr=udp%3A%2F%2Faegir.sexy%3A6969%2Fannounce&tr=udp%3A%2F%2Fastrr.ru%3A6969%2Fannounce&tr=udp%3A%2F%2Fbt.ktrackers.com%3A6666%2Fannounce&tr=udp%3A%2F%2Fcutscloud.duckdns.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fdht.bt251.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ffe.dealclub.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ffree.publictracker.xyz%3A6969%2Fannounce&tr=udp%3A%2F%2Fhtz3.noho.st%3A6969%2Fannounce&tr=udp%3A%2F%2Fipv4.tracker.harry.lu%3A80%2Fannounce&tr=udp%3A%2F%2Fipv6.69.mu%3A6969%2Fannounce&tr=udp%3A%2F%2Fipv6.tracker.monitorit4.me%3A6969%2Fannounce&tr=udp%3A%2F%2Flaze.cc%3A6969%2Fannounce&tr=udp%3A%2F%2Fmail.artixlinux.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fmirror.aptus.co.tz%3A6969%2Fannounce&tr=udp%3A%2F%2Fmoonburrow.club%3A6969%2Fannounce&tr=udp%3A%2F%2Fmovies.zsw.ca%3A6969%2Fannounce&tr=udp%3A%2F%2Fnew-line.net%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.4ever.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.com%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.dstud.io%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.free-tracker.ga%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.publictracker.xyz%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.tracker.ink%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentor.org%3A2710%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Fprivate.anonseed.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fpsyco.fr%3A6969%2Fannounce&tr=udp%3A%2F%2Fpublic-tracker.ml%3A6969%2Fannounce&tr=udp%3A%2F%2Frep-art.ynh.fr%3A6969%2Fannounce&tr=udp%3A%2F%2Fstatic.54.161.216.95.clients.your-server.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ft.133335.xyz%3A6969%2Fannounce&tr=udp%3A%2F%2Fthagoat.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Fthetracker.org%3A80%2Fannounce&tr=udp%3A%2F%2Ftorrents.artixlinux.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.4.babico.name.tr%3A3131%2Fannounce&tr=udp%3A%2F%2Ftracker.altrosky.nl%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.artixlinux.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.auctor.tv%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.beeimg.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.bitsearch.to%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.jonaslsa.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.leech.ie%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.lelux.fi%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.monitorit4.me%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.openbtba.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.pimpmyworld.to%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.skynetcloud.site%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.skyts.net%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tcp.exchange%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.theoks.net%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker1.bt.moack.co.kr%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker1.myporn.club%3A9337%2Fannounce&tr=udp%3A%2F%2Ftracker2.dler.com%3A80%2Fannounce&tr=udp%3A%2F%2Ftrackerb.jonaslsa.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fuploads.gamecoast.net%3A6969%2Fannounce&tr=udp%3A%2F%2Fv1046920.hosted-by-vdsina.ru%3A6969%2Fannounce&tr=udp%3A%2F%2Fwww.peckservers.com%3A9000%2Fannounce&tr=udp%3A%2F%2Fzecircle.xyz%3A6969%2Fannounce&tr=ws%3A%2F%2Fhub.bugout.link%3A80%2Fannounce&ws=http%3A%2F%2F130.61.51.241%3A11111%2F&ws=http%3A%2F%2F130.61.51.241%3A11112%2F&ws=http%3A%2F%2F130.61.51.241%3A11113%2F&ws=http%3A%2F%2F130.61.82.9%3A11111%2F&ws=http%3A%2F%2F130.61.82.9%3A11112%2F&ws=http%3A%2F%2F130.61.82.9%3A11113%2F&ws=http%3A%2F%2F130.162.254.83%3A11111%2F&ws=http%3A%2F%2F130.162.254.83%3A11112%2F&ws=http%3A%2F%2F130.162.254.83%3A11113%2F&ws=http%3A%2F%2F130.61.110.156%3A11111%2F&ws=http%3A%2F%2F130.61.110.156%3A11112%2F&ws=http%3A%2F%2F130.61.110.156%3A11113%2F&ws=http%3A%2F%2F130.162.224.45%3A11111%2F&ws=http%3A%2F%2F130.162.224.45%3A11112%2F&ws=http%3A%2F%2F130.162.224.45%3A11113%2F`
	gid, err := newAria2.DownloadByMagnet(url)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(gid)

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
