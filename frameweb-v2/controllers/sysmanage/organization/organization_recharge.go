package organization

import (
	"phagego/frameweb-v2/controllers/sysmanage"
	"github.com/astaxie/beego"
	. "phagego/frameweb-v2/models"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/astaxie/beego/orm"
	"fmt"
	"time"
)

type OrganizationRechargeIndexController struct {
	sysmanage.BaseController
}

func (c *OrganizationRechargeIndexController) Get() {
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(OrganizationRecharge).Paginate(page, limit, c.LoginAdminRefId)
	pagination.SetPaginator(c.Ctx, limit, total)
	var orgMap = make(map[int64]string)
	if len(list) > 0 {
		orgIds := make([]int64, 0)
		for _, v := range list {
			orgIds = append(orgIds, v.RefId)
		}
		o := orm.NewOrm()
		var orgs []Organization
		o.QueryTable(new(Organization)).Filter("Id__in", orgIds).All(&orgs, "Id", "Name")
		for _, v := range orgs {
			orgMap[v.Id] = v.Name
		}
	}

	c.Data["orgMap"] = orgMap
	c.Data["dataList"] = &list
	c.TplName = "sysmanage/organization_recharge/index.html"
}

func (c *OrganizationRechargeIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	if c.LoginAdminRefId != 0 {
		msg = "没有权限"
		return
	}
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	_, err1 := o.Delete(&OrganizationRecharge{Id: id}, "Id")
	if err1 != nil {
		beego.Error("Delete OrganizationRecharge eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (c *OrganizationRechargeIndexController) SetSuccess() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	if c.LoginAdminRefId != 0 {
		msg = "非法操作"
		return
	}
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "参数错误"
		return
	}
	o := orm.NewOrm()
	orgRecharge := OrganizationRecharge{Id: id}
	if err := o.Read(&orgRecharge); err != nil {
		msg = "异常，请刷新后重试"
		return
	}
	if orgRecharge.Status != 0 {
		msg = "状态已变更，请刷新并确认"
		return
	}
	o.Begin()
	if _, err := o.QueryTable(new(OrganizationRecharge)).Filter("Id", id).Filter("Status", 0).Update(orm.Params{
		"Version":  orm.ColValue(orm.ColAdd, 1),
		"Status":   1,
		"Modifior": c.LoginAdminId,
	}); err != nil {
		o.Rollback()
		msg = "更新状态失败"
		return
	}
	if _, err := o.QueryTable(new(Organization)).Filter("Id", orgRecharge.RefId).Update(orm.Params{
		"Version":        orm.ColValue(orm.ColAdd, 1),
		"RechargeAmount": orm.ColValue(orm.ColAdd, orgRecharge.Amount),
		"Modifior":       c.LoginAdminId,
	}); err != nil {
		o.Rollback()
		msg = "加款失败"
		return
	}
	o.Commit()
	code = 1
	msg = "成功"
}

func (c *OrganizationRechargeIndexController) SetFail() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	if c.LoginAdminRefId != 0 {
		msg = "非法操作"
		return
	}
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "参数错误"
		return
	}
	remark := c.GetString("remark")
	o := orm.NewOrm()
	if _, err := o.QueryTable(new(OrganizationRecharge)).Filter("Id", id).Filter("Status", 0).Update(orm.Params{
		"Version":  orm.ColValue(orm.ColAdd, 1),
		"Status":   2,
		"Remark":   remark,
		"Modifior": c.LoginAdminId,
	}); err != nil {
		msg = "更新状态失败"
		return
	}
	code = 1
	msg = "更新成功"
}

type OrganizationRechargeAddController struct {
	sysmanage.BaseController
}

func (c *OrganizationRechargeAddController) Get() {
	c.TplName = "sysmanage/organization_recharge/add.html"
}

func (c *OrganizationRechargeAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	model := OrganizationRecharge{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	amountF, err := c.GetFloat("amount")
	if err != nil {
		msg = "金额必须为数字，且最多两位小数"
		return
	}
	model.Amount = int(amountF * 1000 / 10)
	model.RefId = c.LoginAdminRefId
	model.OrderNo = fmt.Sprintf("%d%s", c.LoginAdminId, time.Now().Format("20060102150405"))
	model.Status = 0
	model.Creator = c.LoginAdminId
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&model); err != nil {
		msg = "提交失败"
		beego.Error("提交失败", err)
	} else {
		code = 1
		msg = "提交成功"
	}
}

type OrganizationRechargeEditController struct {
	sysmanage.BaseController
}

func (c *OrganizationRechargeEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := OrganizationRecharge{}
	model.Id = id
	o.Read(&model)

	c.Data["data"] = model
	c.TplName = "sysmanage/organization_recharge/edit.html"
}

func (c *OrganizationRechargeEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("OrganizationRechargeIndexController.get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := OrganizationRecharge{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	amountF, err := c.GetFloat("amount")
	if err != nil {
		msg = "金额必须为数字，且最多两位小数"
		return
	}
	model.Amount = int(amountF * 1000 / 10)
	o := orm.NewOrm()
	if _, err := o.QueryTable(new(OrganizationRecharge)).Filter("Id", model.Id).Filter("Status", 0).Update(orm.Params{
		"Amount": model.Amount,
		"Remark": model.Remark,
		"Modifior": c.LoginAdminId,
	}); err != nil {
		msg = "更新失败"
		beego.Error("更新失败", err)
		return
	}
	code = 1
	msg = "更新成功"
}
