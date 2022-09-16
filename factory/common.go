//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"github.com/imroc/req/v3"
	"github.com/jmoiron/sqlx"
	"meido-anime-server/config"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/tool"
)

func NewConfig() (ret *config.Config) {
	panic(wire.Build(
		config.NewConfig,
	))
	return
}

func NewDB() (ret *sqlx.DB) {
	panic(wire.Build(
		NewConfig,
		common.NewDB,
	))
	return
}

func NewDBClient() (ret common.DBClient) {
	panic(wire.Build(
		NewDB,
		common.NewDBClient,
	))
	return
}

func NewSqlTool() (ret *tool.Sql) {
	panic(wire.Build(
		tool.NewSql,
	))
	return
}

func NewQB() (ret *common.QB) {
	panic(wire.Build(
		NewConfig,
		common.NewQB,
	))
	return
}

func NewBangumiClient() (ret *req.Client) {
	panic(wire.Build(
		common.NewBangumiClient,
	))
	return
}
