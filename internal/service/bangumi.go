package service

import (
	"meido-anime-server/internal/model"
)

func (this *Service) GetCalendar() (ret []model.BangumiCalendar, total int, err error) {
	ret, err = this.bangumi.GetCalendar()
	for _, item := range ret {
		total += len(item.Items)
	}
	return
}

func (this *Service) GetSubject(id int) (ret model.BangumiSubject, err error) {
	return this.bangumi.GetSubject(id)
}

func (this *Service) Search(name string, class int) (ret []model.BangumiSearchSubjectItem, total int, err error) {
	searchRet, err := this.bangumi.Search(name, class)
	if err != nil {
		return
	}
	ret = searchRet.List
	total = searchRet.Results
	return
}
