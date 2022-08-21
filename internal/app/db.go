package app

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"meido-anime-server/etc"
	"meido-anime-server/global"
)

func NewDB(conf *etc.Config) *sqlx.DB {
	global.DBOnce.Do(func() {
		var err error
		global.DB, err = sqlx.Open("sqlite3", conf.Db.Path)
		if err != nil {
			panic(err)
			return
		}
		global.DB.SetMaxOpenConns(conf.Db.MaxCons)
	})
	return global.DB
}

func InitDB() {
	NewDB(etc.NewConfig())
}
