package common

import (
	"github.com/imroc/req/v3"
	"sync"
	"time"
)

var bangumiClient *req.Client
var bangumiClientOnce sync.Once

func NewBangumiClient() *req.Client {
	bangumiClientOnce.Do(func() {
		bangumiClient = req.C().
			SetTimeout(5 * time.Second).
			SetBaseURL("https://api.bgm.tv")
	})
	return bangumiClient
}
