package initial

import (
	"fmt"
	"net"
	. "phagego/frameweb-v2/models"
	. "phagego/frameweb-v2/utils"
	. "phagego/common/utils"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
)

func InitFilter() {
	var adminRouter = beego.AppConfig.String("adminrouter")
	//beego.InsertFilter(adminRouter+"/login", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/sys/index", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/admin/*", beego.BeforeRouter, filterLicense)
	//beego.InsertFilter(adminRouter+"/site/*", beego.BeforeRouter, filterLicense)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeRouter, filterAuth)
	beego.InsertFilter(adminRouter+"/*", beego.BeforeExec, filterXSRFToken)
}

var filterXSRFToken = func(ctx *context.Context) {
	expire := int64(beego.BConfig.WebConfig.XSRFExpire)
	ctx.Input.SetData("xsrf_token", ctx.XSRFToken(beego.BConfig.WebConfig.XSRFKey, expire))
}

/**
 * 登录验证、鉴权
 */
var filterAuth = func(ctx *context.Context) {
	// 不需要鉴权的url
	switch ctx.Request.RequestURI {
	case beego.URLFor("LoginController.Get"):
		return
	case beego.URLFor("LoginController.Logout"):
		return
	case beego.URLFor("LoginController.LoginVerify"):
		return
	}
	// 登录验证
	lid, ok := ctx.Input.Session("loginAdminId").(int64)
	if !ok {
		ctx.Redirect(302, beego.URLFor("LoginController.Get"))
	}
	// token验证
	sestoken, ok := ctx.Input.Session("token").(string)
	var cactoken string
	GetCache(fmt.Sprintf("loginAdminId%d", lid), &cactoken)
	if !ok || sestoken == "" || sestoken != cactoken {
		ctx.ResponseWriter.Write([]byte("登录过期，请重新登录"))
		ctx.Abort(401, "登录过期，请重新登录")
	}

	// 鉴权
	o := orm.NewOrm()
	var arList orm.ParamsList
	_, err := o.QueryTable(new(AdminRole)).Filter("AdminId", lid).ValuesFlat(&arList, "RoleId")
	if err != nil {
		beego.Error("FilterAuth Query AdminRole error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	_, err = o.QueryTable(new(Role)).Filter("Id__in", arList).Filter("Enabled", 1).ValuesFlat(&arList, "Id")
	if err != nil {
		beego.Error("FilterAuth Query AdminRole error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	var rpList orm.ParamsList
	_, err = o.QueryTable(new(RolePermission)).Filter("RoleId__in", arList).Distinct().ValuesFlat(&rpList, "PermissionId")
	if err != nil {
		beego.Error("FilterAuth Query RolePermission error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	var permList orm.ParamsList
	_, err = o.QueryTable(new(Permission)).Filter("Id__in", rpList).Filter("Enabled", 1).ValuesFlat(&permList, "Url")
	if err != nil {
		beego.Error("FilterAuth Query Permission error", err)
		ctx.Abort(500, "内部错误, 请联系管理员")
		return
	}
	var currentUrl = ctx.Request.URL.EscapedPath()
	var isAuth = false
	for _, perm := range permList {
		if perm != nil && perm.(string) != "" && beego.URLFor(perm.(string)) == currentUrl {
			isAuth = true
		}
	}
	// 没有权限
	if !isAuth {
		ctx.ResponseWriter.Write([]byte("没有权限或页面不存在"))
		ctx.Abort(401, "没有权限或页面不存在")
	}
}

var filterLicense = func(ctx *context.Context) {
	beego.Info("filter license")
	// 不需要登录的url
	switch ctx.Request.RequestURI {
	case beego.URLFor("SysIndexController.Systeminfo"):
		return
	}
	lic := beego.AppConfig.String("serverlicense")
	if lic == "" {
		ctx.ResponseWriter.Write([]byte("当前系统为试用版，请购买正版"))
		ctx.Abort(500, "当前系统为试用版，请购买正版")
		beego.Error("License not found, please config!")
		return
	}
	netInterfaces, err := net.Interfaces()
	if err != nil {
		ctx.ResponseWriter.Write([]byte("内部错误，请重试"))
		ctx.Abort(500, "内部错误，请重试")
		beego.Error("fail to get net interfaces:", err)
		return
	}
	indexlen, err := strconv.ParseInt(SubString(lic, 23, 1), 10, 64)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("内部错误，请重试500-1"))
		ctx.Abort(500, "内部错误，请重试500-1")
		beego.Error("fail to conv indexlen:", err)
		return
	}
	index, err := strconv.ParseInt(SubString(lic, 11, int(indexlen)), 10, 64)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("内部错误，请重试500-2"))
		ctx.Abort(500, "内部错误，请重试500-2")
		beego.Error("fail to conv index:", err)
		return
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 || netInterface.Index != int(index) {
			continue
		}
		sign := Md5(strconv.FormatInt(index, 10), Pubsalt, netInterface.HardwareAddr.String())
		if strings.ToUpper(sign) == SubString(lic, 58, 32) {
			return
		}
		break
	}
	ctx.ResponseWriter.Write([]byte("当前系统为试用版，请购买正版"))
	ctx.Abort(500, "当前系统为试用版，请购买正版")
	return
}
