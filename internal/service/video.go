package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"meido-anime-server/internal/global"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/repo"
	"meido-anime-server/internal/tool"
	"os"
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

func (this *VideoService) CacheSourcePath() (err error) {
	list, _, err := this.repo.SelectItemList(context.TODO(), 0)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range list {
		global.LinkPathCache.Store(item.LinkPath, struct{}{})
	}
	return
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
			log.Println(err)
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

	downloadPath := fmt.Sprintf("%s/%s/S%02d", strings.TrimRight(global.QBDownloadPath, "/"), request.Title, season)
	sourcePath := filepath.Join(global.SourcePath, request.Title, fmt.Sprintf("S%02d", season))

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
		SourceDir: sourcePath,
		Time: model.Time{
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		},
	}

	qbRule := model.QBRule{
		Enabled:                   true,
		MustContain:               request.MustContain,
		MustNotContain:            request.MustNotContain,
		UseRegex:                  false,
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
	log.Printf(`[bangumi:%d][%s][第%d季] 订阅成功,qbittorrent下载路径:[%s],meido-anime读取路径:[%s]`, request.BangumiId, request.Title, season, downloadPath, sourcePath)
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

func (this *VideoService) Link() {
	// 获取番剧列表
	list, total, err := this.repo.SelectList(context.TODO(), vo.VideoGetListRequest{})
	if err != nil {
		log.Printf("[硬链接失败][获取视频列表失败][error]:%v\n", err)
		return
	}
	start := time.Now()

	videoItemList := make([]model.VideoItem, 0, total)
	linkSourceList := make([]model.LinkSource, 0, total)

	for _, item := range list {
		linkDir := filepath.Join(global.MediaPath, item.Title, fmt.Sprintf("S%02d", item.Season))
		// 创建媒体文件夹
		if err := os.MkdirAll(linkDir, 0666); err != nil {
			log.Printf("[硬链接失败][创建媒体库路径失败] %s [已跳过该番剧][error]: %v\n", item.SourceDir, err)
			continue
		}

		// 检测下载路径是否存在
		exists, err := tool.IsFileExists(item.SourceDir)
		if err != nil {
			log.Printf("[硬链接失败][获取番剧路径信息失败] %s [已跳过该番剧][error]: %v\n", item.SourceDir, err)
			continue
		}
		if !exists {
			log.Printf("[硬链接失败][番剧路径不存在] %s [已跳过该番剧][error]: %v\n", item.SourceDir, err)
			continue
		}

		// 获取下载路径下的所有文件
		dir, err := ioutil.ReadDir(item.SourceDir)
		if err != nil {
			log.Printf("[硬链接失败][读取番剧路径失败] %s [已跳过该番剧][error]: %v\n", item.SourceDir, err)
			continue
		}

		// 构造源文件列表
		for _, file := range dir {
			if file.IsDir() {
				continue
			}

			filename := file.Name()
			// 从下载视频文件提取集数
			episode, err := tool.GetEpisode(filename)
			if err != nil {
				log.Printf("[硬链接失败][获取集数失败] %s [已跳过该集][error]: %v\n", filename, err)
				continue
			}
			// 硬链接文件标题
			title := fmt.Sprintf("%s - S%02dE%02d%s", item.Title, item.Season, episode, filepath.Ext(filename))
			linkPath := filepath.Join(linkDir, title)

			// 跳过已存在的link path
			if _, ok := global.LinkPathCache.Load(linkPath); !ok {
				linkSourceList = append(linkSourceList, model.LinkSource{
					Id:             item.Id,
					Episode:        int64(episode),
					SourceFilePath: filepath.Join(item.SourceDir, filename),
					LinkFilePath:   linkPath,
				})
				global.LinkPathCache.Store(linkPath, struct{}{})
			}

		}
	}

	if len(linkSourceList) == 0 {
		return
	}

	for _, item := range linkSourceList {
		// 判断下载文件是否存在
		ok, err := tool.IsFileExists(item.SourceFilePath)
		if err != nil {
			log.Printf("[硬链接失败][获取下载视频文件失败] %s [已跳过该文件][error]:%v\n", item.SourceFilePath, err)
			continue
		}
		if !ok {
			log.Printf("[硬链接失败][下载视频文件不存在] %s [已跳过该文件] \n ", item.SourceFilePath)
			continue
		}

		// 判断硬链接文件是否存在
		ok, err = tool.IsFileExists(item.LinkFilePath)
		if err != nil {
			log.Printf("[硬链接失败][尝试获取硬链接文件信息错误] %s [已跳过该文件][error]:%v\n", item.LinkFilePath, err)
			continue
		}
		if ok {
			log.Printf("[硬链接失败][硬链接文件已存在] %s [已跳过该文件] \n ", item.LinkFilePath)
			continue
		}

		// 构造插入数据
		videoItemList = append(videoItemList, model.VideoItem{
			Pid:        item.Id,
			Episode:    item.Episode,
			SourcePath: item.SourceFilePath,
			LinkPath:   item.LinkFilePath,
			Time: model.Time{
				CreateTime: time.Now().Unix(),
				UpdateTime: time.Now().Unix(),
			},
		})
	}
	// 插入数据
	cnt, err := this.repo.InsertItemList(context.TODO(), videoItemList)
	if err != nil {
		log.Printf("[硬链接失败][插入数据库失败] %v \n", err)
		return
	}

	end := time.Now()
	difference := end.Sub(start)
	log.Printf("[硬链接任务完成] [共计执行成功%d次硬链接] [共计耗时 %s ] \n", cnt, difference.String())
}
