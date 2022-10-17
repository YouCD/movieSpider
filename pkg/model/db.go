package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"math/rand"
	"movieSpider/pkg"
	"movieSpider/pkg/config"
	"movieSpider/pkg/log"
	"movieSpider/pkg/types"
	"os"
	"strings"
	"sync"
	"time"
)

var MovieDB = new(movieDB)

type movieDB struct {
	db *sql.DB
}

var once sync.Once

func NewMovieDB() (*movieDB, error) {
	once.Do(func() {

		var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.Database) // 连接数据库
		mdb, err := sql.Open("mysql", dsn)                                                                                                                  // 不校验数据库信息，只对数据库信息做校验
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		err = mdb.Ping()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		MovieDB.db = mdb
	})
	return MovieDB, nil
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

//func queryMore() {
//
//	// 定义sql
//	sqlstr := `select id,namespace,name,podIP,hostIP,image,status,createTime,startTime,labels,kind,kindName from k8s_pod;`
//	rows, e := m.db.Query(sqlstr)
//
//	if e != nil {
//		return
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		//var data K8sPodTable
//		e := rows.Scan(&data.Id, &data.Namespace, &data.Name, &data.PodIP, &data.HostIP, &data.Image, &data.Status, &data.CreateTime, &data.StartTime, &data.Labels, &data.Kind, &data.KindName)
//		if e != nil {
//			return
//		}
//		fmt.Printf("%#v\n", data)
//	}
//	// 执行结果并获得结果
//	//row:=m.db.QueryRow(sqlstr)
//	//row.Scan(&data.Id, &data.Name, &data.Image,&data.CreateTime)
//	// 获取结果
//
//}

func (m *movieDB) CreatFeedVideo(video *types.FeedVideo) (err error) {
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
			return errors.WithMessage(pkg.ErrDBExist, video.Name)
		} else {
			return errors.WithMessage(err, video.Name)
		}
	}
	return
}
func (m *movieDB) CreatDouBanVideo(video *types.DouBanVideo) (err error) {
	v, err := m.FetchOneDouBanVideoByDouBanID(video.DoubanID)
	if err != nil {
		log.Warn(err)
	}
	if v != nil {
		// 将该记录变更为 可播放
		if v.Playable != video.Playable {
			v.Playable = video.Playable
			err = m.UpDateDouBanVideo(v)
			return errors.WithMessagef(err, "UpDateDouBanVideo%s", v.Names)
		}
		return nil
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
	return
}

func (m *movieDB) RandomOneDouBanVideo() (video *types.DouBanVideo, err error) {
	video = new(types.DouBanVideo)
	sql := `select id,names,douban_id,playable from douban_video where imdb_id="";`

	rows, err := m.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []*types.DouBanVideo
	for rows.Next() {
		var v types.DouBanVideo
		err = rows.Scan(&v.ID, &v.Names, &v.DoubanID, &v.Playable)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &v)
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(videos))
	video = videos[index]
	return
}

func (m *movieDB) FetchOneDouBanVideoByDouBanID(DouBanID string) (video *types.DouBanVideo, err error) {
	video = new(types.DouBanVideo)
	// 定义sql
	sql := `select id,names,douban_id,playable from douban_video where douban_id=? ;`
	row := m.db.QueryRow(sql, DouBanID)

	err = row.Scan(&video.ID, &video.Names, &video.DoubanID, &video.Playable)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errors.WithMessagef(pkg.ErrDBNotFound, "DouBanID: %s", DouBanID)
		}
		return nil, errors.WithMessagef(err, "DouBanID: %s", DouBanID)
	}
	return
}

func (m *movieDB) UpDateDouBanVideo(video *types.DouBanVideo) (err error) {
	// 定义sql
	sql := `update douban_video set imdb_id=?,row_data=?,playable=?,type=?,names=? where id=?;`
	_, err = m.db.Exec(sql, video.ImdbID, video.RowData, video.Playable, video.Type, video.Names, video.ID)
	if err != nil {
		return errors.WithMessage(err, video.Names)
	}
	return
}

//func (m *movieDB) FetchOneDouBanVideoByName(name string) (video *types.DouBanVideo, err error) {
//	video = new(types.DouBanVideo)
//	// 定义sql
//	sql := `select id,names,douban_id,playable from douban_video where names like=? ;`
//	row := m.db.QueryRow(sql, name)
//
//	err = row.Scan(&video.ID, &video.Names, &video.DoubanID, &video.Playable)
//	if err != nil {
//		if strings.Contains(err.Error(), "no rows in result set") {
//			return nil, errors.WithMessagef(pkg.ErrDBNotFound, "name: %s", name)
//		}
//		return nil, errors.WithMessagef(err, "name: %s", name)
//	}
//	return
//}

// FetchDouBanMovies 获取 所有的 电影名
func (m *movieDB) FetchDouBanMovies() (names []string, err error) {
	sql := `select names from douban_video where type="movie"`
	rows, err := m.db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	var namesA []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			continue
		}
		namesA = append(namesA, name)
	}

	for _, v := range namesA {
		var names1 []string
		if err = json.Unmarshal([]byte(v), &names1); err != nil {
			log.Error(err)
			continue
		}
		for _, n := range names1 {
			names = append(names, n)
		}

	}
	return
}

// FetchMagnetByName 通过电影名 获取磁力连接
func (m *movieDB) FetchMagnetByName(names []string) (videos []*types.FeedVideo, err error) {
	for _, n := range names {
		// 定义sql
		sql := `select id,magnet,name,torrent_name from feed_video where name like ? and  type="movie" and download!=1;`
		rows, err := m.db.Query(sql, fmt.Sprintf("%%.%s.%%", n))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		// 只查找 没有下载过 && 类型为movie数据
		log.Warn("movieDB: FetchMagnetByName 开始第一次查找数据.")
		for rows.Next() {
			var video types.FeedVideo
			err = rows.Scan(&video.ID, &video.Magnet, &video.Name, &video.TorrentName)
			if err != nil {
				return nil, err
			}
			videos = append(videos, &video)
		}

		// 如果 没有一条都没有找到 video
		if len(videos) == 0 {
			log.Warn("movieDB: FetchMagnetByName 开始第二次查找数据.")
			// 查找 没有下载过 && 类型不等于TV的数据
			sql := `select id,magnet,name,torrent_name from feed_video where name like ? and download!=1 and type!="tv";`

			rows, err := m.db.Query(sql, fmt.Sprintf("%%.%s.%%", n))
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var video types.FeedVideo
				err = rows.Scan(&video.ID, &video.Magnet, &video.Name, &video.TorrentName)
				if err != nil {
					return nil, err
				}
				videos = append(videos, &video)
			}
		}

	}

	return
}

func (m *movieDB) UpdateFeedVideoDownloadByID(id int32) (err error) {
	// 定义sql
	sql := `update feed_video set download=? where id=?;`
	_, err = m.db.Exec(sql, 1, id)
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