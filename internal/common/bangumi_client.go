package common

import (
	"github.com/imroc/req/v3"
	"sync"
	"time"
)

var bangumiClient *req.Client
var bangumiClientOnce sync.Once

func InitBangumiClient() {
	bangumiClientOnce.Do(func() {
		bangumiClient = req.C().
			SetTimeout(10 * time.Second).
			SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1370.37").
			SetBaseURL("https://api.bgm.tv")
	})
}
func GetBangumiClient() *req.Client {
	InitBangumiClient()
	return bangumiClient
}
