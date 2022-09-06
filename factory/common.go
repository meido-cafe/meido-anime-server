//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"meido-anime-server/config"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/tool"
)

func NewDB() (ret *sqlx.DB) {
	panic(wire.Build(
		config.NewConfig,
		common.NewDB,
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
		config.NewConfig,
		common.NewQB,
	))
	return
}
