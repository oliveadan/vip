package login

import (
	"fmt"
	"html/template"
	. "phagego/frameweb-v2/models"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/common/utils"
	"time"
	. "phagego/frameweb-v2/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"phagego/plugins/googleauth"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	beego.Warn("LoginController Get from ip:", c.Ctx.Input.IP())
	var pubkey, prikey string
	if c.GetSession("loginpubkey") != nil && c.GetSession("loginprikey") != nil {
		pubkey = c.GetSession("loginpubkey").(string)
	} else {
		pubkey, prikey = RsaGenerateKey(1024)
		c.SetSession("loginprikey", prikey)
		c.SetSession("loginpubkey", pubkey)
	}
	if beego.BConfig.RunMode == "dev" {
		c.Data["username"] = "admin"
		c.Data["pass"] = "111111"
		c.Data["captchaValue"] = "1"
	} else {
		c.Data["username"] = ""
		c.Data["pass"] = ""
		c.Data["captchaValue"] = ""
	}
	c.Data["year"] = time.Now().Year()
	c.Data["pubkey"] = pubkey
	c.Data["siteName"] = GetSiteConfigValue(Scname)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "sysmanage/login/index.html"
}

func (c *LoginController) Post() {
	ret := make(map[string]interface{})
	username := c.GetString("username")
	pwd := c.GetString("password")
	if username == "" {
		ret["msg"] = "用户名不能为空"
	} else if pwd == "" {
		ret["msg"] = "密码不能为空"
	} else if beego.BConfig.RunMode == "prod" && !GetCpt().VerifyReq(c.Ctx.Request) {
		ret["msg"] = "验证码错误"
	} else {
		if c.GetSession("loginprikey") == nil {
			ret["msg"] = "请刷新后再试"
		} else {
			prikey := c.GetSession("loginprikey").(string)
			pwdDecrypt := RsaDecrypt(pwd, prikey)
			//pwdDecrypt := pwd
			o := orm.NewOrm()
			admin := Admin{Username: username}
			err := o.Read(&admin, "Username")
			if err != nil {
				beego.Error("Login error", err)
				ret["msg"] = "用户名或密码错误"
			} else {
				cols := make([]string, 0)
				if admin.Enabled == 0 {
					ret["msg"] = "用户名或密码错误"
				} else if admin.Locked == 1 {
					ret["msg"] = "账号已被锁定，无法登录"
				} else if admin.Password != Md5(pwdDecrypt, Pubsalt, admin.Salt) {
					ret["msg"] = "用户名或密码错误"
					admin.LoginFailureCount += 1
					if admin.LoginFailureCount >= 5 {
						admin.Locked = 1
						cols = append(cols, "Locked")
					}
					cols = append(cols, "LoginFailureCount")
				} else {
					if admin.LoginVerify == 1 { // 需要邮箱验证
						go SendMailVerifyCode(admin.Email)
						ret["code"] = 2
						ret["msg"] = "登录验证，请输入邮箱验证码"
					} else if admin.LoginVerify == 2 { // 需要谷歌安全码验证
						ret["code"] = 3
						ret["msg"] = "登录验证，请输入谷歌安全码"
					} else {
						token := GetGuid()
						SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, 28800)
						c.SetSession("token", token)
						c.SetSession("loginAdminId", admin.Id)
						c.SetSession("loginAdminRefId", admin.RefId)
						c.SetSession("loginAdminName", admin.Name)
						c.SetSession("loginAdminUsername", admin.Username)
						ret["code"] = 1
						ret["msg"] = "登录成功"
						ret["url"] = c.URLFor("BaseController.Index")
					}
					admin.LoginFailureCount = 0
					admin.LoginIp = c.Ctx.Input.IP()
					admin.LoginDate = time.Now()
					cols = append(cols, "LoginFailureCount", "LoginIp", "LoginDate")
				}
				if len(cols) > 0 {
					o.Update(&admin, cols...)
				}
			}
		}
	}
	beego.Warn("LoginController Post from ip:", c.Ctx.Input.IP(), "username:", username)
	c.Data["json"] = ret
	c.ServeJSON()
}

func (c *LoginController) LoginVerify() {
	var code int
	var msg string
	var reurl string
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reurl)
	username := c.GetString("username")
	verifyCode := c.GetString("code")
	verifyType, _ := c.GetInt("verify", 2)
	if verifyCode == "" {
		msg = "验证码不能为空"
		return
	}

	o := orm.NewOrm()
	admin := Admin{Username: username}
	if err := o.Read(&admin, "Username"); err != nil {
		msg = "验证失败，请重试"
		return
	}
	var isVerify bool
	if verifyType == 2 {
		isVerify = VerifyMailVerifyCode(admin.Email, verifyCode)
	} else if verifyType == 3 {
		isVerify, _ = googleauth.VerifyGAuth(admin.GaSecret, verifyCode)
	}
	if !isVerify {
		msg = "验证失败"
		return
	}
	token := GetGuid()
	SetCache(fmt.Sprintf("loginAdminId%d", admin.Id), token, 28800)
	c.SetSession("token", token)
	c.SetSession("loginAdminId", admin.Id)
	c.SetSession("loginAdminRefId", admin.RefId)
	c.SetSession("loginAdminName", admin.Name)
	c.SetSession("loginAdminUsername", admin.Username)
	code = 1
	msg = "验证成功"
	reurl = c.URLFor("BaseController.Index")
}

func (c *LoginController) Logout() {
	DelCache(fmt.Sprintf("loginAdminId%v", c.GetSession("loginAdminId")))
	c.DelSession("loginAdminId")
	c.Redirect(c.URLFor("LoginController.Get"), 302)
}
