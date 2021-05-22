package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Cfg Config

type DbConfig struct {
	UserName string
	PassWord string
	Host 	string
	Port	int
	DbName	string
}

type SettingConfig struct {
	Pn int
	Port int
}

type Config struct {
	Version string
	Author string
	Db DbConfig
	Settings SettingConfig
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func init()  {
	JsonParse := NewJsonStruct()
	Cfg = Config{}
	JsonParse.load("./config.json", &Cfg)
	fmt.Println("Config Init")
	fmt.Println(Cfg.Author, Cfg.Version)
}

func (jst *JsonStruct) load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func InitCfg() Config {
	return Cfg
}
