package etc

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
		Port      int    `yaml:"port"`
		GinMode   string `yaml:"gin_mode"`
		MediaPath string `yaml:"media_path"`
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
		DownloadPath string `yaml:"download_path"`
	} `yaml:"qbittorrent"`
}

var Conf Config
var configOnce sync.Once

func marshal(filename string) {
	var err error

	byt, err := ioutil.ReadFile(filepath.Join("etc", filename+".yaml"))
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(byt, &Conf); err != nil {
		panic(err)
	}
}

func handleLocal() {
	marshal("local")
	global.QBDownloadPath = Conf.QB.DownloadPath
	global.QBCategory = Conf.QB.Category
	global.MediaPath = Conf.Server.MediaPath
}

// dev环境config配置
func handleDev() {
	marshal("dev")

	Conf.QB.Username = os.Getenv("QB_USERNAME")
	Conf.QB.Password = os.Getenv("QB_PASSWORD")
	Conf.QB.Url = os.Getenv("QB_WEB_URL")
	global.QBDownloadPath = os.Getenv("QB_DOWNLOAD_PATH")
	global.QBCategory = os.Getenv("QB_CATEGORY")
	global.MediaPath = os.Getenv("MEDIA_PATH")
}

func handlePro() {
	marshal("pro")

	Conf.QB.Username = os.Getenv("QB_USERNAME")
	Conf.QB.Password = os.Getenv("QB_PASSWORD")
	Conf.QB.Url = os.Getenv("QB_WEB_URL")
	global.QBDownloadPath = os.Getenv("QB_DOWNLOAD_PATH")
	global.QBCategory = os.Getenv("QB_CATEGORY")
	global.MediaPath = os.Getenv("MEDIA_PATH")
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
