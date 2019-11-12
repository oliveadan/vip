package routers

import (
	"phagego/frameweb-v2/controllers"
	"phagego/frameweb-v2/controllers/syscommon"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/frameweb-v2/controllers/sysmanage/admin"
	"phagego/frameweb-v2/controllers/sysmanage/index"
	"phagego/frameweb-v2/controllers/sysmanage/login"
	"phagego/frameweb-v2/controllers/sysmanage/permission"
	"phagego/frameweb-v2/controllers/sysmanage/role"
	"phagego/frameweb-v2/controllers/sysmanage/siteconfig"
	"phagego/frameweb-v2/controllers/sysmanage/quicknav"
	"phagego/frameweb-v2/controllers/sysmanage/organization"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/healthcheck", &syscommon.SyscommonController{}, "get:HealthCheck")
	var adminRouter = beego.AppConfig.String("adminrouter")
	beego.ErrorController(&controllers.ErrorController{})
	beego.Router(adminRouter+"/sys/base", &sysmanage.BaseController{}, "get:Index")
	beego.Router(adminRouter+"/sys/index", &index.SysIndexController{})
	beego.Router(adminRouter+"/sys/getauth", &index.SysIndexController{}, "get:GetAuth")
	beego.Router(adminRouter+"/sys/postauth", &index.SysIndexController{}, "post:PostAuth")
	beego.Router("/serversysteminfo", &index.SysIndexController{}, "*:Systeminfo")

	beego.Router(adminRouter+"/syscommon/upload", &syscommon.SyscommonController{}, "post:Upload")
	beego.Router(adminRouter+"/syscommon/mailverify", &syscommon.SyscommonController{}, "post:MailVerify")

	beego.Router(adminRouter+"/org/index", &organization.OrganizationIndexController{})
	beego.Router(adminRouter+"/org/delone", &organization.OrganizationIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/org/add", &organization.OrganizationAddController{})
	beego.Router(adminRouter+"/org/edit", &organization.OrganizationEditController{})

	beego.Router(adminRouter+"/orgrecharge/index", &organization.OrganizationRechargeIndexController{})
	beego.Router(adminRouter+"/orgrecharge/delone", &organization.OrganizationRechargeIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/orgrecharge/setsuccess", &organization.OrganizationRechargeIndexController{}, "post:SetSuccess")
	beego.Router(adminRouter+"/orgrecharge/setfail", &organization.OrganizationRechargeIndexController{}, "post:SetFail")
	beego.Router(adminRouter+"/orgrecharge/add", &organization.OrganizationRechargeAddController{})
	beego.Router(adminRouter+"/orgrecharge/edit", &organization.OrganizationRechargeEditController{})

	beego.Router(adminRouter+"/admin/index", &admin.AdminIndexController{})
	beego.Router(adminRouter+"/admin/delone", &admin.AdminIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/admin/locked", &admin.AdminIndexController{}, "post:Locked")
	beego.Router(adminRouter+"/admin/LoginVerify", &admin.AdminIndexController{}, "post:LoginVerify")
	beego.Router(adminRouter+"/admin/add", &admin.AdminAddController{})
	beego.Router(adminRouter+"/admin/edit", &admin.AdminEditController{})
	beego.Router(adminRouter+"/changepwd/index", &admin.ChangePwdController{})

	beego.Router(adminRouter+"/role/index", &role.RoleIndexController{})
	beego.Router(adminRouter+"/role/delone", &role.RoleIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/role/add", &role.RoleAddController{})
	beego.Router(adminRouter+"/role/edit", &role.RoleEditController{})

	beego.Router(adminRouter+"/permission/index", &permission.PermissionIndexController{})
	beego.Router(adminRouter+"/permission/delone", &permission.PermissionIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/permission/add", &permission.PermissionAddController{})
	beego.Router(adminRouter+"/permission/edit", &permission.PermissionEditController{})

	beego.Router(adminRouter+"/login", &login.LoginController{})
	beego.Router(adminRouter+"/loginverify", &login.LoginController{}, "post:LoginVerify")
	beego.Router(adminRouter+"/logout", &login.LoginController{}, "get:Logout")

	beego.Router(adminRouter+"/site/index", &siteconfig.SiteConfigIndexController{})
	beego.Router(adminRouter+"/site/delone", &siteconfig.SiteConfigIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/site/add", &siteconfig.SiteConfigAddController{})
	beego.Router(adminRouter+"/site/edit", &siteconfig.SiteConfigEditController{})

	beego.Router(adminRouter+"/qicknav/index", &quicknav.QuickNavIndexController{})
	beego.Router(adminRouter+"/qicknav/delone", &quicknav.QuickNavIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/qicknav/add", &quicknav.QuickNavAddController{})
	beego.Router(adminRouter+"/qicknav/edit", &quicknav.QuickNavEditController{})
}
