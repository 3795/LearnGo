package support

import (
	"LearnGo/BlueNetdisc/code/core"
	"github.com/json-iterator/go"
	"time"
	"unsafe"
)

type TankConfig struct {
	serverPort int
	// 是否安装
	installed bool
	// 文件存放的位置
	matterPath string
	// Mysql地址
	mysqlUrl string
	// 其他配置参数
	item *ConfigItem
}

// tank.json config items.
type ConfigItem struct {
	ServerPort    int
	MatterPath    string
	MysqlPort     int
	MysqlHost     string
	MysqlSchema   string
	MysqlUsername string
	MysqlPassword string
}

//validate whether the config file is ok
func (this *ConfigItem) validate() bool {

	if this.ServerPort == 0 {
		core.LOGGER.Error("ServerPort is not configured")
		return false
	}

	if this.MysqlUsername == "" {
		core.LOGGER.Error("MysqlUsername  is not configured")
		return false
	}

	if this.MysqlPassword == "" {
		core.LOGGER.Error("MysqlPassword  is not configured")
		return false
	}

	if this.MysqlHost == "" {
		core.LOGGER.Error("MysqlHost  is not configured")
		return false
	}

	if this.MysqlPort == 0 {
		core.LOGGER.Error("MysqlPort  is not configured")
		return false
	}

	if this.MysqlSchema == "" {
		core.LOGGER.Error("MysqlSchema  is not configured")
		return false
	}
	return true
}

func (this *TankConfig) Init() {
	jsoniter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", iter.ReadString(), time.Local)
		if err != nil {
			iter.Error = err
			return
		}
		*((*time.Time)(ptr)) = t
	})

	jsoniter.RegisterTypeEncoderFunc("time.Time", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		t := *((*time.Time)(ptr))
		//if use time.UTC there will be 8 hours gap.
		stream.WriteString(t.Local().Format("2006-01-02 15:04:05"))
	}, nil)

	this.serverPort = core.DEFAULT_SERVER_PORT
	this.ReadFromConfigFile()
}

func (this *TankConfig) ReadFromConfigFile() {
	// read from tank.json
}
