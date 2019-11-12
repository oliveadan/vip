package quicknav

import (
	"phagego/frameweb-v2/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "phagego/frameweb-v2/models"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"github.com/astaxie/beego/orm"
)

type QuickNavIndexController struct {
	sysmanage.BaseController
}

func (this *QuickNavIndexController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(QuickNav).Paginate(page, limit)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["dataList"] = list
	this.TplName = "sysmanage/quicknav/index.html"
}

func (this *QuickNavIndexController) Delone() {
	var code int
	var msg string
	url := beego.URLFor("QuickNavIndexController")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	model := QuickNav{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		beego.Error("Delete QuickNav eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type QuickNavAddController struct {
	sysmanage.BaseController
}

func (this *QuickNavAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/quicknav/add.html"
}

func (this *QuickNavAddController) Post() {
	var code int
	var msg string
	var url = beego.URLFor("QuickNavIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	quicknav := QuickNav{}
	if err := this.ParseForm(&quicknav); err != nil {
		msg = "参数异常"
		return
	}
	quicknav.Creator = this.LoginAdminId
	quicknav.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&quicknav); err != nil {
		msg = "添加失败"
		beego.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type QuickNavEditController struct {
	sysmanage.BaseController
}

func (this *QuickNavEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	model := QuickNav{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("QuickNavIndexController.get"), 302)
	}
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.Data["data"] = model
	this.TplName = "sysmanage/quicknav/edit.html"
}

func (this *QuickNavEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("QuickNavIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	quicknav := QuickNav{}
	if err := this.ParseForm(&quicknav); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Name", "WebSite", "Icon", "Seq", "ModifyDate"}
	quicknav.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&quicknav, cols...); err != nil {
		msg = "更新失败"
		beego.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
