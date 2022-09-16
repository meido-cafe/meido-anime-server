//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/service"
)

func NewBangumiRepo() (ret *repo.BangumiRepo) {
	panic(wire.Build(
		NewBangumiClient,
		repo.NewBangumiRepo,
	))
	return
}

func NewBangumiService() (ret *service.BangumiService) {
	panic(wire.Build(
		NewBangumiRepo,
		service.NewBangumiService,
	))
	return
}

func NewBangumiApi() (ret *v1.BangumiApi) {
	panic(wire.Build(
		NewBangumiService,
		v1.NewBangumiApi,
	))
	return
}
