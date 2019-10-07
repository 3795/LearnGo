package util

import (
	"LearnGo/BlueNetdisc/code/tool/result"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(result.BadRequest(err.Error()))
		}
	}
}

func GetGoPath() string {
	return build.Default.GOPATH
}

func GetDevHomePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get dev home path")
	}

	dir := GetDirOfPath(file)
	return GetDirOfPath(dir)
}

func GetHomePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	if EnvDevlopment() {
		exPath = GetDevHomePath() + "/tmp"
	}

	return UniformPath(exPath)
}

/**
规范返回的路径格式
*/
func UniformPath(p string) string {
	p = strings.Replace(p, "\\", "/", -1)
	p = path.Clean(p)
	return strings.TrimSuffix(p, "/")
}

func GetDirOfPath(fullPath string) string {
	// Linux或Mac处理方式
	macIndex := strings.LastIndex(fullPath, "/")
	// Window处理方式
	winIndex := strings.LastIndex(fullPath, "\\")
	index := macIndex
	if winIndex > index {
		index = winIndex
	}

	return fullPath[:index]
}

func GetLogPath() string {
	homePath := GetHomePath()
	return homePath
}
