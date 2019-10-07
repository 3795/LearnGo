package support

import (
	"LearnGo/BlueNetdisc/code/core"
	"LearnGo/BlueNetdisc/code/tool/result"
	"LearnGo/BlueNetdisc/code/tool/util"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
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

	lowerMode := strings.ToLower(this.mode)
	if this.mode == "" || lowerMode == MODE_WEB {
		// 以web方式启动
		this.HandleWeb()
	} else if lowerMode == MODE_VERSION {
		this.HandleVersion()
	} else {
		if this.host == "" {
			this.host = fmt.Sprintf("http://127.0.0.1:%d", core.DEFAULT_SERVER_PORT)
		}

		if this.username == "" {
			panic(result.BadRequest("in mode %s, username is required", this.mode))
		}

		if this.password == "" {
			if util.EnvDevlopment() {
				panic(result.BadRequest("If run in IDE, use -password yourPassword to input password"))
			} else {
				fmt.Print("Enter Password:")
				bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					panic(err)
				}

				this.password = string(bytePassword)
				fmt.Println()
			}
		}

		if lowerMode == MODE_MIRROR {
			this.HandleMirror()
		} else if lowerMode == MODE_MIGRATE_20_TO_30 {
			this.HandleMigrate20to30()
		} else if lowerMode == MODE_CRAWL {
			this.HandleCrawl()
		} else {
			panic(result.BadRequest("cannot handle mode %s \r\n", this.mode))
		}
	}
}

func (this *TankApplication) HandleWeb() {
	//tankLogger := &TankLogger{}
}

func (this *TankApplication) HandleVersion() {

}

func (this *TankApplication) HandleMirror() {

}

func (this *TankApplication) HandleMigrate20to30() {

}

func (this *TankApplication) HandleCrawl() {

}
