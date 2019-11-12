package admin

import (
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/frameweb-v2/models"
	. "phagego/common/utils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ChangePwdController struct {
	sysmanage.BaseController
}

func (this *ChangePwdController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/admin/changepwd.html"
}

func (this *ChangePwdController) Post() {
	var code int
	var msg string
	var url string
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	id := this.LoginAdminId
	if id == 0 {
		code = 1
		msg = "请重新登录"
		url = this.URLFor("LoginController.Get")
		return
	}
	oldPwd := this.GetString("oldPassword")
	newPwd := this.GetString("newPassword")
	reNewPwd := this.GetString("reNewPassword")
	if oldPwd == "" || strings.TrimSpace(oldPwd) == "" {
		msg = "旧密码不能为空"
		return
	} else if newPwd == "" || strings.TrimSpace(newPwd) == "" {
		msg = "新密码不能为空"
		return
	} else if strings.TrimSpace(newPwd) != strings.TrimSpace(reNewPwd) {
		msg = "两次输入的新密码不一致"
		return
	}
	o := orm.NewOrm()
	admin := Admin{Id: id}
	if err := o.Read(&admin); err != nil {
		msg = "用户信息错误，请重试"
	} else if Md5(oldPwd, Pubsalt, admin.Salt) != admin.Password {
		msg = "旧密码错误"
	} else {
		salt := GetGuid()
		pa := Md5(newPwd, Pubsalt, salt)
		admin.Password = pa
		admin.Salt = salt

		if _, err2 := o.Update(&admin, "Password", "Salt", "ModifyDate"); err2 != nil {
			msg = "更新失败"
			beego.Error("Change password error", err2)
		} else {
			code = 1
			msg = "更新成功，请重新登录"
			url = this.URLFor("LoginController.Logout")
		}
	}
}
