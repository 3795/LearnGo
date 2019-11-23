package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort               int      `json:"apiPort"`
	ApiReadTimeout        int      `json:"apiReadTimeout"`
	ApiWriteTimeout       int      `json:"apiWriteTimeout"`
	EtcdEndpoints         []string `json:"etcdEndpoints"`
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`
	WebRoot               string   `json:"webRoot"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
}

// 配置句柄单例
var (
	G_config *Config
)

// 加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	// 读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// Json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	// 单例赋值
	G_config = &conf
	return
}
