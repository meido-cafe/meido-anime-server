package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/service"
)

type RssApi struct {
	service *service.RssService
}

func NewRssApi(service *service.RssService) *RssApi {
	return &RssApi{service: service}
}

func (this *RssApi) GetMikanInfo(ctx *gin.Context) {
	req := vo.GetInfoMikanRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		model.BadBind(ctx)
		return
	}
	if req.SubjectName == "" {
		model.Bad(ctx, "番剧名称不能为空")
		return
	}

	res, err := this.service.GetInfoMikan(req)
	if err != nil {
		model.Error(ctx, "获取番剧mikan信息失败", err.Error())
		return
	}
	model.Data(ctx, res)
}
