package repo

import (
	"github.com/jmoiron/sqlx"
)

type RssRepo struct {
	db *sqlx.DB
}

func NewRssRepo(db *sqlx.DB) RssInterface {
	return &RssRepo{db: db}
}
