package siteconfig

import (
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/frameweb-v2/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strings"
)

func validate(siteConfig *SiteConfig) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(siteConfig.Code, "errmsg").Message("代码必选")
	valid.Required(siteConfig.Value, "errmsg").Message("值必填")
	valid.MaxSize(siteConfig.Value, 1024, "errmsg").Message("值最长1024位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type SiteConfigIndexController struct {
	sysmanage.BaseController
}

func (c *SiteConfigIndexController) Get() {
	var siteConfigList []SiteConfig
	o := orm.NewOrm()
	qs := o.QueryTable(new(SiteConfig))
	qs.All(&siteConfigList)
	// 返回值
	c.Data["dataList"] = siteConfigList
	c.TplName = "sysmanage/siteconfig/index.html"
}

func (c *SiteConfigIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, _ := c.GetInt64("id")
	siteConfig := SiteConfig{Id: id}
	o := orm.NewOrm()
	err := o.Read(&siteConfig)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	} else if siteConfig.IsSystem == 1 {
		msg = "系统内置，不能删除"
		return
	}
	_, err1 := o.Delete(&SiteConfig{Id: id})
	if err1 != nil {
		beego.Error("Delete siteconfig error", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type SiteConfigAddController struct {
	sysmanage.BaseController
}

func (c *SiteConfigAddController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "sysmanage/siteconfig/add.html"
}

func (c *SiteConfigAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	siteConfig := SiteConfig{}
	if err := c.ParseForm(&siteConfig); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&siteConfig); hasError {
		msg = errMsg
		return
	}
	diyCode := strings.TrimSpace(c.GetString("diyCode"))
	if diyCode != "" {
		siteConfig.Code = diyCode
	}
	siteConfig.Creator = c.LoginAdminId
	siteConfig.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&siteConfig, "Code");err != nil {
		msg = "添加失败,请重试"
		beego.Error("Insert siteconfig error", err)
	} else if !created {
		msg = "添加失败，配置已存在"
	} else {
		code = 1
		msg = "添加成功"
	}
}

type SiteConfigEditController struct {
	sysmanage.BaseController
}

func (c *SiteConfigEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	siteConfig := SiteConfig{Id: id}

	err := o.Read(&siteConfig)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("SiteConfigIndexController.get"), 302)
	} else {
		c.Data["data"] = siteConfig
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
		c.TplName = "sysmanage/siteconfig/edit.html"
	}
}

func (c *SiteConfigEditController) Post() {
	var code int
	var msg string
	var reurl = c.URLFor("SiteConfigIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reurl)
	siteConfig := SiteConfig{}
	if err := c.ParseForm(&siteConfig); err != nil {
		msg = "参数异常"
		return
	}
	siteConfig.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&siteConfig, "Value", "ModifyDate"); err != nil {
		msg = "更新失败"
		beego.Error("Update siteconfig error", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
