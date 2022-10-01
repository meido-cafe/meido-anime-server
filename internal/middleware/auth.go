package middleware

import (
	"github.com/gin-gonic/gin"
	"meido-anime-server/config"
	"meido-anime-server/internal/api/response"
	"meido-anime-server/internal/global"
	"time"
)

type Middleware struct {
	conf *config.Config
}

func NewMiddleware() *Middleware {
	return &Middleware{conf: config.GetConfig()}
}

func (this *Middleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// token为空
		token := ctx.GetHeader("token")
		if token == "" {
			response.Bad(ctx, "未登录")
			ctx.Abort()
			return
		}

		// token不存在
		user, ok := global.TokenCache[token]
		if !ok {
			response.Bad(ctx, "非法token")
			ctx.Abort()
			return
		}

		nowUnix := time.Now().Unix()
		// token 过期
		if user.TokenTime+this.conf.Server.TokenExpirationTime < nowUnix {
			delete(global.TokenCache, token)
			response.Bad(ctx, "token已失效,请重新登录")
			ctx.Abort()
			return
		}

		// 刷新token时间
		global.TokenCache[token].TokenTime = nowUnix

		ctx.Next()
	}
}
