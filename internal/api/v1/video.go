package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model/vo"
	"strings"
)

func (this *Api) Link(ctx *gin.Context) {
	this.service.Link()
	response.Success(ctx)
}

func (this *Api) GetList(ctx *gin.Context) {
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

func (this *Api) GetOne(ctx *gin.Context) {
	req := vo.VideoGetOneRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
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

func (this *Api) Subscribe(ctx *gin.Context) {
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
	case req.Category == 0:
		response.Bad(ctx, "分类不能为空")
		return
	}
	category, err := this.service.GetCategory(req.Category)
	if err != nil {
		response.Error(ctx, "订阅失败")
		return
	}

	if category.Id == 0 {
		response.Bad(ctx, "分类不存在")
		return
	}

	exist, err := this.service.GetOne(vo.VideoGetOneRequest{
		BangumiId: req.BangumiId,
	})

	if err != nil {
		response.Error(ctx, "订阅失败")
		return
	}

	if exist.Id > 0 {
		response.Bad(ctx, "番剧 bangumi id 已存在")
		return
	}

	exist, err = this.service.GetOne(vo.VideoGetOneRequest{
		Title: req.Title,
	})
	if err != nil {
		response.Error(ctx, "订阅失败")
		return
	}

	if exist.Id > 0 {
		response.Bad(ctx, "番剧名称已存在")
		return
	}

	if err = this.service.Subscribe(req); err != nil {
		response.Error(ctx, "订阅失败")
		return
	}
	response.Success(ctx)
}

// Add 手动添加订阅
func (this *Api) Add(ctx *gin.Context) {
	req := vo.VideoAdd{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	// 参数校验
	switch {
	case req.Title == "":
		response.Bad(ctx, "番剧名称不能为空")
		return
	case req.Category == 0:
		response.Bad(ctx, "分类不能为空")
		return
	}

	if req.Mode == 1 {
		if req.RssUrl == "" {
			response.Bad(ctx, "rss链接不能为空")
			return
		}
	} else if req.Mode == 2 {
		if len(req.TorrentList) == 0 {
			response.Bad(ctx, "种子列表不能为空")
			return
		}
	} else {
		response.Bad(ctx, "创建模式错误")
		return
	}

	category, err := this.service.GetCategory(req.Category)
	if err != nil {
		response.Error(ctx, "订阅失败")
		return
	}

	if category.Id == 0 {
		response.Bad(ctx, "分类不存在")
		return
	}

	exist, err := this.service.GetOne(vo.VideoGetOneRequest{
		Title: req.Title,
	})
	if err != nil {
		response.Error(ctx, "添加失败")
		return
	}

	if exist.Id > 0 {
		response.Bad(ctx, "番剧名称已存在")
		return
	}

	// 添加
	if err := this.service.Add(req); err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx)
}

func (this *Api) DeleteVideo(ctx *gin.Context) {
	req := vo.DeleteVideoRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}
	if req.Id == 0 {
		response.Bad(ctx, "id不能为空")
		return
	}
	if err := this.service.DeleteVideo(req); err != nil {
		response.Error(ctx, "删除订阅失败")
		return
	}

	response.Success(ctx)
}

func (this *Api) UpdateVideoCategory(ctx *gin.Context) {
	req := vo.UpdateVideoCategoryRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	switch {
	case len(req.Ids) == 0:
		response.Bad(ctx, "ID不能为空")
		return
	case req.Category == 0:
		response.Bad(ctx, "分类不能为空")
		return
	}

	category, err := this.service.GetCategory(req.Category)
	if err != nil {
		response.Bad(ctx, "更新失败")
		return
	}
	if category.Id == 0 {
		response.Bad(ctx, "分类不存在")
		return
	}

	if err = this.service.UpdateVideoCategory(req); err != nil {
		response.Error(ctx, "更新失败")
		return
	}

	response.Success(ctx)
	return
}

func (this *Api) CreateCategory(ctx *gin.Context) {
	req := vo.CreateCategoryRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}
	if req.Name == "" {
		response.Bad(ctx, "分类名称不能为空")
		return
	}

	category, err := this.service.GetCategoryByName(req.Name)
	if err != nil {
		response.Error(ctx, "创建失败")
		return
	}

	if category.Id > 0 {
		response.Bad(ctx, "分类已存在")
		return
	}

	if err := this.service.CreateCategory(req); err != nil {
		response.Error(ctx, "创建失败")
		return
	}
	response.Success(ctx)
	return
}

func (this *Api) GetCategoryList(ctx *gin.Context) {
	list, err := this.service.GetCategoryList()
	if err != nil {
		response.Error(ctx, "获取失败")
		return
	}
	response.List(ctx, list, len(list))
	return
}

func (this *Api) DeleteCategory(ctx *gin.Context) {
	req := vo.DeleteCategoryRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}
	if req.Id == 0 {
		response.Bad(ctx, "ID不能为空")
		return
	}

	if err := this.service.DeleteCategory(req.Id); err != nil {
		response.Error(ctx, "删除失败")
		return
	}
	response.Success(ctx)
	return
}

func (this *Api) UpdateCategoryName(ctx *gin.Context) {
	req := vo.UpdateCategoryNameRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	switch {
	case req.Id == 0:
		response.Bad(ctx, "ID不能为空")
		return
	case req.Name == "":
		response.Bad(ctx, "分类名称不能为空")
		return
	}

	result, err := this.service.GetCategoryByName(req.Name)
	if err != nil {
		response.Error(ctx, "更新失败")
		return
	}
	if result.Id > 0 {
		if result.Id == req.Id {
			response.Success(ctx)
			return
		}
		response.Bad(ctx, "分类已存在")
		return
	}

	// 更新
	if err = this.service.UpdateCategoryName(req); err != nil {
		response.Error(ctx, "更新失败")
		return
	}
	response.Success(ctx)
	return
}
