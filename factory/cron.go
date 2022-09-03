//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/service"
)

func NewCronRepo() (ret *repo.CronRepo) {
	panic(wire.Build(
		NewDB,
		repo.NewCronRepo,
	))
	return
}

func NewCronService() (ret *service.CronService) {
	panic(wire.Build(
		NewCronRepo,
		service.NewCronService,
	))
	return
}
