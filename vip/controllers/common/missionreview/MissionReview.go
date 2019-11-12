package missionreview

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"math"
	"os"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
	"strings"
	"time"
)

type IndexMissionReview struct {
	sysmanage.BaseController
}

func (this *IndexMissionReview) Get() {
	//导出
	isExport, _ := this.GetInt("isExport", 0)
	if isExport == 1 {
		this.Export()
		return
	}
	account := strings.TrimSpace(this.GetString("account"))
	timeStart := strings.TrimSpace(this.GetString("timeStart"))
	timeEnd := strings.TrimSpace(this.GetString("timeEnd"))
	status, _ := this.GetInt("status", 0)
	missionId, _ := this.GetInt64("missionId", 0)
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MissionReview).Paginate(page, limit, account, timeStart, timeEnd, status, missionId)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["condArr"] = map[string]interface{}{"account": account,
		"timeStart": timeStart,
		"timeEnd":   timeEnd,
		"status":    status,
		"missionId": missionId}
	//获得任务列表
	this.Data["missionList"] = GetMissions()
	this.Data["dataList"] = list
	this.TplName = "common/missionreview/index.tpl"
}

func (this *IndexMissionReview) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	mr := MissionReview{Id: id}
	o := orm.NewOrm()
	err := o.Read(&mr)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect("IndexMissionReview.get", 302)
	}
	_, err1 := o.Delete(&mr, "Id")
	if err1 != nil {
		beego.Error("删除任务审核失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
		return
	}
}

func (this *IndexMissionReview) Review() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	status, _ := this.GetInt("status")
	mr := MissionReview{Id: id}
	mr.Status = status
	mr.DeliveredTime = time.Now()
	o := orm.NewOrm()
	_, err := o.Update(&mr, "Status")
	if err != nil {
		logs.Error("更新审核状态失败", err)
		msg = "更新失败"
		return
	}
	code = 1
	msg = "更新审核状态成功"
	return
}

func (this *IndexMissionReview) ReviewBatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	mid, e := this.GetInt64("missionId", 0)
	if e != nil {
		msg = "获取活动ID失败"
		return
	}
	if mid == 0 {
		msg = "请选择要一键审核的活动"
		return
	}
	o := orm.NewOrm()
	exist := o.QueryTable(new(MissionReview)).Filter("MissionId", mid).Filter("Status", 0).Exist()
	if !exist {
		msg = "没有可标记的数据"
		return
	}
	i, err := o.QueryTable(new(MissionReview)).Filter("MissionId", mid).Update(orm.Params{"Status": 1, "DeliveredTime": time.Now()})
	if err != nil {
		logs.Error("udpate MissionReview batch error", err)
		msg = "一键审核失败"
		return
	}
	code = 1
	msg = fmt.Sprintf("成功审核%d条数据", i)
}

func (this *IndexMissionReview) Export() {
	account := strings.TrimSpace(this.GetString("account"))
	timeStart := strings.TrimSpace(this.GetString("timeStart"))
	timeEnd := strings.TrimSpace(this.GetString("timeEnd"))
	status, _ := this.GetInt("status", 0)
	missionId, _ := this.GetInt64("missionId", 0)

	page := 1
	limit := 1000
	list, total := new(MissionReview).Paginate(page, limit, account, timeStart, timeEnd, status, missionId)
	totalInt := int(total)
	if totalInt > limit {
		page1 := (float64(totalInt) - float64(limit)) / float64(limit)
		page2 := int(math.Ceil(page1))
		for page = 2; page <= (page2 + 1); page++ {
			list1, _ := new(MissionReview).Paginate(page, limit, account, timeStart, timeEnd, status, missionId)
			for _, v := range list1 {
				list = append(list, v)
			}
		}
	}

	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "ID")
	xlsx.SetCellValue("Sheet1", "B1", "会员账号")
	xlsx.SetCellValue("Sheet1", "C1", "活动描述")
	xlsx.SetCellValue("Sheet1", "D1", "申请时间")
	xlsx.SetCellValue("Sheet1", "E1", "审核状态")
	xlsx.SetCellValue("Sheet1", "F1", "审核时间")
	for i, value := range list {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.Account)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), utils.GetMissionDescribe(value.MissionId))
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.CreateDate.Format("2006-01-02 15:04:05"))
		if value.Status == 0 {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), "审核中")
		} else if value.Status == 1 {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), "审核通过")
		} else if value.Status == 2 {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), "审核拒绝")

		}
		if !value.DeliveredTime.IsZero() {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), value.DeliveredTime.Format("2006-01-02 15:04:05"))
		}
	}
	// Save xlsx file by the given path.
	fileName := fmt.Sprintf("./tmp/excel/reviewist_%s.xlsx", time.Now().Format("20060102150405"))
	err := xlsx.SaveAs(fileName)
	if err != nil {
		beego.Error("Export reward error", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}
