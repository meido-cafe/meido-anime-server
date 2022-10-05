package response

import "github.com/gin-gonic/gin"

func Custom(ctx *gin.Context, code int, message string, data any) {
	res := make(gin.H)
	res["code"] = code
	if message != "" {
		res["message"] = message
	}
	if data != nil {
		res["data"] = data
	}
	ctx.JSON(200, res)
}
func Success(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(200, gin.H{
		"code":    500,
		"message": msg,
	})
}
func Bad(ctx *gin.Context, msg string) {
	ctx.JSON(200, gin.H{
		"code":    400,
		"message": msg,
	})
}

func Data(ctx *gin.Context, data any) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func List[T int | int64](ctx *gin.Context, list any, total T) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"items": list,
			"total": total,
		},
	})
}

func BadBind(ctx *gin.Context) {
	Bad(ctx, "参数错误")
}
