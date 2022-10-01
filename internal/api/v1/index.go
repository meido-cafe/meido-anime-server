package v1

import "meido-anime-server/internal/service"

type Api struct {
	service *service.Service
}

func NewApi() *Api {
	return &Api{service: service.NewService()}
}
