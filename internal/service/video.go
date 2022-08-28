package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meido-anime-server/internal/global"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/tool"
	"path/filepath"
	"strings"
	"time"
)

type VideoService struct {
	repo repo.VideoInterface
}

func NewVideoService(repo repo.VideoInterface) *VideoService {
	return &VideoService{
		repo: repo,
	}
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
		season, matchStr, err = tool.GetSeason(request.Title)
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

	downloadPath := filepath.Join(global.QBDownloadPath, request.Title, fmt.Sprintf("S%02d", season))

	data := model.Video{
		BangumiId:    request.BangumiId,
		Origin:       1,
		Category:     request.Category,
		Title:        request.Title,
		Season:       season,
		Cover:        request.Cover,
		Total:        request.Total,
		RssUrl:       request.RssUrl,
		PlayTime:     request.PlayTime,
		DownloadPath: downloadPath,
		Time: model.Time{
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		},
	}

	qbRule := model.QBRule{
		Enabled:                   true,
		MustContain:               request.MustContain,
		MustNotContain:            request.MustNotContain,
		UseRegex:                  true,
		EpisodeFilter:             request.EpisodeFilter,
		SmartFilter:               request.SmartFilter,
		PreviouslyMatchedEpisodes: nil,
		AffectedFeeds:             []string{request.RssUrl},
		IgnoreDays:                0,
		LastMatch:                 "",
		AddPaused:                 false,
		AssignedCategory:          global.QBCategory,
		SavePath:                  downloadPath,
	}

	marshal, err := json.Marshal(qbRule)
	if err != nil {
		log.Printf("[%s]解析失败:%v", request.Title, err)
		return err
	}
	setRule := model.QBSetRule{
		RuleName: fmt.Sprintf("%s - S%02d", request.Title, season),
		RuleDef:  string(marshal),
	}

	if err = this.repo.InsertOne(context.TODO(), data, setRule); err != nil {
		log.Println(err)
		return
	}
	log.Printf(`[bangumi:%d][%s][第%d季] 订阅成功,下载路径:[%s]`, request.BangumiId, request.Title, season, downloadPath)
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

func (this *VideoService) GetQBLogs() (res vo.GetQBLogsResponse, err error) {
	logs, err := this.repo.GetQBLogs(context.TODO())
	if err != nil {
		log.Println(err)
		return
	}
	for _, item := range logs {
		res.Items = append(res.Items, vo.GetQBLogsResponseItem{
			QBLog: item,
			Time:  time.UnixMilli(item.Timestamp).Format("2006-01-02 15:04:05"),
		})
	}

	return
}
