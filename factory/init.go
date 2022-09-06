//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"meido-anime-server/internal/service"
)

func NewInitService() (ret *service.InitService) {
	panic(wire.Build(
		NewCronService,
		NewVideoService,
		service.NewInitService,
	))
	return
}
