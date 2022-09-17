//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	v1 "meido-anime-server/internal/api/v1"
)

func NewUserApi() (ret *v1.UserApi) {
	panic(wire.Build(
		v1.NewUserApi,
	))
	return
}
