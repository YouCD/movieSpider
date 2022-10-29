package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"math/rand"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	types2 "movieSpider/internal/types"
	"os"
	"strings"
	"sync"
	"time"
)

type movieDB struct {
	db *sql.DB
}

var (
	once          sync.Once
	movieDatabase = new(movieDB)

	ErrorDataExist = errors.New("数据已存在")
)

func NewMovieDB() *movieDB {
	once.Do(func() {

		var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.Database) // 连接数据库
		mdb, err := sql.Open("mysql", dsn)                                                                                                                  // 不校验数据库信息，只对数据库信息做校验
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		err = mdb.Ping()
		if err != nil {
			if strings.Contains(err.Error(), "Unknown database") {
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port) // 连接数据库
				db, err := sql.Open("mysql", dsn)
				if err != nil {
					log.Error("fdsfsdfsdfsdf", err)
				}

				sql := "create database " + config.MySQL.Database + " charset utf8mb4;"
				_, err = db.Exec(sql)
				if err != nil {
					log.Error(err)
				}
				_, err = db.Exec("USE " + config.MySQL.Database)
				if err != nil {
					log.Error(err)
				}
				movieDatabase.db = mdb
				movieDatabase.InitDBTable()
				return
			}
			log.Error(err)
			os.Exit(1)

		}
		movieDatabase.db = mdb
	})
	return movieDatabase
}

func (m *movieDB) InitDBTable() (err error) {
	doubanVideoSQL := "CREATE TABLE `douban_video` (\n  `id` int(11) NOT NULL AUTO_INCREMENT,\n  `names` varchar(255) NOT NULL COMMENT '片名列表',\n  `douban_id` varchar(255) NOT NULL COMMENT '豆瓣ID',\n  `imdb_id` varchar(255) NOT NULL COMMENT 'imdbID',\n  `row_data` longtext NOT NULL COMMENT '原始数据',\n  `timestamp` bigint(11) NOT NULL COMMENT '修改创建时间',\n  `type` varchar(255) NOT NULL COMMENT '类型',\n  `playable` varchar(255) NOT NULL COMMENT '是否可以播放',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `name` (`names`)\n) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4;"
	_, err = m.db.Exec(doubanVideoSQL)
	if err != nil {
		return err
	}
	feedVideoSQL := "CREATE TABLE `feed_video` (\n  `id` int(11) NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) NOT NULL COMMENT '片名',\n  `torrent_name` varchar(255) NOT NULL COMMENT '种子名',\n  `torrent_url` varchar(255) NOT NULL COMMENT '种子引用地址',\n  `magnet` longtext NOT NULL COMMENT '磁力链接',\n  `year` varchar(255) NOT NULL COMMENT '年份',\n  `type` varchar(255) NOT NULL COMMENT 'tv或movie',\n  `row_data` longtext COMMENT '原始数据',\n  `web` varchar(255) NOT NULL COMMENT '站点',\n  `download` int(11) NOT NULL COMMENT '1:已经下载',\n  `timestamp` bigint(11) NOT NULL COMMENT '修改创建时间',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `name` (`name`,`torrent_name`)\n) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4;"
	_, err = m.db.Exec(feedVideoSQL)
	if err != nil {
		return err
	}
	return
}

//查询单条结果

func (m *movieDB) CreatFeedVideo(video *types2.FeedVideo) (err error) {
	if video.Magnet == "" {
		return errors.New(fmt.Sprintf("CreatFeedVideo Magnet is nill : %#v", video))
	}
	sql := `insert into feed_video(torrent_name,torrent_url,magnet,year,name,type,row_data,web,download,timestamp) value (?,?,?,?,?,?,?,?,?,?);`
	_, err = m.db.Exec(sql,
		video.TorrentName,
		video.TorrentUrl,
		video.Magnet,
		video.Year,
		video.Name,
		video.Type,
		video.RowData,
		video.Web,
		video.Download,
		time.Now().Unix())
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			log.Debugf("CreatFeedVideo 数据已存在 video: %#v", video)
			return errors.WithMessagef(ErrorDataExist, "name: %s type: %s.", video.Name, video.Type)
		}
		return errors.WithMessage(err, video.Name)
	}
	log.Debugf("CreatFeedVideo 数据已添加 video: %#v", video)
	return
}
func (m *movieDB) CreatDouBanVideo(video *types2.DouBanVideo) (err error) {
	v, err := m.FetchOneDouBanVideoByDouBanID(video.DoubanID)
	if err != nil {
		// 忽略 错误信息： sql: no rows in result set
		if !errors.Is(sql.ErrNoRows, err) {
			log.Error("video.DoubanID : %s,err: %s", video.DoubanID, err)
		}
	}
	if v != nil {
		log.Debugf("CreatDouBanVideo已存在 %#v", v)
		// 将该记录变更为 可播放
		if v.Playable != video.Playable {
			v.Playable = video.Playable
			log.Debugf("FetchOneDouBanVideoByDouBanID %#v", v)
			err = m.UpdateDouBanVideo(v)
			return errors.WithMessagef(err, "UpdateDouBanVideo %s", v.Names)
		}
		return nil
	}

	if video.Names == "null" {
		log.Errorf("CreatDouBanVideo 数据错误. video: %#v", video)
		return
	}

	sql := `insert into douban_video(names,douban_id,imdb_id,row_data,type,playable,timestamp) value (?,?,?,?,?,?,?);`
	_, err = m.db.Exec(sql,
		video.Names,
		video.DoubanID,
		video.ImdbID,
		video.RowData,
		video.Type,
		video.Playable,
		time.Now().Unix(),
	)

	if err != nil {
		return errors.WithMessage(err, video.Names)
	}
	log.Debugf("CreatDouBanVideo 数据已添加. video: %#v", video)
	return
}

