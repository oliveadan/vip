package organization

import (
	"phagego/frameweb-v2/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "phagego/frameweb-v2/models"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
	"strings"
)

func validate(org *Organization) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(org.Name, "errmsg").Message("组织名称必填")
	valid.MaxSize(org.BindDomain, 127, "errmsg").Message("绑定域名最长127个字符")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type OrganizationIndexController struct {
	sysmanage.BaseController
}

func (c *OrganizationIndexController) Get() {
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Organization).Paginate(page, limit)
	pagination.SetPaginator(c.Ctx, limit, total)
	c.Data["dataList"] = &list
	c.TplName = "sysmanage/organization/index.html"
}

func (c *OrganizationIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("OrganizationIndexController")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	if c.LoginAdminRefId != 0 {
		msg = "没有权限"
		return
	}
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	_, err1 := o.Delete(&Organization{Id: id}, "Id")
	if err1 != nil {
		beego.Error("Delete Organization eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type OrganizationAddController struct {
	sysmanage.BaseController
}

func (c *OrganizationAddController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "sysmanage/organization/add.html"
}

func (c *OrganizationAddController) Post() {
	var code int
	var msg string
	var url = beego.URLFor("OrganizationIndexController.get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	organization := Organization{}
	if err := c.ParseForm(&organization); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&organization); hasError {
		msg = errMsg
		return
	}
	var zeroTime time.Time
	if organization.ExpireTime == zeroTime {
		organization.ExpireTime = time.Now().Add(time.Hour * 26280) // 默认三年
	}
	organization.Creator = c.LoginAdminId
	organization.Modifior = c.LoginAdminId
	if organization.BindDomain != "" && !strings.HasSuffix(organization.BindDomain, ",") {
		organization.BindDomain = organization.BindDomain + ","
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&organization); err != nil {
		msg = "添加失败"
		beego.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type OrganizationEditController struct {
	sysmanage.BaseController
}

func (c *OrganizationEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := Organization{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("OrganizationIndexController.get"), 302)
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["data"] = model
	c.TplName = "sysmanage/organization/edit.html"
}

func (c *OrganizationEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("OrganizationIndexController.get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	organization := Organization{}
	if err := c.ParseForm(&organization); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&organization); hasError {
		msg = errMsg
		return
	}
	if organization.BindDomain != "" && !strings.HasSuffix(organization.BindDomain, ",") {
		organization.BindDomain = organization.BindDomain + ","
	}
	cols := []string{"Name", "BindDomain", "ModifyDate", "ExpireTime"}
	organization.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&organization, cols...); err != nil {
		msg = "更新失败"
		beego.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
