package service

import "log"

type InitService struct {
	CronService *CronService
	service     *Service
}

func NewInitService() *InitService {
	return &InitService{
		service:     NewService(),
		CronService: NewCronService(),
	}
}

func (this *InitService) Init() {
	if err := this.service.CacheLinkPath(); err != nil {
		log.Println("缓存硬链接资源路径失败:", err)
	}

	defer this.CronService.Start()
}
