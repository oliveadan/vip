package membertotal

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"html/template"
	"net/url"
	"os"
	"phage/utils"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
	"strconv"
	"strings"
	"time"
)

type MemberTotalIndexController struct {
	sysmanage.BaseController
}

func (this *MemberTotalIndexController) Get() {
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
	account := strings.TrimSpace(this.GetString("account"))
	level := strings.TrimSpace(this.GetString("level"))
	integral := strings.TrimSpace(this.GetString("integral"))
	keep := strings.TrimSpace(this.GetString("keep"))
	order := strings.TrimSpace(this.GetString("order"))

	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MemberTotal).Paginate(page, limit, account, level, integral, keep, order)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["condArr"] = map[string]interface{}{"account": account,
		"level":    level,
		"integral": integral,
		"keep":     keep,
		"order":    order}
	this.Data["dataList"] = list
	this.TplName = "common/membertotal/index.html"
}

func (this *MemberTotalIndexController) Count() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	s := this.GetString("start")
	e := this.GetString("end")
	if s == "" || e == "" {
		msg = "开始和结束时间不能为空"
		return
	}
	if s == e {
		msg = "开始和结束时间不能相同"
		return
	}
	start, _ := time.Parse("2006-01-02 15:04:05", s)
	end, _ := time.Parse("2006-01-02 15:04:05", e)
	if start.After(end) {
		msg = "开始时间不能晚于结束时间"
		return
	}
	if end.Before(start) {
		msg = "结束时间不能早于开始时间"
		return
	}
	o := orm.NewOrm()
	//获取要计算的所有会员的数据
	var mts []MemberTotal
	_, err := o.QueryTable(new(MemberTotal)).OrderBy("Bet").Limit(-1).All(&mts)
	if err != nil {
		msg = "获取会员数据失败"
		return
	}
	//用于改变会员保级状态
	var ids = make([]int64, 0)
	for _, v := range mts {
		//如果升级时间晚于开始时间跳过计算
		if v.LevelUpTime.After(start) {
			if v.KeepEnable != 0 {
				ids = append(ids, v.Id)
			}
			continue
		}
		//保级金额为0跳过配置
		var level Level
		err := o.QueryTable(new(Level)).Filter("VipLevel", v.Level).One(&level, "KeepLevelAmount", "KeepLevelDown", "TotalBet")
		if err != nil {
			beego.Error("获取会员%s等级失败", v.Account, err)
			continue
		}
		if level.KeepLevelAmount == 0 {
			if v.KeepEnable == 1 {
				ids = append(ids, v.Id)
			}
			continue
		}
		//对比投注和保级金额，判断会员是否达到保级要求
		var maps []orm.Params
		_, err1 := o.Raw("SELECT SUM(bet)  FROM ph_member_single WHERE account = ? AND create_date >= ? AND create_date <= ? ", v.Account, start, end).Values(&maps)
		if err1 != nil {
			beego.Error("获取会员%s所有金额失败", v.Account, err)
			continue
		}
		bet := maps[0]["SUM(bet)"]
		var bet1 int64
		if bet != nil {
			bet1, _ = strconv.ParseInt(bet.(string), 10, 64)
		}
		if bet == nil || bet1 < level.KeepLevelAmount {
			var levelbet Level
			err := o.QueryTable(new(Level)).Filter("VipLevel", level.KeepLevelDown).One(&levelbet, "TotalBet")
			if err != nil {
				beego.Error("获取%s的倒退等级投注额失败", v.Account, err)
				continue
			}
			o.Begin()
			_, err1 := o.QueryTable(new(MemberTotal)).Filter("Id", v.Id).Update(orm.Params{
				"Bet":         levelbet.TotalBet,
				"keepEnable":  1,
				"Level":       level.KeepLevelDown,
				"LevelUpTime": time.Now()})
			if err1 != nil {
				o.Rollback()
				beego.Error("更新会员%s失败", v.Account, err1)
			}
		} else {
			if v.KeepEnable == 1 {
				ids = append(ids, v.Id)
			}
		}
	}
	if len(ids) >= 1 {
		_, err1 := o.QueryTable(new(MemberTotal)).Filter("Id__in", ids).Update(orm.Params{"KeepEnable": 0})
		if err1 != nil {
			o.Rollback()
			beego.Error("更新保级状态失败", err1)
		}
	}
	o.Commit()
	num, _ := o.QueryTable(new(MemberTotal)).Filter("keepEnable", 1).Count()
	num1, _ := o.QueryTable(new(MemberTotal)).Filter("keepEnable", 0).Count()
	//重置提示
	_, _ = o.QueryTable(new(MemberTotal)).Update(orm.Params{"Tip": 0})
	code = 1
	msg = fmt.Sprintf("成功保级%d个,保级失败%d个", num1, num)
}

