package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/global"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/tool"
	"os"
	"time"
)

type UserApi struct {
}

func NewUserApi() *UserApi {
	return &UserApi{}
}

func (this *UserApi) Login(ctx *gin.Context) {
	req := vo.LoginRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadBind(ctx)
		return
	}

	username := os.Getenv("USERNAME")
	md5Username := tool.MD5Salt(username, global.Salt, 1)
	if req.Username != md5Username {
		log.Println(md5Username)
		response.Bad(ctx, "用户名或密码错误")
		return
	}

	password := os.Getenv("PASSWORD")
	md5pwd := tool.MD5Salt(password, global.Salt, 2)
	if req.Password != md5pwd {
		log.Println(md5pwd)
		response.Bad(ctx, "用户名或密码错误")
		return
	}

	// 生成token
	token := uuid.New().String()

	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")
	unix := now.Unix()
	userAgent := ctx.Request.Header.Get("User-Agent")
	ip := ctx.ClientIP()
	// 缓存token
	global.TokenCache[token] = &model.User{
		LoginTime: unix,
		TokenTime: unix,
		UserAgent: userAgent,
		Ip:        ip,
	}

	log.Printf("[%s][%s][%s] login \n", nowString, ip, userAgent)

	response.Data(ctx, gin.H{
		"token": token,
	})
}

func (this *UserApi) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	delete(global.TokenCache, token)
	response.Success(ctx)
}

func (this *UserApi) UpdateUsername(ctx *gin.Context) {

}
