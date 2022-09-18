package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"meido-anime-server/config"
	"meido-anime-server/factory"
	"meido-anime-server/internal/api"
)

func main() {
	config.InitConfig()
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	gin.SetMode(config.NewConfig().Server.GinMode)
	engine := gin.Default()
	router := engine.Group("")
	api.InitRouter(router)

	init := factory.NewInitService()
	init.Init()

	if err := engine.Run(fmt.Sprintf("%s:%d", "localhost", config.Conf.Server.Port)); err != nil {
		log.Fatalln("启动失败:", err)
	}
}

// TODO 历史番剧下载(只支持种子)
// TODO 系统设置表
