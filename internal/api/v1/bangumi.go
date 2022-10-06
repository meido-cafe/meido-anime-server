package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model/vo"
)

func (this *Api) GetCalendar(ctx *gin.Context) {
	calendarList, total, err := this.service.GetCalendar()
	if err != nil {
		response.Error(ctx, "获取失败")
		return
	}
	response.List(ctx, calendarList, total)
}

func (this *Api) GetIndex(ctx *gin.Context) {
	req := vo.GetIndexRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	switch {
	case req.Sort != "" && req.Sort != "rank" && req.Sort != "date" && req.Sort != "title":
		response.Bad(ctx, "不支持的排序方式")
		return
	case req.Type != "" && req.Type != "tv" && req.Type != "web" && req.Type != "ova" && req.Type != "movie":
		response.Bad(ctx, "不支持的番剧类型")
		return
	}

	index, err := this.service.GetIndex(req)
	if err != nil {
		response.Error(ctx, "获取番剧索引失败")
		return
	}

	response.Data(ctx, index)
}
func (this *Api) GetSubject(ctx *gin.Context) {
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

func (this *Api) Search(ctx *gin.Context) {
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

func (this *Api) GetSubjectCharacters(ctx *gin.Context) {
	req := vo.GetSubjectCharactersRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	if req.Id == 0 {
		response.Bad(ctx, "bangumi id 不能为空")
		return
	}

	res, err := this.service.GetSubjectCharacters(req.Id)
	if err != nil {
		response.Error(ctx, "获取角色信息失败")
		return
	}

	response.List(ctx, res, len(res))
}
