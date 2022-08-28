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
	{
		r := b.Group("video") // /api/v1/video
		c := factory.NewVideoApi()
		r.POST("subscribe", c.Subscribe)      // 订阅番剧
		r.GET("detail", c.GetOne)             // 获取详情
		r.GET("list", c.GetList)              // 获取番剧列表
		r.DELETE("rss", c.DeleteRss)          // 删除rss链接
		r.PUT("rss", c.UpdateRss)             // 更新rss链接
		r.GET("qbittorrent/log", c.GetQBLogs) // 获取QB日志
	}
}
