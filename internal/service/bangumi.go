package service

import (
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/repo"
)

type BangumiService struct {
	repo *repo.BangumiRepo
}

func NewBangumiService(repo *repo.BangumiRepo) *BangumiService {
	return &BangumiService{repo: repo}
}

func (this *BangumiService) GetCalendar() (ret []model.BangumiCalendar, total int, err error) {
	ret, err = this.repo.GetCalendar()
	for _, item := range ret {
		total += len(item.Items)
	}
	return
}

func (this *BangumiService) GetSubject(id int) (ret model.BangumiSubject, err error) {
	return this.repo.GetSubject(id)
}

func (this *BangumiService) Search(name string, class int) (ret []model.BangumiSearchSubjectItem, total int, err error) {
	searchRet, err := this.repo.Search(name, class)
	if err != nil {
		return
	}
	ret = searchRet.List
	total = searchRet.Results
	return
}
