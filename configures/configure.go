package configures

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type ImConfig struct {
	Log struct {
		LogPath string `yaml:"logPath"`
		LogName string `yaml:"logName"`
	} `ymal:"log"`

	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		DbName   string `yaml:"name"`
	} `yaml:"mysql"`

	Wx struct {
		AppId  string `yaml:"appid"`
		Secret string `yaml:"secret"`
	} `yaml:"wx"`
}

var Config ImConfig
var Env string

func InitConfigures() error {
	cfBytes, err := ioutil.ReadFile("conf/config.yml")
	if err == nil {
		var conf ImConfig
		yaml.Unmarshal(cfBytes, &conf)
		Config = conf
		return nil
	} else {
		return err
	}
}
