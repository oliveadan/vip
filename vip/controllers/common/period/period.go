package period

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
)

type PeriodIndexController struct {
	sysmanage.BaseController
}

func (this *PeriodIndexController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Period).Paginate(page, limit)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["dataList"] = &list
	this.TplName = "common/period/index.html"
}

func (this *PeriodIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	period := Period{Id: id}
	err := o.Read(&period)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
	}

	_, err1 := o.Delete(&period, "Id")
	if err1 != nil {
		beego.Error("删除周期分类失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
	}
}

type PeriodAddController struct {
	sysmanage.BaseController
}

func (this *PeriodAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "common/period/add.html"
}

func (this *PeriodAddController) Post() {
	var code int
	var msg string
	url := beego.URLFor("PeriodIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	period := Period{}
	if err := this.ParseForm(&period); err != nil {
		beego.Error("添加周期分类参数异常", err)
		msg = "参数异常"
		return
	}
	_, err1 := period.Create()
	if err1 != nil {
		beego.Error("添加周期分类失败", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type PeriodEditController struct {
	sysmanage.BaseController
}

func (this *PeriodEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	period := Period{Id: id}
	err := o.Read(&period)
	if err != nil {
		this.Redirect("PeriodIndexController.get", 302)
	} else {
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.Data["data"] = period
		this.TplName = "common/period/edit.html"
	}
}

func (this *PeriodEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("PeriodIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	period := Period{}
	if err := this.ParseForm(&period); err != nil {
		beego.Error("修改周期分类参数异常")
		msg = "参数异常"
		return
	}

	cols := []string{"PeriodName", "Rank"}
	o := orm.NewOrm()
	period.Modifior = this.LoginAdminId
	_, err1 := o.Update(&period, cols...)
	if err1 != nil {
		beego.Error("更新周期分类参数异常", err1)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
	}

	_, err2 := o.QueryTable(new(MemberSingle)).Filter("PeriodName", period.PeriodName).Update(orm.Params{
		"PeriodSeq": period.Rank,
	})
	if err2 != nil {
		beego.Error("更新周期分类参数异常(1)", err2)
		msg = "更新失败(1)"
		return
	} else {
		code = 1
		msg = "更新成功"
	}
}
