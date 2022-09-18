package repo

import (
	"github.com/imroc/req/v3"
	"log"
	"meido-anime-server/internal/model"
	"strconv"
)

type BangumiRepo struct {
	client *req.Client
}

func NewBangumiRepo(client *req.Client) *BangumiRepo {
	return &BangumiRepo{client: client}
}

func (this *BangumiRepo) GetCalendar() (ret []model.BangumiCalendar, err error) {
	res := this.client.Get("calendar").Do()
	if err = res.UnmarshalJson(&ret); err != nil {
		log.Println(err)
		return
	}
	return
}
func (this *BangumiRepo) GetSubject(id int) (ret model.BangumiSubject, err error) {
	res, err := this.client.R().
		SetPathParam("id", strconv.Itoa(id)).
		SetResult(&ret).
		Get("v0/subjects/{id}")
	if err != nil {
		log.Println(err)
		return
	}

	if res.IsError() {
		log.Printf("[获取番剧信息失败] [code] %d [response] %s\n", res.StatusCode, res.String())
		return
	}

	return
}

func (this *BangumiRepo) Search(name string, class int) (ret model.BangumiSearchSubjectResponse, err error) {
	res, err := this.client.R().
		SetPathParam("name", name).
		SetQueryParam("type", strconv.Itoa(class)).
		SetResult(&ret).
		Get("/search/subject/{name}")
	if err != nil {
		log.Println(err)
		return
	}
	if res.IsError() {
		log.Printf("[搜索番剧失败] [code] %d [response] %s\n", res.StatusCode, res.String())
		return
	}
	return
}