package timegift

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"phage/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
)

type IndexTimeGiftController struct {
	sysmanage.BaseController
}

func (this *IndexTimeGiftController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(TimeGift).Paginate(page, limit)
	o := orm.NewOrm()
	var vb VipAttribute
	var vbs VipAttribute
	var notice VipAttribute
	o.QueryTable(new(VipAttribute)).Filter("code", utils.TimeGifttime).One(&vb)
	o.QueryTable(new(VipAttribute)).Filter("code", utils.TimeGiftStatus).One(&vbs)
	o.QueryTable(new(VipAttribute)).Filter("code", utils.TimeGiftNotice).One(&notice)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["vipattribute"] = vb
	this.Data["vbs"] = vbs
	this.Data["notice"] = notice
	this.Data["dataList"] = list
	this.TplName = "common/timegift/index.html"
}

func (this *IndexTimeGiftController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	timeGift := TimeGift{Id: id}
	o := orm.NewOrm()
	_, err := o.Delete(&timeGift)
	if err != nil {
		beego.Error("del timeGift error", err)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *IndexTimeGiftController) ModifyAttr() {
	var code int
	var msg string
	time := this.GetString("time")
	notice := this.GetString("notice")
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	if time != "" {
		timeattribute := VipAttribute{Code: utils.TimeGifttime, Value: time}
		o := orm.NewOrm()
		bool, _, err := o.ReadOrCreate(&timeattribute, "Code")
		if err != nil {
			beego.Error("readorcreate vipattribute error", err)
			msg = "修改失败(1)"
			return
		}
		if !bool {
			_, err1 := o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGifttime).Update(orm.Params{"Value": time})
			if err1 != nil {
				beego.Error("update vipattributetime value error", err)
				msg = "修改失败(2)"
				return
			}
			code = 1
			msg = "修改成功"
		} else {
			code = 1
			msg = "修改成功"
		}
	}
	if notice != "" {
		timeattribute := VipAttribute{Code: utils.TimeGiftNotice, Value: notice}
		o := orm.NewOrm()
		bool, _, err := o.ReadOrCreate(&timeattribute, "Code")
		if err != nil {
			beego.Error("readorcreate vipattribute error", err)
			msg = "修改失败(1)"
			return
		}
		if !bool {
			_, err1 := o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGiftNotice).Update(orm.Params{"Value": notice})
			if err1 != nil {
				beego.Error("update vipattributenotice value error", err)
				msg = "修改失败(2)"
				return
			}
			code = 1
			msg = "修改成功"
		} else {
			code = 1
			msg = "修改成功"
		}
	}

}

func (this *IndexTimeGiftController) ModifyStatus() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	timeattribute := VipAttribute{Code: utils.TimeGiftStatus, Value: "1"}
	o := orm.NewOrm()
	bool, _, err := o.ReadOrCreate(&timeattribute, "Code")
	if err != nil {
		beego.Error("readorcreate vipattribute error", err)
		msg = "开启失败(1)"
		return
	}
	if !bool {
		var vb VipAttribute
		_ = o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGiftStatus).One(&vb)
		var s string
		if vb.Value == "1" {
			s = "0"
		} else {
			s = "1"
		}
		_, err1 := o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGiftStatus).Update(orm.Params{"Value": s})
		if err1 != nil {
			beego.Error("update vipattribute value error", err)
			msg = "修改失败(2)"
			return
		}
		code = 1
		msg = "修改成功"
	} else {
		code = 1
		msg = "修改成功"
	}
}

type AddTimeGiftController struct {
	sysmanage.BaseController
}

func (this *AddTimeGiftController) Get() {
	this.TplName = "common/timegift/add.html"
}
func (this *AddTimeGiftController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	timeGift := TimeGift{}
	if err := this.ParseForm(&timeGift); err != nil {
		msg = "参数异常"
		return
	}
	timeGift.Creator = this.LoginAdminId
	timeGift.Modifior = this.LoginAdminId
	if timeGift.MinMoney > timeGift.MaxMoney {
		msg = "最小金额不能大于最大金额"
		return
	} else if timeGift.MinMoney == timeGift.MaxMoney {
		msg = "最小金额不能等于最大金额"
		return
	}
	_, err1 := timeGift.Create()
	if err1 != nil {
		beego.Error("insert timgtift  error", err1)
		msg = "添加失败"
		return
	} else {
		code = 1
		msg = "添加成功"
	}
}

type EditTimeGiftController struct {
	sysmanage.BaseController
}

func (this *EditTimeGiftController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	timeGift := TimeGift{Id: id}
	err := o.Read(&timeGift)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("IndexTimeGiftController.get"), 302)
	} else {
		this.Data["data"] = timeGift
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/timegift/edit.html"
	}
}
func (this *EditTimeGiftController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	timeGift := TimeGift{}
	if err := this.ParseForm(&timeGift); err != nil {
		beego.Error("parameter error", err)
		msg = "参数异常"
		return
	}
	cols := []string{"GiftName", "GiftContent", "GiftLevel", "MinMoney", "MaxMoney", "Category"}
	if timeGift.MinMoney > timeGift.MaxMoney {
		msg = "最小金额不能大于最大金额"
		return
	} else if timeGift.MinMoney == timeGift.MaxMoney {
		msg = "最小金额不能等于最大金额"
		return
	}
	_, err1 := timeGift.Update(cols...)
	if err1 != nil {
		beego.Error("update timegift error", err1)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
		return
	}
}
