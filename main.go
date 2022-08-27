package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"meido-anime-server/etc"
	"meido-anime-server/internal/api"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	etc.InitConfig()

	gin.SetMode(etc.NewConfig().Server.GinMode)
	engine := gin.Default()
	router := engine.Group("")

	api.InitRouter(router)

	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", etc.Conf.Server.Port)); err != nil {
		panic(err)
	}
}