func (this *MemberTotalIndexController) Delbatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	o := orm.NewOrm()
	res, err := o.Raw("DELETE from ph_member_total WHERE id != 0").Exec()
	if err != nil {
		beego.Error("删除所的会员统计失败", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		num, _ := res.RowsAffected()
		msg = fmt.Sprintf("成功删除%d条记录", num)
		return
	}
}

func (this *MemberTotalIndexController) ChangeActivityStatus() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	i, _ := this.GetInt64("id", 0)
	status, _ := this.GetInt("status")
	mt := MemberTotal{Id: i}
	var s int
	if status == 0 {
		s = 1
	} else {
		s = 0
	}
	mt.ActivityStatus = s
	_, err := mt.Update("ActivityStatus")
	if err != nil {
		beego.Error("update ActivityStatus error", err)
		msg = "更新失败"
		return
	}
	code = 1
	msg = "更新成功"
}

func (this *MemberTotalIndexController) Export() {
	o := orm.NewOrm()
	var membertotal []MemberTotal
	_, err := o.QueryTable(new(MemberTotal)).OrderBy("-Bet").Limit(-1).All(&membertotal)
	if err != nil {
		beego.Error("导出会员统计列表失败", err)
	}
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "会员账号")
	xlsx.SetCellValue("Sheet1", "B1", "VIP等级")
	xlsx.SetCellValue("Sheet1", "C1", "会员积分")
	xlsx.SetCellValue("Sheet1", "D1", "升级时间")
	xlsx.SetCellValue("Sheet1", "E1", "保级状态")
	xlsx.SetCellValue("Sheet1", "F1", "时间奖励总收益")
	for i, value := range membertotal {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), value.Account)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), value.Level)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), value.Bet)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), value.LevelUpTime.Format("01/02/06 15:04"))
		if value.KeepEnable == 1 {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), "保级失败")
		} else {
			xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), "保级成功")
		}
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), value.TimeGiftSum)
	}
	fileName := fmt.Sprintf("./tmp/excel/membertotallist_%s.xlsx", time.Now().Format("20060102150405"))
	err1 := xlsx.SaveAs(fileName)
	if err1 != nil {
		beego.Error("导出会员列表失败", err.Error())
	} else {
		defer os.Remove(fileName)
		this.Ctx.Output.Download(fileName)
	}
}

