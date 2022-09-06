package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/service"
	"strings"
)

func NewVideoApi(service *service.VideoService) *VideoApi {
	return &VideoApi{service: service}
}

type VideoApi struct {
	service *service.VideoService
}

func (this *VideoApi) Link(ctx *gin.Context) {
	this.service.Link()
	response.Success(ctx)
}

func (this *VideoApi) GetList(ctx *gin.Context) {
	req := vo.VideoGetListRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		return
	}

	if err := req.CheckPage(); err != nil {
		response.Bad(ctx, err.Error())
		return
	}

	res, err := this.service.GetList(req)
	if err != nil {
		response.Error(ctx, "获取番剧列表失败")
		return
	}
	response.List(ctx, res.Items, res.Total)
}

func (this *VideoApi) GetOne(ctx *gin.Context) {
	req := vo.VideoGetOneRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}
	if req.Id == 0 && req.BangumiId == 0 {
		response.Bad(ctx, "缺少参数")
		return
	}
	res, err := this.service.GetOne(req)
	if err != nil {
		response.Error(ctx, "查询失败")
		return
	}
	response.Data(ctx, res)
}

// 有bangumi 订阅rss
func (this *VideoApi) Subscribe(ctx *gin.Context) {
	req := vo.VideoSubscribeRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.RssUrl = strings.TrimSpace(req.RssUrl)

	switch {
	case req.BangumiId == 0:
		response.Bad(ctx, "bangumi id 不能为空")
		return
	case req.Title == "":
		response.Bad(ctx, "番剧名称不能为空")
		return
	}

	exist, err := this.service.GetOne(vo.VideoGetOneRequest{
		Id:        0,
		BangumiId: req.BangumiId,
	})

	if err != nil {
		response.Error(ctx, "下载记录添加失败")
		return
	}

	if exist.Id > 0 {
		response.Bad(ctx, "番剧已存在")
		return
	}

	if err := this.service.Subscribe(req); err != nil {
		response.Error(ctx, "下载记录添加失败")
		return
	}
	response.Success(ctx)
}

func (this *VideoApi) DeleteRss(ctx *gin.Context) {
	req := vo.DeleteRssRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}
	if req.Id == 0 {
		response.Bad(ctx, "id不能为空")
		return
	}
	if err := this.service.DeleteRss(req); err != nil {
		response.Error(ctx, "")
		return
	}

	response.Success(ctx)
}

func (this *VideoApi) UpdateRss(ctx *gin.Context) {
	req := vo.UpdateRssRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}
	if req.Id == 0 {
		response.Bad(ctx, "id不能为空")
		return
	}
	if req.Rss == "" {
		response.Bad(ctx, "rss链接不能为空")
		return
	}
	if err := this.service.UpdateRss(req); err != nil {
		response.Error(ctx, "")
		return
	}

	response.Success(ctx)

}

func (this *VideoApi) GetQBLogs(ctx *gin.Context) {
	logs, err := this.service.GetQBLogs()
	if err != nil {
		response.Error(ctx, "获取日志失败")
		return
	}
	response.List(ctx, logs.Items, logs.Total)
}
