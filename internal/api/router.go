package api

import (
	"github.com/gin-gonic/gin"
	v1 "meido-anime-server/internal/api/v1"
	"meido-anime-server/internal/middleware"
)

func InitRouter(router *gin.RouterGroup) {
	api := v1.NewApi()
	midWare := middleware.NewMiddleware() // 中间件

	base := router.Group("/api")
	b := base.Group("v1")
	{
		r := b.Group("user")
		r.POST("login", api.Login)                  // 登录
		r.GET("logout", midWare.Auth(), api.Logout) // 退出
	}

	b.Use(midWare.Auth())
	{
		r := b.Group("rss")                   // /api/v1/rss
		r.GET("info/mikan", api.GetMikanInfo) // 获取mikan的番剧信息
		r.GET("search", api.GetSearch)        // 获取mikan的搜索rss
		r.GET("subject", api.GetRssSubject)   // 根据mikan的ID与字幕组ID获取rss
	}
	{
		r := b.Group("category")              // /api/v1/category
		r.POST("", api.CreateCategory)        // 添加分类
		r.GET("list", api.GetCategoryList)    // 获取分类列表
		r.DELETE("", api.DeleteCategory)      // 删除分类
		r.PUT("name", api.UpdateCategoryName) // 更新分类名
	}
	{
		r := b.Group("video") // /api/v1/video

		r.GET("link", api.Link)                    // 手动执行硬链接
		r.POST("subscribe", api.Subscribe)         // 订阅番剧
		r.POST("add", api.Add)                     // 手动添加
		r.GET("detail", api.GetOne)                // 获取详情
		r.GET("list", api.GetList)                 // 获取番剧列表
		r.DELETE("", api.DeleteVideo)              // 删除番剧
		r.PUT("category", api.UpdateVideoCategory) // 更改分类
	}
	{
		r := b.Group("bangumi")            // /api/v1/bangumi
		r.GET("calendar", api.GetCalendar) // 获取新番日历
		r.GET("subject", api.GetSubject)   // 获取番剧详细信息
		r.GET("search", api.Search)        // 搜索番剧
	}
}
