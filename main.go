package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"meido-anime-server/config"
	"meido-anime-server/internal/api"
	"meido-anime-server/internal/service"
)

func main() {
	config.InitConfig()

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	gin.SetMode(config.GetConfig().Server.GinMode)
	engine := gin.Default()
	router := engine.Group("")
	api.InitRouter(router)

	init := service.NewInitService()
	init.Init()

	if err := engine.Run(fmt.Sprintf("%s:%d", "0.0.0.0", config.Conf.Server.Port)); err != nil {
		log.Fatalln("启动失败:", err)
	}
}

// TODO 历史番剧下载(只支持种子)
// TODO 系统设置表
