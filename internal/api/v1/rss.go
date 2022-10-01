package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model/vo"
)

func (this *Api) GetMikanInfo(ctx *gin.Context) {
	req := vo.GetRssInfoMikanRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}

	switch {
	case req.SubjectName == "":
		response.Bad(ctx, "番剧名称不能为空")
		return
	}

	res, err := this.service.GetInfoMikan(req)
	if err != nil {
		response.Error(ctx, "获取番剧mikan信息失败")
		return
	}
	response.Data(ctx, res)
}

func (this *Api) GetSearch(ctx *gin.Context) {
	req := vo.GetRssSearchRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}

	switch {
	case req.SubjectName == "":
		response.Bad(ctx, "番剧名称不能为空")
		return
	}

	res, err := this.service.GetSearch(req)
	if err != nil {
		response.Error(ctx, "获取rss信息失败")
		return
	}
	response.Data(ctx, res)
}

func (this *Api) GetRssSubject(ctx *gin.Context) {
	req := vo.GetRssSubjectRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		log.Println(err)
		response.BadBind(ctx)
		return
	}

	switch {
	case req.MikanId == 0:
		response.Bad(ctx, "mikan 番剧id不能为空")
		return
	case req.MikanGroupId == 0:
		response.Bad(ctx, "mikan 字幕组id不能为空")
		return
	}

	res, err := this.service.GetRssSubject(req)
	if err != nil {
		response.Error(ctx, "获取rss信息失败")
		return
	}
	response.Data(ctx, res)
}
