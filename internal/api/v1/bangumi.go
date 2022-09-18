package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/service"
)

type BangumiApi struct {
	service *service.BangumiService
}

func NewBangumiApi(service *service.BangumiService) *BangumiApi {
	return &BangumiApi{service: service}
}

func (this *BangumiApi) GetCalendar(ctx *gin.Context) {
	calendarList, total, err := this.service.GetCalendar()
	if err != nil {
		response.Error(ctx, "获取失败")
		return
	}
	response.List(ctx, calendarList, total)
}

func (this *BangumiApi) GetSubject(ctx *gin.Context) {
	req := vo.GetSubjectRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	if req.Id == 0 {
		response.Bad(ctx, "bangumi id 不能为空")
		return
	}

	subject, err := this.service.GetSubject(req.Id)
	if err != nil {
		response.Error(ctx, "获取失败")
		return
	}
	response.Data(ctx, subject)
}

func (this *BangumiApi) Search(ctx *gin.Context) {
	req := vo.SearchRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	if req.Name == "" {
		response.Bad(ctx, "番剧名称不能为空")
		return
	}

	search, total, err := this.service.Search(req.Name, req.Class)
	if err != nil {
		response.Error(ctx, "搜索失败")
		return
	}
	response.List(ctx, search, total)
}