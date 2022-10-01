package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"meido-anime-server/internal/global"
	"meido-anime-server/internal/tool"
	"path/filepath"
	"reflect"
	"runtime"
)

type CronService struct {
	cron    *cron.Cron
	service *Service
}

func NewCronService() *CronService {
	return &CronService{
		cron:    cron.New(),
		service: NewService(),
	}
}

func (this *CronService) Start() {

	list := []cronFunc{
		this.handleVideoLink,
		this.handleClearDir,
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
		this.service.Link()
	}}
}

// 清理媒体文件夹下的空文件夹
func (this *CronService) handleClearDir() producer {
	return producer{"0 */10 * * * ?", func() {
		list, err := this.service.GetCategoryList()
		if err != nil {
			log.Println("[定时清理媒体库空目录失败] [获取分类列表错误]:", err)
			return
		}
		for _, item := range list {
			go func(dir string) {
				ret, err := tool.RemoveEmptyDirAll(dir, false)
				if err != nil {
					log.Printf("[定时清理空目录][失败] [%s] error: %v\n", dir, err)
					return
				}
				for _, item := range ret {
					log.Printf("[定时清理空目录][成功] 检测到 [%s] 为空目录,已删除\n", item)
				}
			}(filepath.Join(global.MediaPath, item.Name))
		}
	}}
}
