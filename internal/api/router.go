package api

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/factory"
)

func InitRouter(router *gin.RouterGroup) {
	base := router.Group("/api")
	b := base.Group("v1")
	{
		r := b.Group("demo")
		c := factory.NewDemoHander()
		r.GET("hello", c.Hello)
	}
}
