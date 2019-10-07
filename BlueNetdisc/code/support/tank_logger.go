package support

import (
	"LearnGo/BlueNetdisc/code/tool/util"
	"fmt"
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
	//this.Info()
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

	var consoleFormat = fmt.Sprintf("%s%s %s:%d %s", prefix, util.ConverTimeToTimeString(time.Now()),
		util.GetFilenameOfPath(file), line, content)
	fmt.Printf(consoleFormat)
	this.goLogger.SetPrefix(prefix)

	err := this.goLogger.Output(3, content)
	if err != nil {
		fmt.Printf("occur error while logging %s \r\n", err.Error())
	}
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
