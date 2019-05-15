package main

import (
	"fmt"
	"project/lottery/project/bootstrap"
	"project/lottery/project/web/route"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("Go抽奖系统", "Seven.iwi")
	app.Bootstrap()
	app.Configure(route.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
