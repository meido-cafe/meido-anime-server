//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/service"
)

func NewRssService() (ret *service.RssService) {
	panic(wire.Build(
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
