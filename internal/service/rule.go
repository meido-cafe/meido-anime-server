package service

import (
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
)

func (this *Service) GetRuleList() (res []model.Rule, err error) {
	return this.repo.RuleSelectList()
}

func (this *Service) DeleteRule(req vo.DeleteRuleRequest) (err error) {
	return this.repo.RuleDeleteList(req.ID)
}

func (this *Service) AddRuleList(req vo.AddRuleListRequest) (err error) {

	data := make([]model.Rule, 0, len(req.Rule))
	for _, item := range req.Rule {
		data = append(data, model.Rule{
			Name:           item.Name,
			MustContain:    item.MustContain,
			MustNotContain: item.MustNotContain,
			UseRegex:       item.UseRegex,
			EpisodeFilter:  item.EpisodeFilter,
			SmartFilter:    item.SmartFilter,
			Time:           model.NewTime(),
		})
	}
	return this.repo.RuleInsertList(data)
}

func (this *Service) UpdateRule(req vo.UpdateRuleRequest) (err error) {
	return this.repo.RuleUpdateOne(model.Rule{
		Id:             req.Id,
		Name:           req.Name,
		MustContain:    req.MustContain,
		MustNotContain: req.MustNotContain,
		UseRegex:       req.UseRegex,
		EpisodeFilter:  req.EpisodeFilter,
		SmartFilter:    req.SmartFilter,
		Time:           model.NewTime(),
	})
}

func (this *Service) CheckRuleNameRepeated(nameList []string) (res []string, err error) {
	return this.repo.RuleSelectNameListByName(nameList)
}
