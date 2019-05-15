package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"project/lottery/project/conf"
	"sync"
)

// 连接到Mysql

var dbLock sync.Mutex
var masterInstance *xorm.Engine
var slaveInstance *xorm.Engine

func InstanceDBMaster() *xorm.Engine {
	if masterInstance != nil {
		return masterInstance
	}
	dbLock.Lock()
	defer dbLock.Unlock()

	// 第二次判断
	if masterInstance != nil {
		return masterInstance
	}

	return NewDBMaster()
}

func NewDBMaster() *xorm.Engine {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.DBMaster.User,
		conf.DBMaster.Pwd,
		conf.DBMaster.Host,
		conf.DBMaster.Port,
		conf.DBMaster.Database)

	instance, err := xorm.NewEngine(conf.DriverName, sourceName)
	if err != nil {
		log.Fatalln("dbHelper.InstanceDBMaster.NewEngine error: ", err)
		return nil
	}
	instance.ShowSQL(true)
	masterInstance = instance
	return masterInstance
}