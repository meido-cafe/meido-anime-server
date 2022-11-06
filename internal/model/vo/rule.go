package vo

import "meido-anime-server/internal/model"

type DeleteRuleRequest struct {
	ID int64 `json:"id"`
}
type AddRuleListRequest struct {
	Rule []model.Rule `json:"rule" form:"rule"`
}
type UpdateRuleRequest struct {
	model.Rule
}
