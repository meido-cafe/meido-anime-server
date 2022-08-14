package repo

import "database/sql"

type Rss struct {
	db *sql.DB
}

func NewRssRepo(db *sql.DB) RssInterface {
	return &Rss{db: db}
}
