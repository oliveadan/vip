package mission

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"

	. "phagego/phage-vip4-web/models/common"
)

type IndexMissionController struct {
	sysmanage.BaseController
}

func (this *IndexMissionController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Mission).Paginate(page, limit)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["dataList"] = list
	this.TplName = "common/mission/index.html"
}

func (this *IndexMissionController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	mission := Mission{Id: id}
	o := orm.NewOrm()
	err := o.Read(&mission)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect("IndexMissionController.get", 302)
	}
	_, err1 := o.Delete(&mission, "Id")
	if err1 != nil {
		beego.Error("删除任务失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
		return
	}
}

type AddMissionController struct {
	sysmanage.BaseController
}

func (this *AddMissionController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "common/mission/add.html"
}

func (this *AddMissionController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	mission := Mission{}
	if err := this.ParseForm(&mission); err != nil {
		beego.Error("任务参数异常", err)
		msg = "参数异常"
		return
	}
	o := orm.NewOrm()
	_, err1 := o.Insert(&mission)
	if err1 != nil {
		beego.Error("添加任务失败", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type EditMissionController struct {
	sysmanage.BaseController
}

func (this *EditMissionController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	mission := Mission{Id: id}
	err := o.Read(&mission)
	if err != nil {
		this.Redirect("IndexMissionController.get", 302)
	} else {
		this.Data["data"] = mission
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/mission/edit.html"
	}
}

func (this *EditMissionController) Post() {
	var code int
	var msg string
	url := beego.URLFor("IndexMissionController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	mission := Mission{}
	if err := this.ParseForm(&mission); err != nil {
		beego.Error("参数异常", err)
		msg = "参数异常"
		return
	}
	cols := []string{"Describe", "CountEnble", "Integral", "StartTime", "EndTime", "SumEnable"}
	mission.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	_, err := o.Update(&mission, cols...)
	if err != nil {
		beego.Error("更新任务失败", err)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
	}
}
