package bootstrap

import (
	"github.com/kataras/iris"
	"io/ioutil"
	"project/lottery/project/conf"
	"time"
)

type Configurator func(bootstrapper *Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName		string
	AppOwner	string
	AppSpawDate		time.Time
}

func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper{
	b := &Bootstrapper {
		Application: iris.New(),
		AppName: appName,
		AppOwner: appOwner,
		AppSpawDate: time.Now(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

var pathPrefix = "D:/Develop/Go/src/project/lottery/project/web"

func (b *Bootstrapper) SetupViews(viewDir string) {
	htmlEngine := iris.HTML(viewDir, ".html").Layout("/shared/layout.html")
	htmlEngine.Reload(true)

	//添加定制方法
	htmlEngine.AddFunc("FromUnixtimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeFormShort)
	})

	htmlEngine.AddFunc("FromUnixtime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeForm)
	})

	b.RegisterView(htmlEngine)
}

func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonouput := ctx.URLParamExists("json"); jsonouput {
			_, _ = ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		_ = ctx.View("/shared/error.html")
	})
}

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

const (
	Favicon      = "/favicon.ico"
)

var StaticAssets = pathPrefix + "/public"

func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./views")
	b.SetupErrorHandlers()

	// static files
	b.Favicon(StaticAssets + Favicon)
	b.StaticWeb(StaticAssets[1:], StaticAssets)
	indexhtml, err := ioutil.ReadFile(StaticAssets + "/index.html")
	if err == nil {
		b.StaticContent(StaticAssets[1:]+"/", "text/html", indexhtml)
	}
	// 不要把目录末尾"/"省略掉
	iris.WithoutPathCorrectionRedirection(b.Application)

	return b
}

func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	_ = b.Run(iris.Addr(addr), cfgs...)
}
