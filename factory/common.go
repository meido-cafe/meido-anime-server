//go:build wireinject
// +build wireinject

package factory

import (
	"database/sql"
	"github.com/google/wire"
	"meido-anime-server/etc"
	"meido-anime-server/internal/app"
)

func NewDB() (ret *sql.DB) {
	panic(wire.Build(
		etc.NewConfig,
		app.NewDB,
	))
	return
}
