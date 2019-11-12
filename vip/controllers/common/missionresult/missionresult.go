package missionresult

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"os"
	"phage/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
	"strings"
	"time"
)

type IndexMissionResultController struct {
	sysmanage.BaseController
}

func (this *IndexMissionResultController) Get() {
	//导出
	isExport, _ := this.GetInt("isExport", 0)
	if isExport == 1 {
		this.Export()
		return
	}
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	missionid, _ := this.GetInt64("id")
	account := strings.TrimSpace(this.GetString("account"))
	starttime := strings.TrimSpace(this.GetString("startTime"))
	endtime := strings.TrimSpace(this.GetString("endTime"))
	status, _ := this.GetInt("status")
	limit, _ := beego.AppConfig.Int("pagelimit")

	var period MissionResult
	o := orm.NewOrm()
	err1 := o.QueryTable(new(MissionResult)).Filter("CreateDate__isnull", false).Filter("MissionId", missionid).OrderBy("-CreateDate").Distinct().Limit(-1).One(&period, "CreateDate")
	if err1 != nil {
		beego.Error("获取任务数据期数失败", err1)
	}
	if !period.CreateDate.IsZero() && starttime == "" {
		starttime = period.CreateDate.Format("2006-01-02 15:04:05")
	}

	list, total := new(MissionResult).Paginate(page, limit, missionid, account, starttime, endtime, status)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["condArr"] = map[string]interface{}{
		"account":   account,
		"status":    status,
		"endtime":   endtime,
		"starttime": starttime}
	this.Data["dataList"] = list
	this.Data["Missionid"] = missionid
	//任务内容
	this.Data["describe"] = utils.GetMissionDescribe(missionid)
	this.TplName = "common/missionresult/index.html"
}

func (this *IndexMissionResultController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	mr := MissionResult{Id: id}
	o := orm.NewOrm()
	err := o.Read(&mr)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect("IndexMissionResultController.get", 302)
		return
	}
	_, err1 := o.Delete(&mr, "Id")
	if err1 != nil {
		beego.Error("删除计算结果失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *IndexMissionResultController) Delbatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	missionid, _ := this.GetInt64("id")
	s := this.GetString("startTime")
	e := this.GetString("endTime")
	if s == "" || e == "" {
		msg = "请选择要删除的时间区间"
		return
	}
	o := orm.NewOrm()
	num, err := o.QueryTable(new(MissionResult)).Filter("MissionId", missionid).Filter("CreateDate__gte", s).Filter("CreateDate__lte", e).Delete()
	if err != nil {
		beego.Error("删除计算结果失败", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = fmt.Sprintf("成功删除%d条数据", num)
	}

}

func (this *IndexMissionResultController) Reviewbatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	missionid, _ := this.GetInt64("id")
	s := this.GetString("startTime")
	e := this.GetString("endTime")
	if s == "" || e == "" {
		msg = "请选择批量审核通过的时间区间"
		return
	}
	o := orm.NewOrm()
	num, err := o.QueryTable(new(MissionResult)).Filter("MissionId", missionid).Filter("CreateDate__gte", s).Filter("CreateDate__lte", e).Update(orm.Params{"Status": 1})
	if err != nil {
		beego.Error("批量标记失败", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = fmt.Sprintf("成功标记%d条数据", num)
	}
}

func (this *IndexMissionResultController) Review() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	status, _ := this.GetInt("status")
	mr := MissionResult{Id: id}
	mr.Status = status
	mr.ModifyDate = time.Now()
	mr.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	_, err := o.Update(&mr, "Status")
	if err != nil {
		beego.Error("更新审核状态失败", err)
		msg = "更新失败"
		return
	}
	code = 1
	msg = "更新审核状态成功"
	return
}

func (this *IndexMissionResultController) Export() {
	missionid, _ := this.GetInt64("id")
	o := orm.NewOrm()
	var mr []MissionResult
	_, err := o.QueryTable(new(MissionResult)).Filter("MissionId", missionid).Limit(-1).All(&mr)
	if err != nil {
		beego.Error("导出失败", err)
		return
	}

	xlxs := excelize.NewFile()
	xlxs.SetCellValue("Sheet1", "A1", "编号")
	xlxs.SetCellValue("Sheet1", "B1", "创建时间")
	xlxs.SetCellValue("Sheet1", "C1", "会员账号")
	xlxs.SetCellValue("Sheet1", "D1", "奖品")
	xlxs.SetCellValue("Sheet1", "E1", "是否领取")
	xlxs.SetCellValue("Sheet1", "F1", "领取时间")
	for i, value := range mr {
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.Id)
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.CreateDate.Format("2006-01-02 15:04:05"))
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.Account)
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.Prize)
		if value.Enable == 1 {
			xlxs.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), "是")
		} else {
			xlxs.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), "否")
		}
		xlxs.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), value.GetTime.Format("2006-01-02 15:04:05"))
	}
	fileName := fmt.Sprintf("./tmp/excel/MissionResul%s.xlsx", time.Now().Format("20060102150405"))
	err1 := xlxs.SaveAs(fileName)
	if err1 != nil {
		beego.Error("Export MissionResul error", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}
