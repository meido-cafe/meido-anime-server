//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/service"
)

func NewRssRepo() (ret repo.RssInterface) {
	panic(wire.Build(
		NewDB,
		repo.NewRssRepo,
	))
	return
}

func NewRssService() (ret *service.RssService) {
	panic(wire.Build(
		NewRssRepo,
		service.NewRssService,
	))
	return
}

func NewRssApi() (ret *v1.RssApi) {
	panic(wire.Build(
		NewRssService,
		v1.NewRssApi,
	))
	return
}
