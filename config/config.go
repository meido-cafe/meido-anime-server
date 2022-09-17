package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"meido-anime-server/internal/global"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Env    string `yaml:"env"`
	Server struct {
		Port                int    `yaml:"port"`
		GinMode             string `yaml:"gin_mode"`
		TokenExpirationTime int64  `yaml:"token_expiration_time"` // token过期时间
		MediaPath           string `yaml:"media_path"`            // 媒体目录 (硬链接目录)
		SourcePath          string `yaml:"source_path"`           // 资源目录 (下载目录)
	} `yaml:"server"`
	Db struct {
		Path    string `yaml:"path"`
		MaxCons int    `yaml:"max_cons"`
	} `yaml:"database"`
	QB struct {
		Url          string `yaml:"url"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		Category     string `yaml:"category"`
		DownloadPath string `yaml:"download_path"` // QB下载目录 (QB下载文件时指定的路径)
	} `yaml:"qbittorrent"`
}

var Conf Config
var configOnce sync.Once

func marshal(filename string) {
	var err error

	byt, err := ioutil.ReadFile(filepath.Join("config", filename+".yaml"))
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(byt, &Conf); err != nil {
		panic(err)
	}
}

func handleLocal() {
	marshal("local")
	global.QBCategory = Conf.QB.Category
	global.QBDownloadPath = Conf.QB.DownloadPath
	global.SourcePath = Conf.Server.SourcePath
	global.MediaPath = Conf.Server.MediaPath
}

// dev环境config配置
func handleDev() {
	marshal("dev")

	Conf.QB.Username = os.Getenv("QB_USERNAME")
	Conf.QB.Password = os.Getenv("QB_PASSWORD")
	Conf.QB.Url = os.Getenv("QB_WEB_URL")
	global.QBCategory = os.Getenv("QB_CATEGORY")
	global.QBDownloadPath = os.Getenv("QB_DOWNLOAD_PATH")
	global.SourcePath = os.Getenv("SOURCE_PATH")
	global.MediaPath = os.Getenv("MEDIA_PATH")
}

func handlePro() {
	marshal("pro")

	Conf.QB.Username = os.Getenv("QB_USERNAME")
	Conf.QB.Password = os.Getenv("QB_PASSWORD")
	Conf.QB.Url = os.Getenv("QB_WEB_URL")
	global.QBCategory = os.Getenv("QB_CATEGORY")
	global.QBDownloadPath = os.Getenv("QB_DOWNLOAD_PATH")
	global.MediaPath = os.Getenv("MEDIA_PATH")
	global.SourcePath = os.Getenv("SOURCE_PATH")
}

func InitConfig() {
	configOnce.Do(func() {
		log.Println("初始化配置")
		marshal("common")
		switch Conf.Env {
		case "local":
			handleLocal()
		case "dev":
			handleDev()
		case "pro":
			handlePro()
		default:
			panic(errors.New("错误的环境类型"))
		}
		log.Println("配置初始化完成")
	})
}

func NewConfig() *Config {
	InitConfig()
	return &Conf
}
