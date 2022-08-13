package service

import (
	"context"
	"log"
	"meido-anime-server/internal/repo"
)

type DemoService struct {
	repo repo.IDemoRepo
}

func NewDemoService(repo repo.IDemoRepo) *DemoService {
	return &DemoService{repo: repo}
}

func (this DemoService) Hello() {
	err := this.repo.Hello(context.TODO())
	if err != nil {
		log.Println(err)
		return
	}
}
