package service

import (
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
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

func (this *Service) GetSubjectCharacters(id int) (ret []model.BangumiSubjectCharacter, err error) {
	return this.bangumi.GetSubjectCharacters(id)
}

func (this *Service) GetIndex(request vo.GetIndexRequest) (ret model.BangumiIndexResponse, err error) {
	limit := request.PageSize
	offset := (request.Page - 1) * request.PageSize
	ret, err = this.bangumi.GetIndex(limit, offset, request.GetIndexRequestBody)
	if err != nil {
		return
	}
	return
}
