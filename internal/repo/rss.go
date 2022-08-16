package repo

import "database/sql"

type RssRepo struct {
	db *sql.DB
}

func NewRssRepo(db *sql.DB) RssInterface {
	return &RssRepo{db: db}
}
