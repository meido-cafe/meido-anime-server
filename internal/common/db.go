package common

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"meido-anime-server/etc"
	"sync"
)

var db *sqlx.DB
var dbOnce sync.Once

func NewDB(conf *etc.Config) *sqlx.DB {
	dbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("sqlite3", conf.Db.Path)
		if err != nil {
			panic(err)
			return
		}
		db.SetMaxOpenConns(conf.Db.MaxCons)

		log.Println("正在初始化数据库")
		byt, err := ioutil.ReadFile("./sql/init.sql")
		if err != nil {
			panic(err)
		}
		log.Println("数据库初始化成功")

		if _, err = db.Exec(string(byt)); err != nil {
			panic(err)
		}
	})
	return db
}
