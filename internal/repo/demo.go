package repo

import (
	"context"
	"database/sql"
	"log"
)

type IDemoRepo interface {
	Hello(ctx context.Context) (err error)
}

func NewDemoRepo(db *sql.DB) IDemoRepo {
	return &DemoRepo{db: db}
}

type DemoRepo struct {
	db *sql.DB
}

func (d DemoRepo) Hello(ctx context.Context) (err error) {
	log.Println("ok")
	return nil
}
