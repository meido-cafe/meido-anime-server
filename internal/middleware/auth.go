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
	const CODE = 401
	return func(ctx *gin.Context) {

		// token为空
		token := ctx.GetHeader("token")
		if token == "" {
			response.Custom(ctx, CODE, "用户未登录", nil)
			ctx.Abort()
			return
		}

		// token不存在
		user, ok := global.TokenCache[token]
		if !ok {
			response.Custom(ctx, CODE, "token不存在", nil)
			ctx.Abort()
			return
		}

		nowUnix := time.Now().Unix()
		// token 过期
		if user.TokenTime+this.conf.Server.TokenExpirationTime < nowUnix {
			delete(global.TokenCache, token)
			response.Custom(ctx, CODE, "token已失效,请重新登录", nil)
			ctx.Abort()
			return
		}

		// 刷新token时间
		global.TokenCache[token].TokenTime = nowUnix

		ctx.Next()
	}
}
