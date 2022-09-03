package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"meido-anime-server/internal/repo"
	"reflect"
	"runtime"
)

type CronService struct {
	cron *cron.Cron
	repo *repo.CronRepo
}

func NewCronService(cronRepo *repo.CronRepo) *CronService {
	return &CronService{
		cron: cron.New(),
		repo: cronRepo,
	}
}

func (this *CronService) Start() {
	if err := this.register(); err != nil {
		log.Fatalln(err)
	}
	this.cron.Run()
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
