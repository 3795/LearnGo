package util

import (
	"os"
	"os/user"
	"strings"
)

// 判断是否是Window下的开发环境
func EnvWinDevelopment() bool {
	// 获取执行路径
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	systemUser, err := user.Current()
	if systemUser != nil {
		return strings.HasPrefix(ex, systemUser.HomeDir+"\\AppData\\Local\\Temp")
	}

	return false
}

func EnvMacDevelopment() bool {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return strings.HasPrefix(ex, "/private/var/folders")
}

// 判断是否是开发环境（是否在IDE中运行）
func EnvDevlopment() bool {
	return EnvWinDevelopment() || EnvMacDevelopment()
}
