package support

import (
	"flag"
	"fmt"
)

const (
	// 以web的方式启动
	MODE_WEB = "web"
	// 将本地文件发送到云盘
	MODE_MIRROR = "mirror"
	// crawl remote file to EyeblueTank
	MODE_CRAWL = "crawl"
	// 当前版本
	MODE_VERSION = "version"
	// 将2.0版本迁移到3.0
	MODE_MIGRATE_20_TO_30 = "migrate20to30"
)

type TankApplication struct {
	mode      string
	host      string
	username  string
	password  string
	src       string
	dest      string
	overwrite bool
	filename  string
}

func (this *TankApplication) Start() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR: %v\r\n", err)
		}
	}()

	modePtr := flag.String("mode", this.mode, "cli mode web/mirror/crawl")
	hostPtr := flag.String("host", this.host, "tank host")
	usernamePtr := flag.String("username", this.username, "username")
	passwordPtr := flag.String("password", this.password, "password")
	srcPtr := flag.String("src", this.src, "src absolute path")
	destPtr := flag.String("dest", this.dest, "destination path in tank")
	overwritePtr := flag.Bool("overwriter", this.overwrite, "whether same file overwrite")
	filenamePtr := flag.String("filename", this.filename, "filename when crawl")

	flag.Parse()

	this.mode = *modePtr
	this.host = *hostPtr
	this.username = *usernamePtr
	this.password = *passwordPtr
	this.src = *srcPtr
	this.dest = *destPtr
	this.overwrite = *overwritePtr
	this.filename = *filenamePtr
}
