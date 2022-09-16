//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"
	"meido-anime-server/internal/repo"
)

func NewQbittorrentClient() (ret repo.TorrentClientInterface) {
	panic(wire.Build(
		NewQB,
		repo.NewQbittorrent,
	))
	return
}
