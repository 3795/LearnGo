package support

import (
	"LearnGo/BlueNetdisc/code/core"
	"LearnGo/BlueNetdisc/code/tool/util"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type TankLogger struct {
	sync.Mutex

	goLogger *log.Logger

	file *os.File
}

func (this *TankLogger) Init() {
	this.openFile()
	this.Info("[cron job] Every day 00:00 maintain log file")
	expression := "0 0 0 * * ?"
	cronJob := cron.New()
	err := cronJob.AddFunc(expression, this.maintain)
	core.PanicError(err)
	cronJob.Start()

}

func (this *TankLogger) Destroy() {
	this.closeFile()
}

func (this *TankLogger) Log(prefix string, format string, v ...interface{}) {
	content := fmt.Sprintf(format+"\r\n", v...)

	// print to console with line number
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	var consoleFormat = fmt.Sprintf("%s%s %s:%d %s", prefix, util.ConvertTimeToTimeString(time.Now()),
		util.GetFilenameOfPath(file), line, content)
	fmt.Printf(consoleFormat)
	this.goLogger.SetPrefix(prefix)

	err := this.goLogger.Output(3, content)
	if err != nil {
		fmt.Printf("occur error while logging %s \r\n", err.Error())
	}
}

func (this *TankLogger) Debug(format string, v ...interface{}) {
	this.Log("[DEBUG]", format, v...)
}

func (this *TankLogger) Info(format string, v ...interface{}) {
	this.Log("[INFO ]", format, v...)
}

func (this *TankLogger) Warn(format string, v ...interface{}) {
	this.Log("[WARN ]", format, v...)
}

func (this *TankLogger) Error(format string, v ...interface{}) {
	this.Log("[ERROR]", format, v...)
}

func (this *TankLogger) Panic(format string, v ...interface{}) {
	this.Log("[PANIC]", format, v...)
	panic(fmt.Sprintf(format, v...))
}

func (this *TankLogger) openFile() {
	fmt.Printf("use log file %s \r\n", this.fileName())
	f, err := os.OpenFile(this.fileName(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("cannot open log file" + err.Error())
	}

	this.goLogger = log.New(f, "", log.Ltime|log.Lshortfile)

	if this.goLogger == nil {
		fmt.Printf("Error: cannot create goLogger \r\n")
	}

	this.file = f
}

func (this *TankLogger) closeFile() {
	if this.file != nil {
		err := this.file.Close()
		if err != nil {
			panic("occur error while closing log file:" + err.Error())
		}
	}
}

func (this *TankLogger) fileName() string {
	return util.GetLogPath() + "/tank.log"
}

func (this *TankLogger) maintain() {
	this.Info("maintain log")

	this.Lock()
	defer this.Unlock()

	// 关闭文件流，为重命名做准备
	this.closeFile()

	destPath := util.GetLogPath() + "/tank-" + util.ConvertTimeToDateString(util.Yesterday()) + ".log"

	err := os.Rename(this.fileName(), destPath)
	if err != nil {
		this.Error("occur error while renaming log file %s %s", destPath, err.Error())
	}

	// 重新创建一个名叫tank.log的文件，记录新一天的日志
	this.openFile()

	// 删除一个月前的今天的日志文件
	monthAgo := time.Now()
	monthAgo = monthAgo.AddDate(0, -1, 0)
	oldDestPath := util.GetLogPath() + "/tank-" + util.ConvertTimeToDateString(monthAgo) + ".log"
	this.Log("try to delete log file %s", oldDestPath)

	exists := util.PathExists(oldDestPath)
	if exists {
		err = os.Remove(oldDestPath)
		if err != nil {
			this.Error("occur error while deleting log file %s %s", oldDestPath, err.Error())
		}
	} else {
		this.Error("log file %s not exists, skip", oldDestPath)
	}
}
