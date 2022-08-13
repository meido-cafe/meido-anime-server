//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/service"
)

func NewDemoRepo() (ret repo.IDemoRepo) {
	panic(wire.Build(
		NewDB,
		repo.NewDemoRepo,
	))
	return
}

func NewDemoService() (ret *service.DemoService) {
	panic(wire.Build(
		NewDemoRepo,
		service.NewDemoService,
	))
	return
}

func NewDemoHander() (ret *v1.DemoHandler) {
	panic(wire.Build(
		NewDemoService,
		v1.NewDemoHandler,
	))
	return
}
