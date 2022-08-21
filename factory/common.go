//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"meido-anime-server/etc"
	"meido-anime-server/internal/app"
)

func NewDB() (ret *sqlx.DB) {
	panic(wire.Build(
		etc.NewConfig,
		app.NewDB,
	))
	return
}
