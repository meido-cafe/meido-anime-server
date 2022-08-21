//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/service"
)

func NewVideoRepo() (ret repo.VideoInterface) {
	panic(wire.Build(
		NewDB,
		repo.NewVideoRepo,
	))
	return
}

func NewVideoService() (ret *service.VideoService) {
	panic(wire.Build(
		NewVideoRepo,
		service.NewVideoService,
	))
	return
}

func NewVideoApi() (ret *v1.VideoApi) {
	panic(wire.Build(
		NewVideoService,
		v1.NewVideoApi,
	))
	return
}
