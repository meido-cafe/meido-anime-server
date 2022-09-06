package service

import "log"

type InitService struct {
	CronService  *CronService
	videoService *VideoService
}

func NewInitService(videoService *VideoService, cronService *CronService) *InitService {
	return &InitService{
		videoService: videoService,
		CronService:  cronService,
	}
}

func (this *InitService) Init() {
	if err := this.videoService.CacheSourcePath(); err != nil {
		log.Println("缓存硬链接资源路径失败:", err)
	}

	defer this.CronService.Start()
}
