package missiondata

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/pagination"
	"net/url"
	"phage/controllers/sysmanage"
	utils2 "phagego/common/utils"
	. "phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
	"strconv"
	"strings"
	"time"
)

type IndexMissionDateController struct {
	sysmanage.BaseController
}

func (this *IndexMissionDateController) Get() {
	page, err := this.GetInt("p")
	if err != nil {
		page = 1
	}
	missionid, _ := this.GetInt64("id")
	starttime := strings.TrimSpace(this.GetString("startTime"))
	endtime := strings.TrimSpace(this.GetString("endTime"))
	account := strings.TrimSpace(this.GetString("account"))
	//获取期数
	var period MissionDate
	o := orm.NewOrm()
	err1 := o.QueryTable(new(MissionDate)).Filter("Period__isnull", false).Filter("MissionId", missionid).OrderBy("-Period").Distinct().Limit(-1).One(&period, "Period")
	if err1 != nil {
		beego.Error("获取任务数据期数失败", err1)
	}
	if period.Period != "" && starttime == "" {
		starttime = period.Period
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(MissionDate).Paginate(page, limit, missionid, starttime, endtime, account)
	pagination.SetPaginator(this.Ctx, limit, total)
	this.Data["condArr"] = map[string]interface{}{"starttime": starttime,
		"endtime": endtime,
		"account": account}
	this.Data["dataList"] = list
	this.Data["Missionid"] = missionid
	//任务内容
	this.Data["describe"] = utils.GetMissionDescribe(missionid)
	this.TplName = "common/missiondata/index.html"
}

func (this *IndexMissionDateController) Delbatch() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	starttime := strings.TrimSpace(this.GetString("startTime"))
	endtime := strings.TrimSpace(this.GetString("endTime"))
	missionid, _ := this.GetInt64("id")
	if starttime == "" || endtime == "" {
		msg = "请选择要删除的时间区间"
		return
	}
	o := orm.NewOrm()
	num, err := o.QueryTable(new(MissionDate)).Filter("MissionId", missionid).Filter("Period__gte", starttime).Filter("Period__lte", endtime).Delete()
	if err != nil {
		beego.Error("删除任务数据失败", err)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = fmt.Sprintf("成功删除%d条数据", num)
	}
}

func (this *IndexMissionDateController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	md := MissionDate{Id: id}
	o := orm.NewOrm()
	err := o.Read(&md)
	if err == orm.ErrMissPK || err == orm.ErrNoRows {
		this.Redirect("IndexMissionDateController.get", 302)
	}
	_, err1 := o.Delete(&md, "Id")
	if err1 != nil {
		beego.Error("删除任务数据失败", err1)
		msg = "删除失败"
		return
	} else {
		code = 1
		msg = "删除成功"
	}
}

