package common

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"meido-anime-server/config"
	"meido-anime-server/internal/global"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var dbOnce sync.Once
var db *sqlx.DB

func InitSqlite() {
	dbOnce.Do(func() {
		var err error

		db, err = sqlx.Open("sqlite3", config.Conf.Db.Path)
		if err != nil {
			log.Fatalln("连接数据库失败:", err)
		}
		db.SetMaxOpenConns(config.Conf.Db.MaxCons)

		log.Println("正在初始化数据库")
		byt, err := ioutil.ReadFile("./sql/init.sql")
		if err != nil {
			log.Fatalln("初始化数据库失败,读取sql文件失败:", err)
		}

		if _, err = db.Exec(string(byt)); err != nil {
			log.Fatalln("初始化数据库失败:", err)
		}

		now := time.Now().Unix()
		initCategory := []string{"未分类", "anime", "movie", "ova"}
		stmt, err := db.Preparex(` insert or ignore into category (name,origin,create_time,update_time) values (?,?,?,?) `)
		if err != nil {
			log.Fatalln(err)
		}
		for _, item := range initCategory {
			if _, err = stmt.Exec(item, 1, now, now); err != nil {
				log.Fatalln("初始化数据库失败:", err)
			}

			if err = os.MkdirAll(filepath.Join(global.MediaPath, item), 0644); err != nil {
				log.Fatalln("创建分类目录失败:", err)
			}
		}

		log.Println("数据库初始化成功")
	})
}

func GetSqlite() *sqlx.DB {
	InitSqlite()
	return db
}

type DBClient interface {
	sqlx.Queryer
	sqlx.Execer
	sqlx.Preparer
	sqlx.QueryerContext
	sqlx.ExecerContext
	sqlx.PreparerContext

	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}
