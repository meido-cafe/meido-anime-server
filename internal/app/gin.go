package app

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/etc"
)

func NewGin() (*gin.Engine, *gin.RouterGroup) {
	gin.SetMode(etc.NewConfig().Server.GinMode)
	engine := gin.Default()
	router := engine.Group("")
	return engine, router
}