func (this *IndexMissionDateController) Count() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	//获取需期数
	starttime := this.GetString("startTime")
	endtime := this.GetString("endTime")
	missionid, _ := this.GetInt64("id")
	o := orm.NewOrm()
	//获取任务配置详情
	var missiondetails []MissionDetail
	_, err := o.QueryTable(new(MissionDetail)).Filter("MissionId", missionid).OrderBy("-MinLevel", "-Content").All(&missiondetails)
	if err != nil {
		beego.Error("获取配置详情失败", err)
	}
	//获取任务配置（判断是否是累计计算）
	var mission Mission
	errr := o.QueryTable(new(Mission)).Filter("Id", missionid).One(&mission)
	if errr != nil {
		beego.Error("获取任务配置失败", errr)
		msg = "获取任务配置失败"
		return
	}

	var min int64
	for _, v := range missiondetails {
		min = v.Content
	}
	//获取要计算的任务数据
	var mds []MissionDate
	var lists []orm.ParamsList
	var missionresults []MissionResult
	var ids []int64
	var accounts []string
	//判断是否是累加计算
	if mission.SumEnable == 1 {
		//判断是否有可计算数据
		exist := o.QueryTable(new(MissionDate)).Filter("MissionId", missionid).Filter("Enable", 0).Exist()
		if !exist {
			msg = "没有可计算的数据"
			return
		}
		_, err := o.Raw("SELECT account ,sum(data) FROM ph_mission_date where mission_id = ? and account in(SELECT account FROM ph_mission_date WHERE enable = 0 AND mission_id = ?)  GROUP BY account", missionid, missionid).ValuesList(&lists)
		if err != nil {
			beego.Error("获取累计计算数据失败", err)
			msg = "获取累计计算数据失败"
			return
		}
		if len(lists) == 0 {
			msg = "没有要计算的数据"
			return
		}
		//开始计算
		for _, v := range lists {
			var missiondates []MissionDate
			_, err1 := o.QueryTable(new(MissionDate)).Filter("MissionId", missionid).Filter("Enable", 0).Filter("Account", v[0]).All(&missiondates, "Id")
			if err1 != nil {
				beego.Error("获取ID失败", err1)
			}
			for _, h := range missiondates {
				ids = append(ids, h.Id)
			}

			data, _ := strconv.ParseInt(fmt.Sprintf("%s", v[1]), 10, 64)
			if data < min {
				continue
			}
			//获取会员当前VIP等级
			var level MemberTotal
			err := o.QueryTable(new(MemberTotal)).Filter("Account", v[0]).One(&level, "Level")
			if err != nil {
				beego.Info(v[0], "未查询到VIP等级", err)
				continue
			}
			for _, j := range missiondetails {
				if int64(level.Level) >= j.MinLevel && int64(level.Level) <= j.MaxLevel && data >= j.Content {
					var model MissionResult
					model.CreateDate = time.Now()
					model.ModifyDate = time.Now()
					model.Modifior = this.LoginAdminId
					model.Creator = this.LoginAdminId
					model.Version = 0
					model.MissionId = missionid
					model.Account = v[0].(string)
					model.Prize = j.Award
					model.Enable = 0
					missionresults = append(missionresults, model)
					//收集会员账号，更新积分
					accounts = append(accounts, v[0].(string))
					break
				}
			}
		}
	} else {
		if starttime == "" || endtime == "" {
			msg = "请选择要都计算的区间"
			return
		}
		_, err1 := o.QueryTable(new(MissionDate)).Filter("Enable", 0).Filter("Period__gte", starttime).Filter("Period__lte", endtime).Filter("MissionId", missionid).All(&mds)
		if err1 != nil {
			beego.Error("获取要计算的任务数据失败", err1)
			msg = "获取要计算的任务数据失败"
			return
		}
		if len(mds) == 0 {
			msg = "没有可计算的数据"
			return
		}
		//开始计算
		for _, v := range mds {
			ids = append(ids, v.Id)
			if v.Data < min {
				continue
			}
			//获取会员当前VIP等级
			var level MemberTotal
			err := o.QueryTable(new(MemberTotal)).Filter("Account", v.Account).One(&level, "Level")
			if err != nil {
				beego.Info(v.Account, "未查询到VIP等级", err)
				continue
			}
			for _, j := range missiondetails {
				if int64(level.Level) >= j.MinLevel && int64(level.Level) <= j.MaxLevel && v.Data >= j.Content {
					var model MissionResult
					model.CreateDate = time.Now()
					model.ModifyDate = time.Now()
					model.Modifior = this.LoginAdminId
					model.Creator = this.LoginAdminId
					model.Version = 0
					model.MissionId = missionid
					model.Account = v.Account
					model.Prize = j.Award
					model.Enable = 0
					missionresults = append(missionresults, model)
					//收集会员账号，更新积分
					accounts = append(accounts, v.Account)
					break
				}
			}
		}
	}
	//在计算后生成会员统计列表
	var susNums int64
	//将数组拆分导入，一次1000条
	o.Begin()
	mlen := len(missionresults)
	if mlen > 0 {
		for i := 0; i <= mlen/1000; i++ {
			end := 0
			if (i+1)*1000 >= mlen {
				end = mlen
			} else {
				end = (i + 1) * 1000
			}
			if i*1000 == end {
				continue
			}
			tmpArr := missionresults[i*1000 : end]
			if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
				o.Rollback()
				beego.Error("插入计算结果失败", err)
				return
			} else {
				susNums += nums
			}
		}
	}
	//更新为已计算
	_, err6 := o.QueryTable(new(MissionDate)).Filter("Id__in", ids).Update(orm.Params{"Enable": 1})
	if err6 != nil {
		beego.Error("标记已计算失败", err6)
		o.Rollback()
	}
	//更新会员积分
	_, err13 := o.QueryTable(new(MemberTotal)).Filter("Account__in", accounts).Update(orm.Params{"MissionIntegral": orm.ColValue(orm.ColAdd, mission.Integral)})
	if err13 != nil {
		beego.Error("更新会员积分失败", err13)
		o.Rollback()
	}
	o.Commit()
	code = 1
	msg = fmt.Sprintf("计算成功，生成%d条中奖记录", susNums)
}

