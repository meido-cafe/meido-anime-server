package v1

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/service"
)

type UserApi struct {
	service *service.UserService
}

func NewUserApi(userService *service.UserService) *UserApi {
	return &UserApi{service: userService}
}

func (this *UserApi) Login(ctx *gin.Context) {
	req := vo.LoginRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	if req.Username == "" || req.Password == "" {
		response.Bad(ctx, "用户名或密码不能为空")
		return
	}

	token, err := this.service.Login(req, model.User{
		UserAgent: ctx.GetHeader("User-Agent"),
		Ip:        ctx.ClientIP(),
	})
	if err != nil {
		response.Bad(ctx, err.Error())
		return
	}

	response.Data(ctx, gin.H{
		"token": token,
	})
}

func (this *UserApi) Logout(ctx *gin.Context) {
	this.service.Logout(ctx.GetHeader("token"))
	response.Success(ctx)
}
