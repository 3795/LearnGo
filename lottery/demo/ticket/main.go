package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"time"
)

type lotterController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotterController{})
	return app
}

func main() {
	app := newApp()
	_ = app.Run(iris.Addr(":8080"))
}

// 刮刮乐模式
func (c *lotterController) Get() string {
	var prize string
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(10)
	switch {
	case code == 1:
		prize = "一等奖"
	case code >= 2 && code <= 3:
		prize = "二等奖"
	case code >= 4 && code <= 5:
		prize = "三等奖"
	default:
		return fmt.Sprintf("非酋，你的彩票 %d 没有中奖", code)
	}

	return fmt.Sprintf("RMB玩家，你的彩票 %d 中奖啦，等级为 %s", code, prize)
}

func (c *lotterController) GetPrize() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	var prize [7]int
	// 6个红色球
	for i:=0; i<6; i++ {
		prize[i] = r.Intn(33) + 1
	}
	// 最后一位蓝色球，1-16
	prize[6] = r.Intn(16)
	return fmt.Sprintf("中奖号码是 %v", prize)
}
