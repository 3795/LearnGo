package core

import (
	"LearnGo/BlueNetdisc/code/tool/cache"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Context interface {
	http.Handler

	GetDB() *gorm.DB

	GetBean(bean Bean) Bean

	GetSessionCache() *cache.Table

	GetControllerMap() map[string]Controller

	InstallOk()

	CleanUp()
}
