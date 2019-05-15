package controller

import (
	"github.com/kataras/iris"
	"project/lottery/project/models"
	"project/lottery/project/service"
)

type IndexController struct {
	Ctx iris.Context
	ServiceUser    service.UserService
	ServiceGift    service.PrizeService
	ServiceCode    service.CouponService
	ServiceResult  service.LotteryRecordService
	ServiceUserday service.UserdayService
	ServiceBlackip service.BlackipService
}

// http://localhost:8080/
func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "welcome to GO-Lottery, <a href='/public/index.html'>开始抽奖</a>"
}

//http://localhost:8080/gifts
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	datalist := c.ServiceGift.GetAll()
	list := make([]models.TbPrize, 0)
	for _, data := range datalist {
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}
	rs["gift"] = list
	return rs
}

// http://localhost:8080/newprize
func (c *IndexController) GetNewprize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	gifts := c.ServiceGift.GetAll()
	var giftIds []int
	for _, data := range gifts {
		// 虚拟券或者实物奖才需要放到外部榜单中展示
		if data.Gtype > 1 {
			giftIds = append(giftIds, data.Id)
		}
	}
	list := c.ServiceResult.GetNewPrize(50, giftIds)
	rs["prize_list"] = list
	return rs
}

// http://localhost:8080/myprize
//func (c *IndexController) GetMyprize() map[string]interface{} {
//	rs := make(map[string]interface{})
//	rs["code"] = 0
//	rs["msg"] = ""
//
//	//验证是否已登录
//	loginuser := common.GetLoginUser(c.Ctx.Request())
//	if loginuser == nil || loginuser.Uid < 1 {
//		rs["code"] = 101
//		rs["msg"] = "请先登录"
//		return rs
//	}
//	// 只读取最新的100次中奖纪录
//	list := c.ServiceResult.SearchByUser(loginuser.Uid, 1, 100)
//	rs["prize_list"] = list
//	//今天抽奖次数
//	day, _ := strconv.Atoi(common.FormatFromUnixTimeShort(time.Now().Unix()))
//	num := c.ServiceUserday.Count(loginuser.Uid, day)
//	rs["prize_num"] = conf.UserPrizeMax - num
//	return rs
//}

////登录 GET /login
//func (c *IndexController) GetLogin() {
//	uid := common.Random(100000)
//	loginuser := models.ObjLoginuser{
//		Uid:      uid,
//		Username: fmt.Sprintf("admin-%d", uid),
//		Now:      common.NowUnix(),
//		Ip:       common.ClientIP(c.Ctx.Request()),
//	}
//	refer := c.Ctx.GetHeader("Referer")
//	if refer == "" {
//		refer = "/public/index.html?from=login"
//	}
//	common.SetLoginuser(c.Ctx.ResponseWriter(), &loginuser)
//	common.Redirect(c.Ctx.ResponseWriter(), refer)
//}
//
////登出 GET /logout
//func (c *IndexController) GetLogout() {
//	refer := c.Ctx.GetHeader("Referer")
//	if refer == "" {
//		refer = "/public/index.html?from=logout"
//	}
//	common.SetLoginuser(c.Ctx.ResponseWriter(), nil)
//	common.Redirect(c.Ctx.ResponseWriter(), refer)
//}
