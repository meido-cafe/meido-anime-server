package main

import (
	"fmt"
	"log"
	"meido-anime-server/etc"
	"meido-anime-server/internal/api"
	"meido-anime-server/internal/app"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	etc.InitConfig()
	app.InitDB()
	engine, router := app.NewGin()
	api.InitRouter(router)

	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", etc.Conf.Server.Port)); err != nil {
		panic(err)
	}
}
