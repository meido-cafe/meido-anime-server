package config

import (
	"github.com/gin-gonic/gin"
	"log"
	"meido-anime-server/internal/global"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Server struct {
		Port                int    `yaml:"port"`                  // 服务端口
		GinMode             string `yaml:"gin_mode"`              // gin框架的模式
		TokenExpirationTime int64  `yaml:"token_expiration_time"` // token过期时间
		MediaPath           string `yaml:"media_path"`            // 媒体目录 (硬链接目录)
		SourcePath          string `yaml:"source_path"`           // 资源目录 (下载目录)
	}
	Db struct {
		Path    string `yaml:"path"`     // sqlite 路径
		MaxCons int    `yaml:"max_cons"` // 最大连接数
	}
	QB struct {
		Url          string `yaml:"url"`           // qb的web url
		Username     string `yaml:"username"`      // qb的web 用户名
		Password     string `yaml:"password"`      // qb的web 密码
		Category     string `yaml:"category"`      // qb中的下载分类
		DownloadPath string `yaml:"download_path"` // QB下载目录 (QB下载文件时指定的路径)
	}
}

var Conf Config
var configOnce sync.Once

func InitConfig() {
	configOnce.Do(func() {
		log.Println("初始化配置")

		serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
		if err != nil {
			log.Fatalln("初始化配置失败, error:", err)
		}

		tokenTime, err := strconv.ParseInt(os.Getenv("SERVER_TOKEN_EXPIRED_TIME"), 10, 64)
		if err != nil {
			log.Fatalln("初始化配置失败, error:", err)
		}

		Conf.Server.Port = serverPort
		Conf.Server.TokenExpirationTime = tokenTime
		Conf.Server.MediaPath = os.Getenv("MEDIA_PATH")
		Conf.Server.SourcePath = os.Getenv("SOURCE_PATH")

		Conf.Db.Path = os.Getenv("DB_PATH")
		dbMaxCons, err := strconv.Atoi(os.Getenv("DB_MAX_CONS"))
		if err != nil {
			log.Fatalln("初始化配置失败, error:", err)
		}
		Conf.Db.MaxCons = dbMaxCons

		Conf.QB.Username = os.Getenv("QB_USERNAME")
		Conf.QB.Password = os.Getenv("QB_PASSWORD")
		Conf.QB.Url = os.Getenv("QB_WEB_URL")
		Conf.QB.Category = os.Getenv("QB_CATEGORY")
		Conf.QB.DownloadPath = os.Getenv("QB_DOWNLOAD_PATH")

		switch os.Getenv("_MODE") {
		case "local":
			Conf.Server.GinMode = gin.DebugMode
		case "dev":
			Conf.Server.GinMode = gin.DebugMode
		case "pro_test":
			Conf.Server.GinMode = gin.ReleaseMode
		case "pro":
			Conf.Server.GinMode = gin.ReleaseMode
		}

		global.QBCategory = os.Getenv("QB_CATEGORY")
		global.QBDownloadPath = os.Getenv("QB_DOWNLOAD_PATH")
		global.SourcePath = os.Getenv("SOURCE_PATH")
		global.MediaPath = os.Getenv("MEDIA_PATH")

		log.Println("配置初始化完成")
	})
}

func GetConfig() *Config {
	InitConfig()
	return &Conf
}
