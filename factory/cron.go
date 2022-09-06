//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"meido-anime-server/internal/service"
)

func NewCronService() (ret *service.CronService) {
	panic(wire.Build(
		NewVideoService,
		service.NewCronService,
	))
	return
}
