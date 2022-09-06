package common

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"meido-anime-server/config"
	"sync"
)

var db *sqlx.DB
var dbOnce sync.Once

func NewDB(conf *config.Config) *sqlx.DB {
	dbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("sqlite3", conf.Db.Path)
		if err != nil {
			log.Fatalln("连接数据库失败:", err)
		}
		db.SetMaxOpenConns(conf.Db.MaxCons)

		log.Println("正在初始化数据库")
		byt, err := ioutil.ReadFile("./sql/init.sql")
		if err != nil {
			log.Fatalln("初始化数据库失败,读取sql文件失败:", err)
		}

		if _, err = db.Exec(string(byt)); err != nil {
			log.Fatalln("初始化数据库失败:", err)
		}

		log.Println("数据库初始化成功")
	})
	return db
}
