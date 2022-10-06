package repo

import (
	"github.com/imroc/req/v3"
	"log"
	"meido-anime-server/internal/common"
	"meido-anime-server/internal/model"
	"strconv"
)

type BangumiRepo struct {
	client *req.Client
}

func NewBangumiRepo() *BangumiRepo {
	return &BangumiRepo{client: common.GetBangumiClient()}
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

func (this *BangumiRepo) GetSubjectCharacters(id int) (ret []model.BangumiSubjectCharacter, err error) {
	res, err := this.client.R().
		SetPathParam("id", strconv.Itoa(id)).
		SetResult(&ret).
		Get("/v0/subjects/{id}/characters")
	if err != nil {
		log.Println(err)
		return
	}
	if res.IsError() {
		log.Printf("[获取番剧角色信息失败] [code] %d [response] %s\n", res.StatusCode, res.String())
		return
	}
	return
}

func (this *BangumiRepo) GetIndex(sort string, t string, page int, year int, month int) (ret model.BangumiIndexResponse, err error) {
	//https://bangumi.tv/anime/browser/{type}/airtime/{year}-{month}?sort=&page=
	//var class string
	//if t != "" {
	//	class = t + "/"
	//}
	//baseUrl := `https://bangumi.tv/anime/browser`
	//if t != "" {
	//	baseUrl += "/" + t
	//}
	//if year > 0 {
	//
	//}
	//url := fmt.Sprintf("https://bangumi.tv/anime/browser/?sort=%s&page=%d", sort, page)
	return
}
