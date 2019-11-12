package level

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
)

type LevelController struct {
	sysmanage.BaseController
}

func (this *LevelController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit := 1000
	list, total := new(Level).Paginate(page, limit)
	pagination.SetPaginator(this.Ctx, limit, total)

	this.Data["dataList"] = list
	this.TplName = "common/level/index.html"
}

func (this *LevelController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	level := Level{Id: id}
	o := orm.NewOrm()
	_, err := o.Delete(&level, "Id")
	if err != nil {
		beego.Error("删除VIP等级失败", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
	}
}

type LevelAddController struct {
	sysmanage.BaseController
}

func (this *LevelAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "common/level/add.html"
}

func (this *LevelAddController) Post() {
	var code int
	var msg string
	url := beego.URLFor("LevelController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)

	level := Level{}
	if err := this.ParseForm(&level); err != nil {
		beego.Error("vip等级异常", err)
		msg = "数据异常"
		return
	}
	level.Creator = this.LoginAdminId
	level.Modifior = this.LoginAdminId
	_, err1 := level.Create()
	if err1 != nil {
		beego.Error("添加VIP等级失败", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type LevelEditController struct {
	sysmanage.BaseController
}

func (this *LevelEditController) Get() {
	id, _ := this.GetInt64("id")
	level := Level{Id: id}
	o := orm.NewOrm()
	err := o.Read(&level)

	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect(beego.URLFor("LevelController.get"), 302)
	}
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.Data["data"] = level
	this.TplName = "common/level/edit.html"
}

func (this *LevelEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("LevelController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)

	level := Level{}
	err := this.ParseForm(&level)
	if err != nil {
		beego.Error("修改参数异常", err)
		msg = "参数异常"
		return
	}

	cols := []string{"VipLevel", "TotalBet", "LevelGift", "MonthBet", "MonthGift", "VipName", "Bgimg", "Colorimg", "Wbimg", "KeepLevelAmount", "KeepLevelDown"}
	level.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err1 := o.Update(&level, cols...); err1 != nil {
		beego.Error("更新VIP等级失败", err1)
		msg = "更新失败"
	} else {
		code = 1
		msg = "更新成功"
	}
}
