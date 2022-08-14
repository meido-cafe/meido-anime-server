package api

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/factory"
)

func InitRouter(router *gin.RouterGroup) {
	base := router.Group("/api")
	b := base.Group("v1")
	{
		r := b.Group("rss") // /api/v1/rss
		c := factory.NewRssApi()
		r.GET("info/mikan", c.GetMikanInfo)
		r.GET("search", c.GetSearch)
		r.GET("subject", c.GetSubject)
	}
}
