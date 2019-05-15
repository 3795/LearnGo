package route

import (
	"github.com/kataras/iris/_benchmarks/iris-mvc-templates/controllers"
	"github.com/kataras/iris/mvc"
	"project/lottery/project/bootstrap"
	"project/lottery/project/service"
)

func Configure(b *bootstrap.Bootstrapper) {
	userService := service.NewUserService()
	giftService := service.NewPrizeService()
	codeService := service.NewCodeService()
	resultService := service.NewResultService()
	userdayService := service.NewUserdayService()
	blackipService := service.NewBlackipService()

	index := mvc.New(b.Party("/"))
	index.Register(userService, giftService, codeService, resultService, userdayService, blackipService)
	index.Handle(new(controllers.IndexController))

	//admin := mvc.New(b.Party("/admin"))
	//admin.Router.Use(middleware.BasicAuth)
	//admin.Register(userService, giftService, codeService, resultService, userdayService, blackipService)
	//admin.Handle(new(controllers.AdminController))
}
