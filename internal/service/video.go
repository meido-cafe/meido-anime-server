package service

import (
	"context"
	"log"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/repo"
	"meido-anime-server/pkg"
	"strings"
	"time"
)

func NewVideoService(repo repo.VideoInterface) *VideoService {
	return &VideoService{
		repo: repo,
	}
}

type VideoService struct {
	repo repo.VideoInterface
}

func (this *VideoService) GetOne(request vo.VideoGetOneRequest) (ret model.Video, err error) {
	ret, err = this.repo.SelectOne(context.TODO(), request.Id, request.BangumiId)
	if err != nil {
		log.Println("获取番剧信息失败")
		return
	}
	return
}

// Subscribe 订阅RSS
func (this *VideoService) Subscribe(request vo.VideoSubscribeRequest) (err error) {
	// 如果参数没有传第几季, 则从标题正则获取, 获取不到视为第一季
	var season int64
	var matchStr string
	if request.Season > 0 {
		season = request.Season
	} else {
		season, matchStr, err = pkg.GetSeason(request.Title)
		if err != nil {
			return
		}
		if matchStr == "" {
			season = 1
			log.Println("没有从", request.Title, "中获取季相关信息,已设置为默认的第一季")
		} else {
			old := request.Title
			request.Title = strings.ReplaceAll(request.Title, matchStr, "")
			request.Title = strings.TrimSpace(request.Title)
			log.Println(old, " 重命名为 ==> ", request.Title)
		}
	}

	data := model.Video{
		BangumiId: request.BangumiId,
		Origin:    1,
		Category:  request.Category,
		Title:     request.Title,
		Season:    season,
		Cover:     request.Cover,
		Total:     request.Total,
		RssUrl:    request.RssUrl,
		PlayTime:  request.PlayTime,
		Time: model.Time{
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		},
	}

	if err = this.repo.InsertOne(context.TODO(), data); err != nil {
		log.Println(err)
		return
	}
	log.Printf(`[bangumi:%d][%s][第%d季] 添加成功`, request.BangumiId, request.Title, season)
	return
}

func (this *VideoService) GetList(request vo.VideoGetListRequest) (response vo.VideoGetListResponse, err error) {
	list, total, err := this.repo.SelectList(context.TODO(), request)
	if err != nil {
		log.Println(err)
		return
	}
	response.Items = list
	response.Total = total
	return
}

func (this *VideoService) DeleteRss(request vo.DeleteRssRequest) (err error) {
	err = this.repo.UpdateRss(context.TODO(), request.Id, "")
	if err != nil {
		log.Println(err)
	}
	return

}
func (this *VideoService) UpdateRss(request vo.UpdateRssRequest) (err error) {
	err = this.repo.UpdateRss(context.TODO(), request.Id, request.Rss)
	if err != nil {
		log.Println(err)
	}
	return
}
