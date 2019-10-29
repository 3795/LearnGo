package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort               int      `json:"api_port"`
	ApiReadTimeout        int      `json:"api_read_timeout"`
	ApiWriteTimeout       int      `json:"api_write_timeout"`
	EtcdEndpoints         []string `json:"etcd_endpoints"`
	EtcdDialTimeout       int      `json:"etcd_dial_timeout"`
	WebRoot               string   `json:"web_root"`
	MongodbUrl            string   `json:"mongodb_url"`
	MongodbConnectTimeout int      `json:"mongodb_connect_timeout"`
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
