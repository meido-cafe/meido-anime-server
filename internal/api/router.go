package api

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/factory"
)

func InitRouter(router *gin.RouterGroup) {
	videoApi := factory.NewVideoApi()
	rssApi := factory.NewRssApi()
	userApi := factory.NewUserApi()
	bangumiApi := factory.NewBangumiApi()
	middleware := factory.NewMiddleware()

	base := router.Group("/api")
	b := base.Group("v1")
	{
		r := b.Group("user")
		r.POST("login", userApi.Login)                     // 登录
		r.GET("logout", middleware.Auth(), userApi.Logout) // 退出
	}

	b.Use(middleware.Auth())
	{
		r := b.Group("rss")                      // /api/v1/rss
		r.GET("info/mikan", rssApi.GetMikanInfo) // 获取mikan的番剧信息
		r.GET("search", rssApi.GetSearch)        // 获取mikan的搜索rss
		r.GET("subject", rssApi.GetSubject)      // 根据mikan的ID与字幕组ID获取rss
	}
	{
		r := b.Group("video") // /api/v1/video

		r.GET("link", videoApi.Link)                    // 手动执行硬链接
		r.POST("subscribe", videoApi.Subscribe)         // 订阅番剧
		r.GET("detail", videoApi.GetOne)                // 获取详情
		r.GET("list", videoApi.GetList)                 // 获取番剧列表
		r.DELETE("", videoApi.DeleteVideo)              // 删除番剧
		r.PUT("category", videoApi.UpdateVideoCategory) // 更改分类

	}
	{
		r := b.Group("category")                   // /api/v1/category
		r.POST("", videoApi.CreateCategory)        // 添加分类
		r.GET("list", videoApi.GetCategoryList)    // 获取分类列表
		r.DELETE("", videoApi.DeleteCategory)      // 删除分类
		r.PUT("name", videoApi.UpdateCategoryName) // 更新分类名
	}
	{
		r := b.Group("bangumi")                   // /api/v1/bangumi
		r.GET("calendar", bangumiApi.GetCalendar) // 获取新番日历
		r.GET("subject", bangumiApi.GetSubject)   // 获取番剧详细信息
		r.GET("search", bangumiApi.Search)        // 搜索番剧
	}
}
