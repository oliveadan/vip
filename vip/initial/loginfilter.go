package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func InitLoginFilter() {
	beego.InsertFilter("/vipcenterindex", beego.BeforeRouter, FilterLogin)
}

var FilterLogin = func(ctx *context.Context) {
	if ctx.Request.RequestURI == beego.URLFor("VipCenterController.post") {
		return
	}
	session := ctx.Input.Session("loginAccount")
	query := ctx.Input.Query("name")
	if session == nil || session.(string) != query {
		ctx.Redirect(301, beego.URLFor("FrontIndexController.get"))
	}
}
