package missiondetail

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
)

type IndexMissionDetailController struct {
	sysmanage.BaseController
}

func (this *IndexMissionDetailController) Get() {
	missionid, _ := this.GetInt64("missionid")
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	o := orm.NewOrm()
	var mission Mission
	err1 := o.QueryTable(new(Mission)).Filter("Id", missionid).One(&mission, "CountEnble")
	if err1 != nil {
		beego.Error("get Mission CountEnble fauil", err)
	}
	describe := utils.GetMissionDescribe(missionid)
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MissionDetail).Paginate(page, limit, missionid)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["countenable"] = mission.CountEnble
	this.Data["describe"] = describe
	this.Data["missionid"] = missionid
	this.Data["dataList"] = list
	this.TplName = "common/missiondetail/index.html"
}

func (this *IndexMissionDetailController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	md := MissionDetail{Id: id}
	err := o.Read(&md)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect("IndexMissionDetailController.get", 302)
	}
	_, err1 := o.Delete(&md, "Id")
	if err1 != nil {
		beego.Error("删除任务详情失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
		return
	}
}

type AddMissionDetailController struct {
	sysmanage.BaseController
}

func (this *AddMissionDetailController) Get() {
	missionid, _ := this.GetInt64("missionid")
	//判断是否可计算
	o := orm.NewOrm()
	var mission Mission
	err := o.QueryTable(new(Mission)).Filter("Id", missionid).One(&mission, "CountEnble")
	if err != nil {
		beego.Error("get Mission CountEnble fauil", err)
	}
	this.Data["countenble"] = mission.CountEnble
	this.Data["missionid"] = missionid
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "common/missiondetail/add.html"
}

func (this *AddMissionDetailController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	md := MissionDetail{}
	if err := this.ParseForm(&md); err != nil {
		beego.Error("任务详情异常", err)
		msg = "参数异常"
		return
	}
	o := orm.NewOrm()
	_, err1 := o.Insert(&md)
	if err1 != nil {
		beego.Error("添加任务详情失败", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type EditMissionDetailController struct {
	sysmanage.BaseController
}

func (this *EditMissionDetailController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	md := MissionDetail{Id: id}
	//查询是否是可计算
	var mission Mission
	err1 := o.QueryTable(new(Mission)).Filter("Id", id).One(&mission, "CountEnble")
	if err1 != nil {
		beego.Error("get Mssion CountEnble failure", err1)
	}
	err := o.Read(&md)
	if err != nil {
		this.Redirect("IndexMissionDetailController.get", 302)
	} else {
		this.Data["countenble"] = mission.CountEnble
		this.Data["data"] = md
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/missiondetail/edit.html"
	}
}

func (this *EditMissionDetailController) Post() {
	var code int
	var msg string
	mi, _ := this.GetInt64("MissionId")
	url := beego.URLFor("IndexMissionDetailController.get", "missionid", mi)
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	md := MissionDetail{}
	if err := this.ParseForm(&md); err != nil {
		beego.Error("参数异常", err)
		msg = "参数异常"
		return
	}
	cols := []string{"Content", "Award", "MinLevel", "MaxLevel"}
	md.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	_, err := o.Update(&md, cols...)
	if err != nil {
		beego.Error("更新任务详情失败", err)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
	}
}
