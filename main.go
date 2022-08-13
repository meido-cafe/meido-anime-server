package main

import (
	"fmt"
	"meido-anime-server/etc"
	"meido-anime-server/internal/api"
	"meido-anime-server/internal/app"
)

func main() {
	etc.InitConfig()
	app.InitDB()
	engine, router := app.NewGin()
	api.InitRouter(router)

	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", etc.Conf.Server.Port)); err != nil {
		panic(err)
	}
}
