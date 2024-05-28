package bot

import (
	"movieSpider/internal/aria2"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/types"
	"testing"
)

func init() {
	config.InitConfig("/home/ycd/self_data/source_code/go-source/tools-cmd/movieSpider/config.local.yaml")
	model.NewMovieDB()
}
func TestTGBot_SendDatePublishedMsg(t1 *testing.T) {
	t := NewTgBot(config.Config.TG.BotToken, config.Config.TG.TgIDs)
	go t.StartBot()
	obj := &types.DouBanVideo{
		ID:            99119,
		Names:         `["阿凡达3：带种者"]`,
		DoubanID:      "878",
		ImdbID:        "444444444",
		RowData:       `{"name":"死侍3 Deadpool 3","url":"/subject/26957900/","image":"https://img9.doubanio.com/view/photo/s_ratio_poster/public/p2881030676.jpg","director":[{"type":"Person","url":"/celebrity/1289935/","name":"肖恩·利维 Shawn Levy"}],"author":[{"type":"Person","url":"/celebrity/1028169/","name":"罗伯·莱菲尔德 Rob Liefeld"},{"type":"Person","url":"/celebrity/1481720/","name":"温蒂·莫利纽兹 Wendy Molyneux"},{"type":"Person","url":"/celebrity/1481718/","name":"莉兹·莫利纽兹-罗格林 Lizzie Molyneux-Logelin"},{"type":"Person","url":"/celebrity/1014682/","name":"法比安·尼切扎 Fabian Nicieza"},{"type":"Person","url":"/celebrity/1327382/","name":"瑞特·里斯 Rhett Reese"},{"type":"Person","url":"/celebrity/1053623/","name":"瑞安·雷诺兹 Ryan Reynolds"},{"type":"Person","url":"/celebrity/1327381/","name":"保罗·韦尼克 Paul Wernick"},{"type":"Person","url":"/celebrity/1355474/","name":"扎布·威尔斯 Zeb Wells"}],"actor":[{"type":"Person","url":"/celebrity/1053623/","name":"瑞安·雷诺兹 Ryan Reynolds"},{"type":"Person","url":"/celebrity/1010508/","name":"休·杰克曼 Hugh Jackman"},{"type":"Person","url":"/celebrity/1414676/","name":"艾玛·科林 Emma Corrin"},{"type":"Person","url":"/celebrity/1006964/","name":"欧文·威尔逊 Owen Wilson"},{"type":"Person","url":"/celebrity/1017908/","name":"莫蕾娜·巴卡琳 Morena Baccarin"},{"type":"Person","url":"/celebrity/1025138/","name":"马修·麦克费登 Matthew Macfadyen"},{"type":"Person","url":"/celebrity/1018349/","name":"莱斯利·格塞斯 Leslie Uggams"},{"type":"Person","url":"/celebrity/1268159/","name":"罗伯·德兰尼 Rob Delaney"},{"type":"Person","url":"/celebrity/1329261/","name":"卡兰·索尼 Karan Soni"},{"type":"Person","url":"/celebrity/1351117/","name":"布里安娜·希德布兰德 Brianna Hildebrand"},{"type":"Person","url":"/celebrity/1275349/","name":"忽那汐里 Kutsuna Shiori"},{"type":"Person","url":"/celebrity/1190421/","name":"斯蒂芬·卡皮契奇 Stefan Kapičić"},{"type":"Person","url":"/celebrity/1010540/","name":"帕特里克·斯图尔特 Patrick Stewart"}],"releaseTimeJob":"2024-11-08","genre":["喜剧","动作","科幻"],"duration":"","description":"瑞安·雷诺兹主演的《死侍3》官宣休·杰克曼惊喜回归，继续饰演金刚狼。","type":"Movie","aggregateRating":{"type":"AggregateRating","ratingCount":"0","bestRating":"10","worstRating":"2","ratingValue":""}}`,
		Timestamp:     0,
		Type:          "444444444",
		Playable:      "三21问问是岁",
		DatePublished: "2023-06-13",
	}
	newAria2, _ := aria2.NewAria2(config.Downloader.Aria2Label)

	newAria2.AddDownloadTask(obj, "6378562f5e923563")
	downLoadChan := make(chan *types.DownloadNotifyVideo)
	defer close(downLoadChan)
	newAria2.Subscribe(downLoadChan)
	log.Errorf("%v", video)
	select {}
}
