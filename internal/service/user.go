package service

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"meido-anime-server/internal/global"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
	"meido-anime-server/internal/tool"
	"os"
	"time"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (this *UserService) Login(req vo.LoginRequest, user model.User) (token string, err error) {
	md5Username := tool.MD5Salt(os.Getenv("USERNAME"), global.Salt, 1)
	if req.Username != md5Username {
		err = errors.New("用户名或密码错误")
		return
	}

	md5pwd := tool.MD5Salt(os.Getenv("PASSWORD"), global.Salt, 2)
	if req.Password != md5pwd {
		err = errors.New("用户名或密码错误")
		return
	}

	// 生成token
	token = uuid.New().String()
	// 缓存token
	now := time.Now()
	user.TokenTime = now.Unix()
	user.LoginTime = now.Unix()
	global.TokenCache[token] = &user

	log.Printf("[%s][%s][%s] login \n", now.Format("2006-01-02 15:04:05"), user.Ip, user.UserAgent)
	return
}

func (this *UserService) Logout(token string) {
	delete(global.TokenCache, token)
}
