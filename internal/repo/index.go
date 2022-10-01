package repo

import (
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/tool"
)

type Repo struct {
	db  common.DBClient
	sql *tool.Sql
}

func NewRepo(db common.DBClient) *Repo {
	return &Repo{
		db:  db,
		sql: tool.NewSql(),
	}
}
