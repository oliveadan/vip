package index

import (
	"fmt"
	"net"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/common/utils"
	"phagego/plugins/googleauth"
	"time"
	"github.com/astaxie/beego/orm"
	. "phagego/frameweb-v2/models"
	"github.com/astaxie/beego"
)

type SysIndexController struct {
	sysmanage.BaseController
}

func (c *SysIndexController) Get() {
	o := orm.NewOrm()
	var admin = Admin{Id: c.LoginAdminId}
	o.Read(&admin)
	c.Data["loginVerify"] = admin.LoginVerify
	c.TplName = "sysmanage/index/index.html"
}

func (c *SysIndexController) GetAuth() {
	user := "LOGIN-" + c.LoginAdminUsername
	ok, secret, qrCode := googleauth.GetGAuthQr(user)
	o := orm.NewOrm()
	if num, err := o.QueryTable(new(Admin)).Filter("Id", c.LoginAdminId).Update(orm.Params{
		"GaSecret":   secret,
		"Modifior":   c.LoginAdminId,
		"ModifyDate": time.Now(),
	}); err != nil || num != 1 {
		beego.Error("SysIndexController GetAuth", err, num)
		ok = false
	}
	c.Data["ok"] = ok
	c.Data["qrCode"] = qrCode
	c.Data["postUrl"] = c.URLFor("SysIndexController.PostAuth")
	c.TplName = "sysmanage/auth/auth.html"
}

func (c *SysIndexController) PostAuth() {
	var code int
	var msg string
	var reUrl string
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reUrl)
	authCode := c.GetString("auth_code")
	o := orm.NewOrm()
	var admin Admin
	if err := o.QueryTable(new(Admin)).Filter("Id", c.LoginAdminId).One(&admin, "GaSecret"); err != nil {
		msg = "绑定失败，请重试"
		return
	}

	if ok, err := googleauth.VerifyGAuth(admin.GaSecret, authCode); err != nil || !ok {
		beego.Error("SysIndexController PostAuth", err, ok)
		msg = "安全码验证失败，请确认"
		return
	}
	if num, err := o.QueryTable(new(Admin)).Filter("Id", c.LoginAdminId).Update(orm.Params{
		"LoginVerify": 2,
		"Modifior":    c.LoginAdminId,
		"ModifyDate":  time.Now(),
	}); err != nil || num != 1 {
		beego.Error("SysIndexController PostAuth", err, num)
		msg = "绑定失败，请重试"
		return
	}
	code = 1
	msg = "绑定成功"
	reUrl = c.URLFor("SysIndexController.Get")
}

func (c *SysIndexController) Systeminfo() {
	var code int
	var msg string
	var data = make([]string, 0)
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &data)
	token := c.GetString("token")
	if token == "" {
		return
	}
	t := time.Now().Format("2006-01-02")
	if token != utils.Md5(t, utils.Pubsalt) {
		return
	}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		msg = fmt.Sprintf("fail to get net interfaces: %v", err)
		return
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		data = append(data, fmt.Sprintf("%d,%s,%d,%s,%s", netInterface.MTU, netInterface.Flags.String(), netInterface.Index, netInterface.HardwareAddr.String(), netInterface.Name))
	}
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		msg = fmt.Sprintf("fail to get net InterfaceAddrs addrs: %v", err)
		return
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				data = append(data, ipNet.IP.String())
			}
		}
	}
	code = 1
}
