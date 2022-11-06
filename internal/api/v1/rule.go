package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model/vo"
)

func (this *Api) GetRuleList(ctx *gin.Context) {
	list, err := this.service.GetRuleList()
	if err != nil {
		response.Error(ctx, "获取规则列表失败")
		return
	}
	response.List(ctx, list, len(list))
}
func (this *Api) DeleteRule(ctx *gin.Context) {
	var req vo.DeleteRuleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	if len(req.ID) == 0 {
		response.Error(ctx, "ID列表不能为空")
		return
	}

	if err := this.service.DeleteRule(req); err != nil {
		response.Error(ctx, "删除规则失败")
		return
	}
	response.Success(ctx)
}
func (this *Api) AddRuleList(ctx *gin.Context) {
	var req vo.AddRuleListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}
	if len(req.Rule) == 0 {
		response.Bad(ctx, "规则列表不能为空")
		return
	}

	nameSet := make(map[string]struct{})
	nameList := make([]string, 0, len(req.Rule))

	for _, item := range req.Rule {
		if _, ok := nameSet[item.Name]; ok {
			response.Bad(ctx, "规则名称不能重复")
			return
		}
		nameSet[item.Name] = struct{}{}

		if item.UseRegex < 1 || item.UseRegex > 2 {
			response.Bad(ctx, "use_regex 不支持的类型")
			return
		}

		if item.SmartFilter < 1 || item.SmartFilter > 2 {
			response.Bad(ctx, "smart_filter 不支持的类型")
			return
		}

		nameList = append(nameList, item.Name)
	}

	ret, err := this.service.CheckRuleNameRepeated(nameList)
	if err != nil {
		response.Error(ctx, "创建规则失败")
		return
	}

	if len(ret) > 0 {
		response.Bad(ctx, "规则名称已存在")
		return
	}

	if err = this.service.AddRuleList(req); err != nil {
		response.Error(ctx, "创建规则失败")
		return
	}

	response.Success(ctx)
}
func (this *Api) UpdateRule(ctx *gin.Context) {
	var req vo.UpdateRuleRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	if req.Id == 0 {
		response.Bad(ctx, "id不能为空")
		return
	}

	if req.UseRegex < 1 || req.UseRegex > 2 {
		response.Bad(ctx, "use_regex 不支持的类型")
		return
	}

	if req.SmartFilter < 1 || req.SmartFilter > 2 {
		response.Bad(ctx, "smart_filter 不支持的类型")
		return
	}

	ret, err := this.service.CheckRuleNameRepeated([]string{req.Name})
	if err != nil {
		response.Error(ctx, "更新规则失败")
		return
	}

	if len(ret) > 0 {
		response.Bad(ctx, "规则名称已存在")
		return
	}

	if err = this.service.UpdateRule(req); err != nil {
		response.Error(ctx, "更新规则失败")
		return
	}
	response.Success(ctx)
}