func (m *movieDB) RandomOneDouBanVideo() (video *types2.DouBanVideo, err error) {
	video = new(types2.DouBanVideo)
	sql := `select id,names,douban_id,playable from douban_video where imdb_id="";`

	rows, err := m.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []*types2.DouBanVideo
	for rows.Next() {
		var v types2.DouBanVideo
		err = rows.Scan(&v.ID, &v.Names, &v.DoubanID, &v.Playable)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &v)
	}
	if len(videos) == 0 {
		return nil, errors.New("RandomOneDouBanVideo data is null")
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(videos))
	video = videos[index]
	log.Debugf("RandomOneDouBanVideo video: %#v", video)
	return
}

func (m *movieDB) FetchOneDouBanVideoByDouBanID(DouBanID string) (video *types2.DouBanVideo, err error) {
	video = new(types2.DouBanVideo)
	// 定义sql
	sql := `select id,names,douban_id,playable from douban_video where douban_id=? ;`
	row := m.db.QueryRow(sql, DouBanID)

	err = row.Scan(&video.ID, &video.Names, &video.DoubanID, &video.Playable)
	if err != nil {
		return nil, err
	}
	log.Debugf("FetchOneDouBanVideoByDouBanID video: %#v", video)
	return
}

func (m *movieDB) UpdateDouBanVideo(video *types2.DouBanVideo) (err error) {
	// 定义sql
	if video.Names == "" {
		return errors.New("空数据")
	}
	sql := `update douban_video set imdb_id=?,row_data=?,playable=?,type=?,names=? where id=?;`
	_, err = m.db.Exec(sql, video.ImdbID, video.RowData, video.Playable, video.Type, video.Names, video.ID)
	if err != nil {
		return errors.WithMessage(err, video.Names)
	}
	return
}

// FetchDouBanVideoByType 获取 所有的 电影名
func (m *movieDB) FetchDouBanVideoByType(typ types2.Resource) (nameList []string, err error) {
	log.Infof("FetchDouBanVideoByType 搜索 %s 类型豆瓣资源.", typ.Typ())
	sql := `select names from douban_video where type=?;`
	rows, err := m.db.Query(sql, typ.Typ())
	if err != nil {
		return
	}
	defer rows.Close()
	var namesA []string
	for rows.Next() {
		var names string
		if err = rows.Scan(&names); err != nil {
			continue
		}
		namesA = append(namesA, names)
	}

	for _, v := range namesA {
		var names1 []string
		if err = json.Unmarshal([]byte(v), &names1); err != nil {
			log.Error(err)
			continue
		}
		for _, n := range names1 {
			nameList = append(nameList, n)
		}

	}
	return
}

// FetchMovieMagnetByName 通过电影名 获取磁力连接
func (m *movieDB) FetchMovieMagnetByName(names []string) (videos []*types2.FeedVideo, err error) {

	var videos1 []*types2.FeedVideo
	log.Info("FetchMovieMagnetByName 开始第一次查找Movie数据.")
	log.Debugf("FetchMovieMagnetByName 开始第一次查找Movie数据: %s.", names)
	for _, n := range names {
		sql := `select id,magnet,name,torrent_name from feed_video where name like ? and magnet!="" and  type="movie" and download=0 ;`
		var likeName string
		likeName = fmt.Sprintf("%s%%", n)
		// 定义sql
		rows, err := m.db.Query(sql, likeName)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		// 只查找 没有下载过 && 类型为movie数据
		for rows.Next() {
			var video types2.FeedVideo
			err = rows.Scan(&video.ID, &video.Magnet, &video.Name, &video.TorrentName)
			if err != nil {
				return nil, err
			}
			videos1 = append(videos1, &video)
		}
	}
	if len(videos1) > 0 {
		return videos1, nil
	}

	log.Info("FetchMovieMagnetByName 开始第二次查找数据.")
	for _, n := range names {
		// 查找 没有下载过 && 类型不等于TV的数据
		sql := `select id,magnet,name,torrent_name from feed_video where name like ? and magnet!="" and download=0 and type!="tv";`
		var likeName string
		if strings.Contains(n, ".") {
			likeName = fmt.Sprintf("%%.%s.%%", n)
		} else {
			likeName = fmt.Sprintf("%%%s%%", n)
		}

		rows, err := m.db.Query(sql, likeName)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var video types2.FeedVideo
			err = rows.Scan(&video.ID, &video.Magnet, &video.Name, &video.TorrentName)
			if err != nil {
				return nil, err
			}
			videos = append(videos, &video)
		}
	}
	return
}

