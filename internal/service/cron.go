package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"reflect"
	"runtime"
)

type CronService struct {
	cron         *cron.Cron
	videoService *VideoService
}

func NewCronService(videoService *VideoService) *CronService {
	return &CronService{
		cron:         cron.New(),
		videoService: videoService,
	}
}

func (this *CronService) Start() {

	list := []cronFunc{
		this.handleVideoLink,
	}

	if err := this.register(list...); err != nil {
		log.Fatalln("定时任务注册失败:", err)
	}
	go this.cron.Run()
}

type producer struct {
	Spec string
	f    func()
}

type cronFunc func() producer

func (this *CronService) register(list ...cronFunc) (err error) {
	var spec string
	var f func()
	for _, item := range list {
		spec = item().Spec
		f = item().f
		if err = this.cron.AddFunc(spec, f); err != nil {
			fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			msg := fmt.Sprintf("[%s] [%s] 创建定时任务失败:%v\n", spec, fname, err)
			err = errors.New(msg)
			return
		}
	}
	return
}

func (this *CronService) handleVideoLink() producer {
	return producer{"0 */10 * * * ?", func() {
		this.videoService.Link()
	}}
}
