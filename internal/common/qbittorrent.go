package common

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"meido-anime-server/etc"
	"meido-anime-server/internal/global"
	"sync"
	"time"
)

var qb *QB
var qbOnce sync.Once

type QB struct {
	Client *req.Client
}

func NewQB(conf *etc.Config) *QB {
	qbOnce.Do(func() {
		qb = new(QB)

		qb.Client = req.C().
			SetTimeout(5 * time.Second).
			SetBaseURL(fmt.Sprintf("%s/api/v2", conf.QB.Url))

		//if conf.Env == "dev" || conf.Env == "local" {
		//	qb.Client = qb.Client.DevMode()
		//}
		log.Printf("正在连接 qbittorrent: [url: %s] [username: %s] \n", conf.QB.Url, conf.QB.Username)
		res, err := qb.Client.R().SetFormDataAnyType(map[string]interface{}{
			"username": conf.QB.Username,
			"password": conf.QB.Password,
		}).Post("/auth/login")

		if err != nil {
			log.Println("qbittorrent 登录失败")
			panic(err)
		}
		if res.IsError() {
			log.Println("qbittorrent 登录失败")
			panic(res.String())
		}

		res, err = qb.Client.R().Get("/torrents/categories")
		if err != nil {
			log.Println(err)
			panic(err)
		}

		if res.IsError() {
			log.Println("获取分类信息失败")
			panic(res.String())
		}

		hash := make(map[string]struct{})
		if err := json.Unmarshal(res.Bytes(), &hash); err != nil {
			panic(err)
		}

		if _, ok := hash[global.QBCategory]; !ok {
			res, err = qb.Client.R().SetFormData(map[string]string{
				"category": global.QBCategory,
			}).Post("/torrents/createCategory")
			if err != nil {
				log.Println("qbittorrent 创建分类失败")
				panic(err)
			}
			if res.IsError() {
				log.Println("qbittorrent 创建分类失败")
				panic(res.String())
			}
		}

		log.Println("qbittorrent 连接成功")
	})
	return qb
}
