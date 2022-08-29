package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// DataBase 数据库模型
type DataBase struct {
	FocusGoDB string `yaml:"FocusGoDB"`
}

type App struct {
	//是否启动 性能分析功能 true:启动
	EnablePProf bool
}

// Redis redis缓存模型
type Redis struct {
	Url    string `yaml:"Url"`
	Prefix string `yaml:"Prefix"`
}

// rabbitmq配置
type RabbitmqConf struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
}

type Config struct {
	Redis    Redis    `yaml:"Redis"`
	RabbitmqConf RabbitmqConf `yaml:"RabbitmqConf"`
	DataBase DataBase `yaml:"DataBase"`
	LogConf  LogConf  `yaml:"LogConf"`
	Etcd     Etcd     `yaml:"Etcd"`
	App      App      `yaml:"App"`
}

type LogConf struct {
	FilePath string `yaml:"FilePath"`
}

type Etcd struct {
	GrpcProxy string `yaml:"GrpcProxy"`
}


var Conf = &Config{}

// SetUp 读取配置信息
func SetUp() {
	path := os.Getenv("LOC_CFG")
	if path == "" {
		panic("Invalid config path to load please use `LOC_CFG` set to os environment!")
	}

	if f, err := os.Open(path); err != nil {
		panic(fmt.Sprintf("Load config from path %s failed: %s", path, err.Error()))
	} else {
		_ = yaml.NewDecoder(f).Decode(Conf)
	}
}
