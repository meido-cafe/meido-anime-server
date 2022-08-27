package etc

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Env    string `yaml:"env"`
	Server struct {
		Port    int    `yaml:"port"`
		GinMode string `yaml:"gin_mode"`
	} `yaml:"server"`
	Db struct {
		Path    string `yaml:"path"`
		MaxCons int    `yaml:"max_cons"`
	} `yaml:"database"`
	QB struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
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
}

// dev环境config配置
func handleDev() {
	marshal("dev")
	Conf.QB.Username = os.Getenv("QB_USERNAME")
	Conf.QB.Password = os.Getenv("QB_PASSWORD")
	Conf.QB.Url = os.Getenv("QB_WEB_URL")
}

func handlePro() {
	marshal("pro")
	Conf.QB.Username = os.Getenv("QB_USERNAME")
	Conf.QB.Password = os.Getenv("QB_PASSWORD")
	Conf.QB.Url = os.Getenv("QB_WEB_URL")
}

func InitConfig() {
	configOnce.Do(func() {
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
	})
}

func NewConfig() *Config {
	InitConfig()
	return &Conf
}
