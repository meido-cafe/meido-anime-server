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
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	config.InitConfig()

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

// TODO 订阅添加字段 是否正则
// TODO 删除番剧是否同时删除媒体资源 是否删除种子文件 是否同时删除订阅(只有订阅的才有该选项)
// TODO 历史番剧下载(只支持种子)
// TODO 系统设置表
