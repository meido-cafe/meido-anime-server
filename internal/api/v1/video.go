package v1

import (
	"github.com/gin-gonic/gin"
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

func (this *VideoApi) GetList(ctx *gin.Context) {
	req := vo.VideoGetListRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
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

// 有bangumi 订阅rss
func (this *VideoApi) Subscribe(ctx *gin.Context) {
	req := vo.VideoSubscribeRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
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
	case req.RssUrl == "":
		response.Bad(ctx, "rss订阅链接不能为空")
		return
	}

	exist, err := this.service.GetByBangumiId(req.BangumiId)
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
