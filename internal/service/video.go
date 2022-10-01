package service

import (
	"context"
	"encoding/json"
	"errors"
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

// 获取分类的目录路径
func (this *Service) getCategoryDir(categoryName string) string {
	return filepath.Join(global.MediaPath, categoryName)
}

// 获取番剧硬链接路径
func (this *Service) getLinkDir(categoryName string, title string, season int64) string {
	return filepath.Join(global.MediaPath, categoryName, title, fmt.Sprintf("S%02d", season))
}

// CacheLinkPath 缓存已经硬链接的路径
func (this *Service) CacheLinkPath() (err error) {
	list, _, err := this.repo.VideoItemSelectList(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	for _, item := range list {
		global.LinkPathCache.Store(item.LinkPath, struct{}{})
	}
	return
}

// GetOne 根据ID或bangumi id 获取一个番剧的信息
func (this *Service) GetOne(request vo.VideoGetOneRequest) (ret model.Video, err error) {
	ret, err = this.repo.VideoSelectOne(context.TODO(), request.Id, request.BangumiId, request.Title)
	if err != nil {
		log.Println("获取番剧信息失败")
		return
	}
	return
}

// Subscribe 订阅RSS
//
//	只支持订阅, 手动添加种子需要使用添加种子接口
//	tx: 插入数据库, 订阅rss
func (this *Service) Subscribe(request vo.VideoSubscribeRequest) (err error) {
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

	// 获取分类信息
	category, err := this.repo.CategorySelectOne(context.TODO(), request.Category)
	if err != nil {
		log.Println("获取分类信息失败")
		return err
	}

	// 公用参数
	downloadPath := fmt.Sprintf("%s/%s/S%02d", strings.TrimRight(global.QBDownloadPath, "/"), request.Title, season)
	sourcePath := filepath.Join(global.SourcePath, request.Title, fmt.Sprintf("S%02d", season))
	linkDir := this.getLinkDir(category.Name, request.Title, season)
	ruleName := fmt.Sprintf("%s - S%02d", request.Title, season)

	// QB参数
	qbRule := model.QBDef{
		Enabled:                   true,
		MustContain:               request.MustContain,
		MustNotContain:            request.MustNotContain,
		UseRegex:                  request.UseRegex,
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

	setRule := map[string]string{
		"ruleName": ruleName,
		"ruleDef":  string(marshal),
	}

	// 插入video 参数
	data := model.Video{
		BangumiId: request.BangumiId,
		Origin:    1,
		Category:  request.Category,
		Title:     request.Title,
		Season:    season,
		Cover:     request.Cover,
		Total:     request.Total,
		RssUrl:    request.RssUrl,
		RuleName:  ruleName,
		PlayTime:  request.PlayTime,
		SourceDir: sourcePath,
		LinkDir:   linkDir,
		Time: model.Time{
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		},
	}

	ctx := context.TODO()
	ret := this.transaction(ctx, func(repo *repo.Repo) error {
		if err := repo.VideoInsertOne(ctx, data); err != nil {
			return err
		}
		if err := this.qb.AddRss(ctx, request.RssUrl, ruleName); err != nil {
			return err
		}
		if err := this.qb.SetRule(ctx, setRule); err != nil {
			return err
		}
		return nil
	})
	if ret.TxError != nil {
		err = ret.TxError
		log.Printf(`[bangumi:%d][%s][第%d季] 订阅失败 [error] %v \n`, request.BangumiId, request.Title, season, err)
		return
	}
	if ret.Error != nil {
		err = ret.Error
		log.Printf(`[bangumi:%d][%s][第%d季] 订阅失败 [error] %v \n`, request.BangumiId, request.Title, season, err)
		return
	}

	log.Printf(`[bangumi:%d][%s][第%d季] 订阅成功,qbittorrent下载路径:[%s],系统读取路径:[%s],硬链接路径:[%s]\n`, request.BangumiId, request.Title, season, downloadPath, sourcePath, linkDir)
	return
}

// Add 手动添加番剧
//
//	mode == 1: 订阅rss, mode==2: 添加种子方式
//	tx: 插入数据库, 判断mode 执行订阅rss或添加种子
func (this *Service) Add(request vo.VideoAdd) (err error) {
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
	// 获取分类信息
	category, err := this.repo.CategorySelectOne(context.TODO(), request.Category)
	if err != nil {
		log.Println("获取分类信息失败")
		return err
	}

	downloadPath := fmt.Sprintf("%s/%s/S%02d", strings.TrimRight(global.QBDownloadPath, "/"), request.Title, season)
	sourcePath := filepath.Join(global.SourcePath, request.Title, fmt.Sprintf("S%02d", season))
	linkDir := this.getLinkDir(category.Name, request.Title, season)
	ruleName := fmt.Sprintf("%s - S%02d", request.Title, season)

	// QB参数
	qbRule := model.QBDef{
		Enabled:                   true,
		MustContain:               request.MustContain,
		MustNotContain:            request.MustNotContain,
		UseRegex:                  request.UseRegex,
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

	setRule := map[string]string{
		"ruleName": ruleName,
		"ruleDef":  string(marshal),
	}

	videoInsertData := model.Video{
		BangumiId: request.BangumiId,
		Origin:    2,
		Category:  request.Category,
		Title:     request.Title,
		Season:    request.Season,
		Cover:     request.Cover,
		Total:     request.Total,
		RssUrl:    request.RssUrl,
		RuleName:  ruleName,
		PlayTime:  request.PlayTime,
		SourceDir: sourcePath,
		LinkDir:   linkDir,
		Time:      model.NewTime(),
	}
	ctx := context.TODO()
	ret := this.transaction(ctx, func(repo *repo.Repo) error {
		if err := repo.VideoInsertOne(ctx, videoInsertData); err != nil {
			return err
		}

		if request.Mode == 1 {
			if err := this.qb.AddRss(ctx, request.RssUrl, ruleName); err != nil {
				return err
			}
			if err := this.qb.SetRule(ctx, setRule); err != nil {
				return err
			}
			return nil
		}

		if err := this.qb.AddTorrents(ctx, request.TorrentList, category.Name, downloadPath); err != nil {
			return err
		}

		return nil
	})
	if ret.Error != nil {
		log.Printf(`[%s][第%d季] 添加番剧失败 [error] %v \n`, request.Title, request.Season, ret.Error)
		return ret.Error
	}
	if ret.TxError != nil {
		log.Printf(`[%s][第%d季] 添加番剧失败 [error] %v \n`, request.Title, request.Season, ret.TxError)
		return ret.TxError
	}

	return
}

// GetList 获取番剧列表
func (this *Service) GetList(request vo.VideoGetListRequest) (response vo.VideoGetListResponse, err error) {
	list, total, err := this.repo.VideoSelectList(context.TODO(), request)
	if err != nil {
		log.Println(err)
		return
	}
	response.Items = list
	response.Total = total
	return
}

// DeleteVideo 删除番剧
func (this *Service) DeleteVideo(request vo.DeleteVideoRequest) (err error) {

	ret := this.transaction(context.TODO(), func(repo *repo.Repo) error {
		// 获取番剧信息
		video, err := repo.VideoSelectOne(context.TODO(), request.Id, 0, "")
		if err != nil {
			return errors.New("获取")
		}

		// 删除表数据
		if err := repo.VideoDeleteList(context.TODO(), []int64{request.Id}); err != nil {
			return err
		}

		videoItemList, _, err := repo.VideoItemSelectList(context.TODO(), []int64{request.Id})

		// 删除番剧的视频数据
		if err := repo.VideoItemDeleteListByPid(context.TODO(), []int64{video.Id}); err != nil {
			return err
		}

		// 删除本地硬链接
		if err := os.RemoveAll(video.LinkDir); err != nil {
			log.Println(err)
			return err
		}

		// 删除linkpath缓存
		for _, item := range videoItemList {
			global.LinkPathCache.Delete(item.LinkPath)
		}

		// 删除本地种子
		if err := os.RemoveAll(video.SourceDir); err != nil {
			log.Println(err)
			return err
		}

		// 如果ruleName == "" 说明不是qb的规则任务
		if video.RuleName == "" {
			return nil
		}
		// 删除下载规则
		if err := this.qb.DeleteRule(context.TODO(), video.RuleName); err != nil {
			return err
		}
		// 删除rss
		if err := this.qb.DeleteRss(context.TODO(), video.RuleName); err != nil {
			return err
		}
		return nil
	})
	if ret.TxError != nil {
		err = ret.TxError
		log.Println(err)
		return
	}
	if ret.Error != nil {
		err = ret.Error
		return
	}

	return
}

// Link 执行硬链接
func (this *Service) Link() error {
	// 获取番剧列表
	list, total, err := this.repo.VideoSelectList(context.TODO(), vo.VideoGetListRequest{})
	if err != nil {
		log.Printf("[硬链接失败][获取视频列表失败][error]:%v\n", err)
		return err
	}

	start := time.Now()
	videoItemList := make([]model.VideoItem, 0, total)
	linkSourceList := make([]model.LinkSource, 0, total)

	for _, item := range list {
		// 创建媒体文件夹
		if err := os.MkdirAll(item.LinkDir, 0666); err != nil {
			log.Printf("[硬链接失败][创建媒体库路径失败] %s [已跳过该番剧][error]: %v\n", item.LinkDir, err)
			continue
		}

		// 检测下载路径是否存在
		exists, err := tool.IsExists(item.SourceDir)
		if err != nil {
			log.Printf("[硬链接失败][获取番剧路径信息失败] %s [已跳过该番剧][error]: %v\n", item.SourceDir, err)
			continue
		}
		if !exists.Exist {
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
			linkPath := filepath.Join(item.LinkDir, title)

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
		return nil
	}

	for _, item := range linkSourceList {
		// 判断下载文件是否存在
		exist, err := tool.IsExists(item.SourceFilePath)
		if err != nil {
			log.Printf("[硬链接失败][获取下载视频文件失败] %s [已跳过该文件][error]:%v\n", item.SourceFilePath, err)
			continue
		}
		if !exist.Exist {
			log.Printf("[硬链接失败][下载视频文件不存在] %s [已跳过该文件] \n ", item.SourceFilePath)
			continue
		}

		// 判断硬链接文件是否存在
		exist, err = tool.IsExists(item.LinkFilePath)
		if err != nil {
			log.Printf("[硬链接失败][尝试获取硬链接文件信息错误] %s [已跳过该文件][error]:%v\n", item.LinkFilePath, err)
			continue
		}
		if exist.Exist {
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
	cnt, err := this.repo.VideoItemInsertList(context.TODO(), videoItemList)
	if err != nil {
		log.Printf("[硬链接失败][插入数据库失败] %v \n", err)
		return err
	}

	end := time.Now()
	difference := end.Sub(start)
	log.Printf("[硬链接任务完成] [共计执行成功%d次硬链接] [共计耗时 %s ] \n", cnt, difference.String())
	return nil
}

func (this *Service) CreateCategory(request vo.CreateCategoryRequest) (err error) {
	ctx := context.TODO()
	ret := this.transaction(ctx, func(repo *repo.Repo) error {
		if err := repo.CategoryInsertOne(context.TODO(), model.Cagetory{
			Name:   request.Name,
			Origin: 2,
			Time:   model.NewTime(),
		}); err != nil {
			return err
		}

		if err := os.MkdirAll(this.getCategoryDir(request.Name), 0644); err != nil {
			return err
		}

		return nil
	})
	if ret.TxError != nil {
		err = ret.TxError
		log.Println(err)
		return
	}
	if ret.Error != nil {
		err = ret.Error
		return
	}
	return
}

func (this *Service) GetCategoryList() (res []model.Cagetory, err error) {
	return this.repo.CategorySelectList(context.TODO())
}

func (this *Service) GetCategory(id int64) (res model.Cagetory, err error) {
	return this.repo.CategorySelectOne(context.TODO(), id)
}

func (this *Service) GetCategoryByName(name string) (res model.Cagetory, err error) {
	return this.repo.CategorySelectOneByName(context.TODO(), name)
}

// DeleteCategory 删除分类
//
//	对应的分类移动到未知分类中
func (this *Service) DeleteCategory(id int64) (err error) {
	defer func() {
		go this.Link()
	}()

	ret := this.transaction(context.TODO(), func(repo *repo.Repo) error {
		// 获取分类信息
		category, err := repo.CategorySelectOne(context.TODO(), id)
		if err != nil {
			return err
		}
		if category.Origin == 1 {
			return errors.New("无法删除系统内置的分类")
		}
		oldCategoryName := category.Name

		// 获取新分类信息
		newCategory, err := repo.CategorySelectOne(context.TODO(), 1)
		if err != nil {
			return err
		}
		newCategoryName := newCategory.Name

		// 获取分类对应的番剧信息
		videoList, _, err := repo.VideoSelectList(context.TODO(), vo.VideoGetListRequest{Category: id})
		if err != nil {
			return err
		}

		// 获取番剧对应的视频信息
		pidList := make([]int64, 0, len(videoList))
		for _, item := range videoList {
			pidList = append(pidList, item.Id)
		}

		videoItemList, _, err := repo.VideoItemSelectList(context.TODO(), pidList)
		if err != nil {
			return err
		}

		// 删除番剧对应的视频
		if err := repo.VideoItemDeleteListByPid(context.TODO(), pidList); err != nil {
			return err
		}

		// 赋值旧分类文件夹路径与新分类文件夹路径
		oldCategoryDir := this.getCategoryDir(oldCategoryName)
		newCategoryDir := this.getCategoryDir(newCategoryName)

		for i := 0; i < len(videoList); i++ {
			videoList[i].Category = 1
			videoList[i].LinkDir = strings.Replace(videoList[i].LinkDir, oldCategoryDir, newCategoryDir, 1)
		}

		// 更新番剧 link_dir
		if err := repo.VideoUpdateCategoryList(context.TODO(), videoList); err != nil {
			return err
		}

		// 删除分类
		if err := repo.CategoryDeleteOne(context.TODO(), id); err != nil {
			return err
		}

		// 从硬链接缓存中删除原来的硬链接路径
		for _, item := range videoItemList {
			global.LinkPathCache.Delete(item.LinkPath)
		}

		// 删除旧目录
		if err := os.RemoveAll(oldCategoryDir); err != nil {
			return err
		}

		return nil
	})

	if ret.TxError != nil {
		err = ret.TxError
		log.Println(err)
		return
	}
	if ret.Error != nil {
		err = ret.Error
		return
	}
	return
}

// UpdateCategoryName 更新分类的名称
func (this *Service) UpdateCategoryName(request vo.UpdateCategoryNameRequest) (err error) {
	defer func() {
		go this.Link()
	}()

	ret := this.transaction(context.TODO(), func(repo *repo.Repo) error {
		// 检查分类是否存在
		find, err := repo.CategorySelectOneByName(context.TODO(), request.Name)
		if err != nil {
			return err
		}

		if find.Id == request.Id {
			return nil
		}

		if find.Id > 0 {
			err = errors.New("分类名称已存在")
			log.Println(err)
			return err
		}
		// 获取分类信息
		category, err := repo.CategorySelectOne(context.TODO(), request.Id)
		if err != nil {
			return err
		}
		oldCategoryName := category.Name

		// 获取分类对应的番剧信息
		videoList, _, err := repo.VideoSelectList(context.TODO(), vo.VideoGetListRequest{Category: request.Id})
		if err != nil {
			return err
		}

		// 获取番剧对应的视频信息
		pidList := make([]int64, 0, len(videoList))
		for _, item := range videoList {
			pidList = append(pidList, item.Id)
		}

		videoItemList, _, err := repo.VideoItemSelectList(context.TODO(), pidList)
		if err != nil {
			return err
		}

		// 删除番剧对应的视频
		if err := repo.VideoItemDeleteListByPid(context.TODO(), pidList); err != nil {
			return err
		}

		// 赋值旧分类文件夹路径与新分类文件夹路径
		oldCategoryDir := this.getCategoryDir(oldCategoryName)
		newCategoryDir := this.getCategoryDir(request.Name)

		for i := 0; i < len(videoList); i++ {
			videoList[i].LinkDir = strings.Replace(videoList[i].LinkDir, oldCategoryDir, newCategoryDir, 1)
		}

		// 更新番剧 link_dir
		if err := repo.VideoUpdateCategoryList(context.TODO(), videoList); err != nil {
			return err
		}

		// 更新分类名
		if err := repo.CategoryUpdateName(context.TODO(), request.Id, request.Name); err != nil {
			return err
		}

		// 从硬链接缓存中删除原来的硬链接路径
		for _, item := range videoItemList {
			global.LinkPathCache.Delete(item.LinkPath)
		}

		// 删除旧目录
		if err := os.RemoveAll(oldCategoryDir); err != nil {
			return err
		}

		return nil
	})

	if ret.TxError != nil {
		err = ret.TxError
		log.Println(err)
		return
	}
	if ret.Error != nil {
		err = ret.Error
		return
	}

	return
}

// UpdateVideoCategory 更新番剧的分类
// /Media/A/title/s01/*  /Media/A/title/s01 remove
//	/Media/B/title/s01/*
func (this *Service) UpdateVideoCategory(request vo.UpdateVideoCategoryRequest) (err error) {
	defer func() {
		go this.Link()
	}()

	ctx := context.TODO()
	ret := this.transaction(ctx, func(repo *repo.Repo) error {
		// 获取新分类
		category, err := repo.CategorySelectOne(ctx, request.Category)
		if err != nil {
			return err
		}

		// 获取要更新的番剧列表
		videoList, err := repo.VideoSelectListById(ctx, request.Ids)
		if err != nil {
			return err
		}

		// 构建更新数据
		removeDir := make([]string, 0, len(videoList))
		pids := make([]int64, 0, len(videoList))
		for i := 0; i < len(videoList); i++ {
			pids = append(pids, videoList[i].Id)
			removeDir = append(removeDir, videoList[i].LinkDir)
			videoList[i].Category = request.Category
			videoList[i].LinkDir = this.getLinkDir(category.Name, videoList[i].Title, videoList[i].Season)
		}

		// 批量更新
		if err := repo.VideoUpdateCategoryList(ctx, videoList); err != nil {
			return err
		}

		videoItemList, _, err := repo.VideoItemSelectList(ctx, pids)
		if err != nil {
			return err
		}

		// 批量删除番剧的视频
		if err := repo.VideoItemDeleteListByPid(context.TODO(), pids); err != nil {
			return err
		}

		// 清除linkpath缓存
		for _, item := range videoItemList {
			global.LinkPathCache.Delete(item.LinkPath)
		}

		// 删除原link目录
		for _, item := range removeDir {
			if err := os.RemoveAll(item); err != nil {
				log.Printf("[更新番剧分类失败][删除原硬链接目录失败] [%s] error: %v \n", item, err)
				return err
			}
		}
		return nil
	})
	if ret.TxError != nil {
		err = ret.TxError
		log.Println(err)
		return
	}
	if ret.Error != nil {
		err = ret.Error
		return
	}
	return
}
