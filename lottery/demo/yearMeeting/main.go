package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var userList []string
var mu sync.Mutex

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	userList = []string{}
	_ = app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("总共参与人数：%d\n", count)
}

func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	mu.Lock()
	defer mu.Unlock()
	count := len(users)
	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("总数：%d，实际导入人数：%d", count, count2)
}

func (c *lotteryController) GetLucky() string {
	mu.Lock()
	defer mu.Unlock()
	count := len(userList)
	if count > 1 {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("当前中奖用户%s", user)
	} else if count == 1 {
		user := userList[0]
		userList = nil
		return fmt.Sprintf("当前中奖用户%s", user)
	} else {
		return fmt.Sprintf("当前没有参与用户")
	}
}