// FetchTVMagnetByName 通过 电视剧名 获取磁力连接
func (m *movieDB) FetchTVMagnetByName(names []string) (videos []*types2.FeedVideo, err error) {

	var videos1 []*types2.FeedVideo
	log.Info("FetchMovieMagnetByName 开始第一次查找tv数据.")
	log.Debugf("FetchMovieMagnetByName 开始第一次查找tv数据: %s.", names)
	for _, n := range names {
		sql := `select id,magnet,name,torrent_name from feed_video where name like ? and magnet!="" and  type="tv" and download=0;`
		var likeName string
		likeName = fmt.Sprintf("%s%%", n)
		// 定义sql
		rows, err := m.db.Query(sql, likeName)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		// 只查找 没有下载过 && 类型为movie数据
		for rows.Next() {
			var video types2.FeedVideo
			err = rows.Scan(&video.ID, &video.Magnet, &video.Name, &video.TorrentName)
			if err != nil {
				return nil, err
			}
			videos1 = append(videos1, &video)
		}
	}
	if len(videos1) > 0 {
		return videos1, nil
	}

	log.Info("FetchMovieMagnetByName 开始第二次查找tv数据.")
	for _, n := range names {
		// 查找 没有下载过 && 类型不等于TV的数据
		sql := `select id,magnet,name,torrent_name from feed_video where name like ? and magnet!="" and download=0 and type!="tv";`
		var likeName string
		if strings.Contains(n, ".") {
			likeName = fmt.Sprintf("%%.%s.%%", n)
		} else {
			likeName = fmt.Sprintf("%%%s%%", n)
		}

		rows, err := m.db.Query(sql, likeName)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var video types2.FeedVideo
			err = rows.Scan(&video.ID, &video.Magnet, &video.Name, &video.TorrentName)
			if err != nil {
				return nil, err
			}
			videos = append(videos, &video)
		}
	}
	return
}

func (m *movieDB) UpdateFeedVideoDownloadByID(id int32, isDownload int) (err error) {
	// 定义sql
	sql := `update feed_video set download=? where id=?;`
	_, err = m.db.Exec(sql, isDownload, id)
	if err != nil {
		return err
	}
	return
}

func (m *movieDB) CountFeedVideo() (counts []*types2.ReportCount, err error) {
	sql := `select count(*)  as count ,web  from  feed_video group by web order by count;`
	rows, err := m.db.Query(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := new(types2.ReportCount)
		err = rows.Scan(&c.Count, &c.Web)
		if err != nil {
			return nil, err
		}
		counts = append(counts, c)
	}
	return
}

func (m *movieDB) FindLikeTVFromFeedVideo(name string) (videos []*types2.FeedVideo, err error) {
	sql := `select id,name from feed_video where name like ?`
	rows, err := m.db.Query(sql, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		v := new(types2.FeedVideo)
		err = rows.Scan(&v.ID, &v.Name)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return
}
func (m *movieDB) UpdateFeedVideoNameByID(id int32, name string, resource types2.Resource) (err error) {
	// 定义sql
	sql := `update feed_video set name=?,type=? where id=?;`
	_, err = m.db.Exec(sql, name, resource.Typ(), id)
	if err != nil {
		return err
	}
	return
}

//func modifiy() {
//	sqlstr := `update k8s_pod set name=? where id=673;`
//	res, e := m.db.Exec(sqlstr, "fuckyou")
//	if e != nil {
//		return
//	}
//	n, e := res.RowsAffected()
//	if e != nil {
//		return
//	} else {
//		fmt.Printf("ID为%d", n)
//	}
//}
//
//func delete(id int) {
//	sqlstr := `delete from k8s_pod  where id=?;`
//	res, e := m.db.Exec(sqlstr, id)
//	if e != nil {
//		return
//	}
//	n, e := res.RowsAffected()
//	if e != nil {
//		return
//	} else {
//		fmt.Printf("ID为%d", n)
//	}
//}
