package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/service"
)

type DemoHandler struct {
	service *service.DemoService
}

func NewDemoHandler(service *service.DemoService) *DemoHandler {
	return &DemoHandler{service: service}
}

func (this DemoHandler) Hello(ctx *gin.Context) {
	this.service.Hello()
	ctx.JSON(200, gin.H{
		"msg": "world",
	})
}