func (this *MemberTotalIndexController) Import() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	f, h, err := this.GetFile("file")
	defer f.Close()
	if err != nil {
		beego.Error("导入会员统计失败", err)
		msg = "导入失败，请重试（1）"
		return
	}
	fname := url.QueryEscape(h.Filename)
	suffix := utils.SubString(fname, len(fname), strings.LastIndex(fname, ".")-len(fname))
	if suffix != ".xlsx" {
		msg = "文件必须为xlsx"
		return
	}

	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		beego.Error("会员统计导入失败", err)
		msg = "读取excel失败，请重试"
		return
	}
	notice := xlsx.GetSheetIndex("活动状态")
	count := xlsx.GetSheetIndex("会员统计")
	otherAccount := xlsx.GetSheetIndex("迁移会员")
	updateAccount := xlsx.GetSheetIndex("更新会员账号")
	if notice == 0 && count == 0 && otherAccount == 0 && updateAccount == 0 {
		msg = "不存在《会员统计》或《活动状态》、《迁移会员》、《更新会员账号》 "
		return
	}
	if count == 1 {
		models := make([]MemberTotal, 0)
		rows := xlsx.GetRows("会员统计")
		for i, row := range rows {
			if i == 0 {
				continue
			}
			var model MemberTotal
			account := strings.TrimSpace(row[0])
			if account == "" {
				msg = fmt.Sprintf("%s第%d行会员账号为空<br>", msg, i+1)
			}
			level := strings.TrimSpace(row[1])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行VIP等级为空<br>", msg, i+1)
			}
			bet := strings.TrimSpace(row[2])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行会员积分为空<br>", msg, i+1)
			}
			uptime := strings.TrimSpace(row[3])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行会员升级时间为空<br>", msg, i+1)
			}
			status := strings.TrimSpace(row[4])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行升级状态为空<br>", msg, i+1)
			}
			model.Account = account
			lev, _ := strconv.Atoi(level)
			model.Level = lev
			b, _ := strconv.ParseInt(bet, 10, 64)
			model.Bet = b
			ss, _ := time.Parse("01/02/06 15:04", uptime)
			model.LevelUpTime = ss
			if status == "保级失败" {
				model.KeepEnable = 1
			} else {
				model.KeepEnable = 0
			}
			model.Creator = this.LoginAdminId
			model.CreateDate = time.Now()
			model.ModifyDate = time.Now()
			models = append(models, model)
		}
		rlen := len(models)
		if rlen == 0 {
			msg = "没有需要导入的数据"
			return
		}
		var susNums int64
		o := orm.NewOrm()
		o.Begin()
		// 将数组拆分导入，一次1000条
		for i := 0; i <= rlen/1000; i++ {
			end := 0
			if (i+1)*1000 >= rlen {
				end = rlen
			} else {
				end = (i + 1) * 1000
			}
			if i*1000 == end {
				continue
			}
			tmpArr := models[i*1000 : end]
			if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
				o.Rollback()
				beego.Error("会员投注记录导入失败", err)
			} else {
				susNums += nums
			}
		}
		o.Commit()
		if msg != "" {
			msg = fmt.Sprintf("请处理以下错误后再导入：<br>%s", msg)
			return
		}
		code = 1
		msg = fmt.Sprintf("%s成功导入%d条记录", msg, susNums)
		return
	}
	if notice == 1 {
		rows := xlsx.GetRows("活动状态")
		var accounts1 []string
		var accounts0 []string
		for i, row := range rows {
			if i == 0 {
				continue
			}
			account := strings.TrimSpace(row[0])
			if account == "" {
				msg = fmt.Sprintf("%s第%d行会员账号为空<br>", msg, i+1)
			}
			status := strings.TrimSpace(row[1])
			if status == "" {
				msg = fmt.Sprintf("%s第%d行活动状态为空等级为空<br>", msg, i+1)
			}
			if status == "1" {
				accounts1 = append(accounts1, account)
			} else {
				accounts0 = append(accounts0, account)
			}
		}
		if len(accounts1) == 0 && len(accounts0) == 0 {
			msg = "没有需要更新的的数据"
			return
		}
		if len(accounts1) != 0 {
			o := orm.NewOrm()
			o.Begin()
			i, err := o.QueryTable(new(MemberTotal)).Filter("Account__in", accounts1).Update(orm.Params{"ActivityStatus": 1})
			if err != nil {
				o.Rollback()
				beego.Error("update membertotal eorror", err)
				msg = "更新活动状态失败"
				return
			}
			o.Commit()
			msg = fmt.Sprintf("成功更新%d个会员状态为禁用", i)
		}
		if len(accounts0) != 0 {
			o := orm.NewOrm()
			o.Begin()
			i, err := o.QueryTable(new(MemberTotal)).Filter("Account__in", accounts0).Update(orm.Params{"ActivityStatus": 0})
			if err != nil {
				o.Rollback()
				beego.Error("update membertotal eorror", err)
				msg = "更新活动状态失败"
				return
			}
			o.Commit()
			msg = fmt.Sprintf("%s<br>成功更新%d个会员状态为正常", msg, i)
		}
	}
	if otherAccount == 1 {
		o := orm.NewOrm()
		models := make([]MemberTotal, 0)
		rows := xlsx.GetRows("迁移会员")
		var levels []Level
		_, err := o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)
		if err != nil {
			beego.Error("获取VIP等级失败", err)
			msg = "VIP等级获取失败，请检查VIP等级配置"
			return
		}
		for i, row := range rows {
			if i == 0 {
				continue
			}
			var model MemberTotal
			account := strings.TrimSpace(row[0])
			if account == "" {
				msg = fmt.Sprintf("%s第%d行会员账号为空<br>", msg, i+1)
			}
			exist := o.QueryTable(new(MemberTotal)).Filter("Account", account).Exist()
			if exist {
				continue
			}
			level := strings.TrimSpace(row[1])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行VIP等级为空<br>", msg, i+1)
			}
			bet := strings.TrimSpace(row[2])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行会员积分为空<br>", msg, i+1)
			}
			uptime := strings.TrimSpace(row[3])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行会员升级时间为空<br>", msg, i+1)
			}
			status := strings.TrimSpace(row[4])
			if level == "" {
				msg = fmt.Sprintf("%s第%d行升级状态为空<br>", msg, i+1)
			}
			model.Account = account
			lev, _ := strconv.Atoi(level)
			model.Level = lev
			b, _ := strconv.ParseInt(bet, 10, 64)
			model.Bet = b
			ss, _ := time.Parse("01/02/06 15:04", uptime)
			model.LevelUpTime = ss
			if status == "保级失败" {
				model.KeepEnable = 1
			} else {
				model.KeepEnable = 0
			}
			model.Creator = this.LoginAdminId
			model.CreateDate = time.Now()
			model.ModifyDate = time.Now()
			model.GetGiftTime = time.Now()
			models = append(models, model)
			//所有等级默认已领取奖励
			for i := 1; i <= lev; i++ {
				for _, p := range levels {
					if i == p.VipLevel {
						ml := MemberLevelLog{Level: p.VipLevel, Account: account, LevelGift: p.LevelGift, EnAble: 1}
						_, _, err := o.ReadOrCreate(&ml, "Level", "Account")
						if err != nil {
							beego.Error("插入VIP失败", err)
							msg = "系统异常2"
							return
						}
					}
				}
			}
		}
		rlen := len(models)
		if rlen == 0 {
			msg = "没有需要导入的数据"
			return
		}
		var susNums int64
		o.Begin()
		// 将数组拆分导入，一次1000条
		for i := 0; i <= rlen/1000; i++ {
			end := 0
			if (i+1)*1000 >= rlen {
				end = rlen
			} else {
				end = (i + 1) * 1000
			}
			if i*1000 == end {
				continue
			}
			tmpArr := models[i*1000 : end]
			if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
				o.Rollback()
				beego.Error("迁移会员投注记录导入失败", err)
				msg = "迁移会员投注记录导入失败" + fmt.Sprintf("%v", err)
			} else {
				susNums += nums
			}
		}
		if msg != "" {
			msg = fmt.Sprintf("请处理以下错误后再导入：<br>%s", msg)
			return
		}
		o.Commit()
		code = 1
		msg = fmt.Sprintf("%s成功导入%d条迁移会员记录", msg, susNums)
		return
	}
	if updateAccount == 1 {
		rows := xlsx.GetRows("更新会员账号")
		var accouts []string
		var identification string
		for i, row := range rows {
			if i == 0 {
				if len(row) < 3 {
					msg = "标识不能为空"
					return
				}
				identification = strings.TrimSpace(row[2])
				if identification == "" {
					msg = "标识不能为空"
					return
				}
				continue
			}
			account := strings.TrimSpace(row[0])
			if account != "" {
				accouts = append(accouts, account)
			}
		}
		if len(accouts) == 0 {
			msg = "没有需要更新的数据"
			return
		}
		qb, _ := orm.NewQueryBuilder("mysql")
		questions := question(len(accouts))
		qb.Update("ph_member_total").Set("account = CONCAT(?,account)").Where("account in(" + questions + ")")
		sql := qb.String()
		o := orm.NewOrm()
		result, err := o.Raw(sql, identification, accouts).Exec()
		i, _ := result.RowsAffected()
		if err != nil {
			o.Rollback()
			beego.Error("update ph_member_total  error", err)
			msg = "更新失败"
			return
		}
		o.Commit()
		code = 1
		msg = fmt.Sprintf("成功更新了%d条数据", i)
	}
}

