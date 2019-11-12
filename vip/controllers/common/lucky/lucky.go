package lucky

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
)

type LuckyController struct {
	sysmanage.BaseController
}

func (this *LuckyController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit := 100
	list, total := new(Lucky).Paginate(page, limit)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["dataList"] = list
	this.TplName = "common/lucky/index.html"
}

func (this *LuckyController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	Lucky := Lucky{Id: id}
	o := orm.NewOrm()
	_, err := o.Delete(&Lucky, "Id")
	if err != nil {
		beego.Error("删除好运金配置失败", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
	}
}

type LuckyAddController struct {
	sysmanage.BaseController
}

func (this *LuckyAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "common/lucky/add.html"
}

func (this *LuckyAddController) Post() {
	var code int
	var msg string
	url := beego.URLFor("LuckyController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)

	Lucky := Lucky{}
	if err := this.ParseForm(&Lucky); err != nil {
		beego.Error("好运金配置异常", err)
		msg = "数据异常"
	}
	Lucky.Creator = this.LoginAdminId
	Lucky.Modifior = this.LoginAdminId
	_, err1 := Lucky.Create()
	if err1 != nil {
		beego.Error("添加好运金配置失败", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type LuckyEditController struct {
	sysmanage.BaseController
}

func (this *LuckyEditController) Get() {
	id, _ := this.GetInt64("id")
	Lucky := Lucky{Id: id}
	o := orm.NewOrm()
	err := o.Read(&Lucky)

	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect(beego.URLFor("LuckyController.get"), 302)
	}
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.Data["data"] = Lucky
	this.TplName = "common/lucky/edit.html"
}

func (this *LuckyEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("LuckyController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)

	Lucky := Lucky{}
	err := this.ParseForm(&Lucky)
	if err != nil {
		beego.Error("修改参数异常", err)
		msg = "参数异常"
		return
	}

	cols := []string{"MinVipLevel", "MaxVipLevel", "MonthBet", "MonthBet", "Luckygift"}
	Lucky.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err1 := o.Update(&Lucky, cols...); err1 != nil {
		beego.Error("更新好运金配置失败", err1)
		msg = "更新失败"
	} else {
		code = 1
		msg = "更新成功"
	}
}
