package common

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"meido-anime-server/config"
	"meido-anime-server/internal/global"
	"sync"
	"time"
)

var qb *QB
var qbOnce sync.Once

type QB struct {
	Client *req.Client
}

func InitQB() {
	qbOnce.Do(func() {
		qb = new(QB)

		qb.Client = req.C().
			SetTimeout(5 * time.Second).
			SetBaseURL(fmt.Sprintf("%s/api/v2", config.Conf.QB.Url))

		//if config.Conf.Env == "dev" || config.Conf.Env == "local" {
		//	qb.Client = qb.Client.DevMode()
		//}
		log.Printf("正在连接 qbittorrent: [url: %s] [username: %s] \n", config.Conf.QB.Url, config.Conf.QB.Username)
		res, err := qb.Client.R().SetFormDataAnyType(map[string]interface{}{
			"username": config.Conf.QB.Username,
			"password": config.Conf.QB.Password,
		}).Post("/auth/login")

		if err != nil {
			log.Fatalln("qbittorrent 登录失败:", err)
		}
		if res.IsError() {
			log.Fatalln("qbittorrent 登录失败:", res.String())
		}

		res, err = qb.Client.R().Get("/torrents/categories")
		if err != nil {
			log.Fatalln("qbittorrent 获取分类信息失败:", err)
		}

		if res.IsError() {
			log.Fatalln("qbittorrent 获取分类信息失败:", res.String())
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
				log.Fatalln("qbittorrent 创建分类失败:", err)
			}
			if res.IsError() {
				log.Fatalln("qbittorrent 创建分类失败:", res.String())
			}
		}

		log.Println("qbittorrent 连接成功")
	})
}
func GetQB() *QB {
	InitQB()
	return qb
}