func (this *IndexMissionDateController) Import() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	f, h, err := this.GetFile("file")
	misisonid, _ := this.GetInt64("missionid")
	defer f.Close()
	if err != nil {
		beego.Error("导入会员投注失败", err)
		msg = "导入失败，请重试（1）"
		return
	}
	fname := url.QueryEscape(h.Filename)
	suffix := utils2.SubString(fname, len(fname), strings.LastIndex(fname, ".")-len(fname))
	if suffix != ".xlsx" {
		msg = "文件必须为xlsx"
		return
	}

	o := orm.NewOrm()
	missiondatas := make([]MissionDate, 0)
	xlsx, err := excelize.OpenReader(f)
	if err != nil {
		beego.Error("任务数据导入失败", err)
		msg = "读取excel失败，请重试"
		return
	}
	if xlsx.GetSheetIndex("任务数据") == 0 {
		msg = "不存在《任务数据》sheet页"
		return
	}
	rows := xlsx.GetRows("任务数据")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 1 {
			msg = fmt.Sprintf("%s第%d行为空<br>", msg, i+1)
			continue
		}
		var md MissionDate
		account := strings.TrimSpace(row[0])
		if account == "" {
			msg = fmt.Sprintf("%s第%d行会员账号为空<br>", msg, i+1)
		}
		data := strings.TrimSpace(row[1])
		if data == "" {
			msg = fmt.Sprintf("%s第%d行任务数据为空<br>", msg, i+1)
		} else {
			data1, _ := strconv.ParseInt(data, 10, 64)
			md.Data = data1
		}
		// 如果会员存在就跳过
		bool := o.QueryTable(new(MissionDate)).Filter("Account", account).Filter("Missionid", misisonid).Filter("Period", time.Now().Format("2006-01-02 15:04:05")).Exist()
		if bool {
			continue
		}
		md.CreateDate = time.Now()
		md.ModifyDate = time.Now()
		md.Modifior = this.LoginAdminId
		md.Creator = this.LoginAdminId
		md.Account = account
		md.MissionId = misisonid
		md.Period = time.Now().Format("2006-01-02 15:04:05")
		missiondatas = append(missiondatas, md)
	}
	if msg != "" {
		msg = fmt.Sprintf("请处理以下错误后再导入：<br>%s", msg)
		return
	}
	rlen := len(missiondatas)
	if rlen == 0 {
		msg = "没有需要导入的数据"
		return
	}
	var susNums int64
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
		tmpArr := missiondatas[i*1000 : end]
		if nums, err := o.InsertMulti(len(tmpArr), tmpArr); err != nil {
			beego.Error("任务数据导入失败", err)
		} else {
			susNums += nums
		}
	}
	code = 1
	msg = fmt.Sprintf("%s成功导入%d条任务数据", msg, susNums)
	return
}
