//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/service"
)

func NewUserService() *service.UserService {
	panic(wire.Build(
		service.NewUserService,
	))
}

func NewUserApi() (ret *v1.UserApi) {
	panic(wire.Build(
		NewUserService,
		v1.NewUserApi,
	))
	return
}
