package rewardlog

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"math"
	"os"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
	"strings"
	"time"
)

type RewardLogIndexController struct {
	sysmanage.BaseController
}

func (this *RewardLogIndexController) Get() {
	// 导出
	isExport, _ := this.GetInt("isExport", 0)
	if isExport == 1 {
		this.Export()
		return
	}
	// 条件 要和export、Deliveredbatch 函数保持一致
	beego.Informational("query rewardLog ")
	account := strings.TrimSpace(this.GetString("account"))
	timeStart := strings.TrimSpace(this.GetString("timeStart"))
	timeEnd := strings.TrimSpace(this.GetString("timeEnd"))
	delivered, _ := this.GetInt8("delivered", 0)

	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(RewardLog).Paginate(page, limit, account, timeStart, timeEnd, delivered)
	pagination.SetPaginator(this.Ctx, limit, total)
	// 返回值
	this.Data["condArr"] = map[string]interface{}{"account": account,
		"timeStart": timeStart,
		"timeEnd":   timeEnd,
		"delivered": delivered}
	this.Data["dataList"] = list
	this.TplName = "rewardlog/index.tpl"
}

func (this *RewardLogIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	rewardLog := RewardLog{Id: id}
	o := orm.NewOrm()
	err := o.Read(&rewardLog)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&rewardLog, "Id")
	if err1 != nil {
		beego.Error("Delete rewardLog error", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *RewardLogIndexController) Delbatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	// 条件 要和Get函数保持一致
	o := orm.NewOrm()
	res, err := o.Raw("DELETE from ph_reward_log WHERE id != 0").Exec()
	if err != nil {
		beego.Error("Delete batch memberLottery error", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		num, _ := res.RowsAffected()
		msg = fmt.Sprintf("成功删除%d条记录", num)
		return
	}
}

func (this *RewardLogIndexController) Delivered() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	delivered, _ := this.GetInt8("delivered")
	rl := RewardLog{Id: id}
	o := orm.NewOrm()
	err := o.Read(&rl)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		msg = "数据不存在，请确认"
		return
	}
	rl.Delivered = delivered
	rl.DeliveredTime = time.Now()
	rl.Modifior = this.LoginAdminId
	_, err1 := rl.Update("Delivered", "DeliveredTime")
	if err1 != nil {
		beego.Error("Enabled RewardLog error", err1)
		msg = "操作失败"
	} else {
		code = 1
		msg = "操作成功"
	}

}

func (this *RewardLogIndexController) Deliveredbatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	// 条件 要和export、Get 函数保持一致
	beego.Informational("query rewardLog ")
	account := strings.TrimSpace(this.GetString("account"))
	timeStart := strings.TrimSpace(this.GetString("timeStart"))
	timeEnd := strings.TrimSpace(this.GetString("timeEnd"))
	delivered, _ := this.GetInt8("delivered", 2)

	o := orm.NewOrm()
	qs := o.QueryTable(new(RewardLog))
	cond := orm.NewCondition()
	if account != "" {
		cond = cond.And("Account__contains", account)
	}
	if timeStart != "" {
		cond = cond.And("CreateDate__gte", timeStart)
	}
	if timeEnd != "" {
		cond = cond.And("CreateDate__lte", timeEnd)
	}
	if delivered != 2 {
		cond = cond.And("Delivered", delivered)
	}
	qs = qs.SetCond(cond)
	if num, err := qs.Update(orm.Params{"Delivered": 1, "DeliveredTime": time.Now(), "Modifior": this.LoginAdminId}); err != nil {
		beego.Error("Delivered batch RewardLog error", err)
		msg = "批量标记失败"
	} else {
		code = 1
		msg = fmt.Sprintf("成功标记%d条记录", num)
	}
}

func (this *RewardLogIndexController) Export() {
	// 条件 要和get、Deliveredbatch函数保持一致
	account := strings.TrimSpace(this.GetString("account"))
	timeStart := strings.TrimSpace(this.GetString("timeStart"))
	timeEnd := strings.TrimSpace(this.GetString("timeEnd"))
	delivered, _ := this.GetInt8("delivered", 2)

	page := 1
	limit := 1000
	list, total := new(RewardLog).Paginate(page, limit, account, timeStart, timeEnd, delivered)
	totalInt := int(total)
	if totalInt > limit {
		page1 := (float64(totalInt) - float64(limit)) / float64(limit)
		page2 := int(math.Ceil(page1))
		for page = 2; page <= (page2 + 1); page++ {
			list1, _ := new(RewardLog).Paginate(page, limit, account, timeStart, timeEnd, delivered)
			for _, v := range list1 {
				list = append(list, v)
			}
		}
	}

	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "ID")
	xlsx.SetCellValue("Sheet1", "B1", "会员账号")
	xlsx.SetCellValue("Sheet1", "C1", "奖品名称")
	xlsx.SetCellValue("Sheet1", "D1", "奖品内容")
	xlsx.SetCellValue("Sheet1", "E1", "是否派奖")
	xlsx.SetCellValue("Sheet1", "F1", "派奖时间")
	for i, value := range list {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.Id)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.Account)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.GiftName)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.GiftContent)
		if value.Delivered == 1 {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), "已派送")
		} else {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), "未派送")
		}
		if !value.DeliveredTime.IsZero() {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), value.DeliveredTime.Format("2006-01-02 15:04:05"))
		}
	}
	// Save xlsx file by the given path.
	fileName := fmt.Sprintf("./tmp/excel/rewardlist_%s.xlsx", time.Now().Format("20060102150405"))
	err := xlsx.SaveAs(fileName)
	if err != nil {
		beego.Error("Export reward error", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}