type MemberTotalEditController struct {
	sysmanage.BaseController
}

func (this *MemberTotalEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	mt := MemberTotal{Id: id}
	err := o.Read(&mt)
	if err != nil {
		this.Redirect("membertotaleIndexController.get", 302)
	} else {
		this.Data["data"] = mt
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "common/membertotal/edit.html"
	}
}

func (this *MemberTotalEditController) Post() {
	var code int
	var msg string
	url := beego.URLFor("MembertotalIndexController.get")
	defer sysmanage.Retjson(this.Ctx, &msg, &code, &url)
	o := orm.NewOrm()
	mt := MemberTotal{}
	if err := this.ParseForm(&mt); err != nil {
		beego.Error("更新会员总投注信息异常", err)
		msg = "参数异常"
		return
	}
	cols := []string{"Level", "Bet"}
	//修改隐藏表的VIP等级
	//获取VIP等级
	var levels []Level
	_, er := o.QueryTable(new(Level)).OrderBy("-TotalBet").All(&levels)
	if er != nil {
		beego.Error("获取VIP等级失败", er)
		msg = "VIP等级获取失败，请检查VIP等级配置"
		return
	}
	exist := o.QueryTable(new(MemberLevelLog)).Filter("Account", mt.Account).Filter("Level", mt.Level).Exist()
	if !exist {
		for i := 1; i <= mt.Level; i++ {
			for _, p := range levels {
				if i == p.VipLevel {
					ml := MemberLevelLog{Level: p.VipLevel, Account: mt.Account, LevelGift: p.LevelGift}
					_, _, err := o.ReadOrCreate(&ml, "Level", "Account")
					if err != nil {
						beego.Error("插入VIP失败", err)
						msg = "系统异常5"
						o.Rollback()
						return
					}
				}
			}
		}
	}
	mt.Modifior = this.LoginAdminId
	_, err := o.Update(&mt, cols...)
	if err != nil {
		beego.Error("更新会员总投注信息失败", err)
		msg = "更新失败"
		return
	} else {
		code = 1
		msg = "更新成功"
	}
}

func question(len int) string {
	var question string
	for i := 0; i < len; i++ {
		if i == len-1 {
			question += "?"
		} else {
			question += "?,"
		}
	}
	return question
}
