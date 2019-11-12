package front

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"math"
	utils2 "phagego/common/utils"
	"phagego/frameweb-v2/controllers/sysmanage"
	"phagego/frameweb-v2/models"
	"phagego/phage-vip4-web/initial"
	. "phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type VipCenterController struct {
	sysmanage.BaseController
}

func (this *VipCenterController) Prepare() {
	this.EnableXSRF = false
}

func (this *VipCenterController) Get() {
	name := this.GetString("name")
	var mll MemberLevelLog
	o := orm.NewOrm()
	//会员可领取的vip等级
	_ = o.QueryTable(new(MemberLevelLog)).Filter("Account", name).OrderBy("-Level").One(&mll)
	var colorlevel []Level
	_, _ = o.QueryTable(new(Level)).Filter("VipLevel__lte", mll.Level).Filter("VipLevel__gt", 0).All(&colorlevel)
	//判断会员奖品是否已经领取
	var mlls []MemberLevelLog
	_, _ = o.QueryTable(new(MemberLevelLog)).Filter("Account", name).OrderBy("Level").All(&mlls)
	//会员不可领取的vip等级
	var wblevel []Level
	_, _ = o.QueryTable(new(Level)).Filter("VipLevel__gt", mll.Level).All(&wblevel)
	//取会员总信息用于前台展示
	var mt MemberTotal
	_ = o.QueryTable(new(MemberTotal)).Filter("Account", name).One(&mt)
	//会员当前VIP等级信息
	var level Level
	_ = o.QueryTable(new(Level)).Filter("VipLevel", mt.Level).One(&level)

	//下一个VIP等级
	var lev Level
	_ = o.QueryTable(new(Level)).Filter("VipLevel", mt.Level+1).One(&lev)

	//当前投注所占百分比
	var bili float64
	bili = float64(mt.Bet) / float64(lev.TotalBet)
	zs := math.Ceil(bili * 100)
	//距离下一级所需的投注占的百分比
	zs1 := 100 - zs
	//快速导航
	var quciknav []models.QuickNav
	_, _ = o.QueryTable(new(models.QuickNav)).OrderBy("Seq").Limit(4).All(&quciknav)
	//计算会员VIP天数
	subDays := int(time.Now().Sub(mt.CreateDate).Hours()) / 24
	//获取会员的领取时间
	var va VipAttribute
	_ = o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGifttime).One(&va)
	duration, _ := time.ParseDuration(va.Value + "m")
	add := mt.GetGiftTime.Add(duration)
	subgifttime := add.Unix() - time.Now().Unix()
	if subgifttime > 0 || mt.GetGiftTime.IsZero() {
		this.Data["getGiftTimeStatus"] = 0
	}
	//获取会员中奖记录
	var rw []RewardLog
	var sumtimegift float64
	_, rwerr := o.QueryTable(new(RewardLog)).Filter("Account", mt.Account).Filter("Category", 1).OrderBy("-CreateDate").All(&rw)
	if rwerr != nil {
		beego.Error("query rw error", rwerr)
	} else {
		for _, v := range rw {
			s, _ := strconv.ParseFloat(v.GiftContent, 32)
			sumtimegift += s
		}
		for _, v := range rw {
			this.Data["nowtimegift"] = v.GiftContent
			break
		}
	}
	//获取公告图片
	noticeimg := initial.GetNoticeImg()
	sc := models.GetSiteConfigMap(utils.Scnotice)
	this.Data["noticeimg"] = noticeimg
	this.Data["sumtimegift"] = fmt.Sprintf("%1.2f", sumtimegift)
	this.Data["getGiftTime"] = subgifttime
	this.Data["subdays"] = subDays
	this.Data["nav"] = quciknav
	this.Data["fontmt"] = mt
	this.Data["bili"] = zs
	this.Data["bili1"] = zs1
	this.Data["mlls"] = mlls
	this.Data["level"] = level
	this.Data["colorlevel"] = colorlevel
	this.Data["wblevel"] = wblevel
	this.Data["blance"] = lev.TotalBet - mt.Bet
	this.Data["nextlevel"] = lev.VipLevel
	this.Data["notice"] = sc[utils.Scnotice]
	this.TplName = "front/vipcenter.html"
}

func (this *VipCenterController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id := this.GetString("id")
	o := orm.NewOrm()
	var mll MemberLevelLog
	//判断会员是否已经领取彩金
	_ = o.QueryTable(new(MemberLevelLog)).Filter("Id", id).One(&mll)
	if mll.EnAble == 1 {
		msg = "您已经成功领取彩金了~"
		return
	}

	var rewardlog RewardLog
	rewardlog.Account = mll.Account
	rewardlog.GiftName = fmt.Sprintf("VIP%d晋级奖励", mll.Level)
	rewardlog.GiftContent = fmt.Sprintf("%d", mll.LevelGift)
	var lock sync.RWMutex
	lock.Lock()
	_, err := rewardlog.Create()
	lock.Unlock()
	if err != nil {
		beego.Error("生成中奖记录失败", err)
		msg = "系统异常，请刷新后重试"
		return
	} else {
		_, err := o.QueryTable(new(MemberLevelLog)).Filter("Id", id).Update(orm.Params{"Enable": 1})
		if err != nil {
			beego.Error("更新晋级礼物失败", err)
		}
	}
	code = 1
	msg = fmt.Sprintf("恭喜您获VIP%d晋级奖励%d元", mll.Level, mll.LevelGift)
}

func (this *VipCenterController) GetTimeGift() {
	var code int
	var msg string
	account := this.GetString("account")
	data := make(map[string]interface{})
	defer Retjson(this.Ctx, &msg, &code, &data)
	o := orm.NewOrm()
	//查询活动是否开启
	var vb VipAttribute
	_ = o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGiftStatus).One(&vb)
	if vb.Value == "0" {
		msg = "活动暂未开启"
		return
	}

	var mt MemberTotal
	one := o.QueryTable(new(MemberTotal)).Filter("Account", account).One(&mt)
	if one != nil {
		beego.Error("query MemberTotal error", one)
		msg = "领取失败请联系客服处理"
		return
	}
	if mt.ActivityStatus == 1 {
		var notice VipAttribute
		o.QueryTable(new(VipAttribute)).Filter("code", utils.TimeGiftNotice).One(&notice)
		if notice.Value != "" {
			msg = notice.Value
		} else {
			msg = "您的账号异常,请咨询客服"
		}
		return
	}
	//验证此时间段奖品是否已经领取
	var va VipAttribute
	_ = o.QueryTable(new(VipAttribute)).Filter("Code", utils.TimeGifttime).One(&va)
	duration, _ := time.ParseDuration(va.Value + "m")
	add := mt.GetGiftTime.Add(duration)
	subgifttime := add.Unix() - time.Now().Unix()
	if subgifttime > 0 {
		msg = "您已经领取过这个时间段的奖励了"
		return
	}
	//获取奖品
	var tg TimeGift
	err := o.QueryTable(new(TimeGift)).Filter("GiftLevel", mt.Level).One(&tg)
	if err != nil {
		msg = "当前等级没有奖品"
		return
	}
	//生成中奖记录
	var rewardlog RewardLog
	n1 := utils2.RandNum(int(tg.MinMoney), int(tg.MaxMoney-1))
	n2 := utils2.RandNum(0, 99)
	var amount string
	if tg.Category == 0 {
		amount = fmt.Sprintf("%d.%d", n1, n2)
	} else {
		n3 := utils2.RandNum(int(tg.MinMoney), int(tg.MaxMoney-1))
		amount = fmt.Sprintf("0.%d%d", n1, n3)
	}

	rewardlog.CreateDate = time.Now()
	rewardlog.Account = mt.Account
	rewardlog.GiftName = time.Now().Format("2006-01-02 15:04:05") + "-奖励"
	rewardlog.GiftContent = amount
	rewardlog.Category = 1
	gift, _ := strconv.ParseFloat(amount, 64)
	sum := mt.TimeGiftSum + gift
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", sum), 64)
	o.Begin()
	//更新会员领取时间
	_, err3 := o.QueryTable(new(MemberTotal)).Filter("Id", mt.Id).Update(orm.Params{"GetGiftTime": time.Now(),
		"TimeGiftSum": f})
	if err3 != nil {
		beego.Error("update memebertotal error", err3)
		o.Rollback()
		msg = "领取失败请联系客服处理(4)"
		return
	}
	_, err1 := rewardlog.Create()
	if err1 != nil {
		beego.Error("insert RewardLog error", err1)
		o.Rollback()
		msg = "领取失败请联系客服处理(3)"
		return
	}
	o.Commit()
	code = 1
	data["gift"] = amount
	msg = "恭喜您奖励领取成功"
	return
}

func (this *VipCenterController) QueryTimeGift() {
	var code int
	var msg string
	data := make(map[string]interface{})
	account := this.GetString("account")
	defer Retjson(this.Ctx, &msg, &code, &data)
	var rws []RewardLog
	o := orm.NewOrm()
	exist := o.QueryTable(new(MemberTotal)).Filter("Account", account).Exist()
	if !exist {
		msg = "您输入的会员账号不存在"
		return
	}
	_, e := o.QueryTable(new(RewardLog)).Filter("Account", account).Filter("Category", 1).All(&rws)
	if e != nil {
		beego.Error("query rewardlog error", e)
		msg = "查询失败(1)"
		return
	}
	type gift struct {
		Time time.Time
		Name string
	}
	var gifts []gift
	for _, v := range rws {
		g := gift{}
		g.Time = v.CreateDate
		g.Name = v.GiftContent
		gifts = append(gifts, g)
	}
	data["gifts"] = gifts
	code = 1
	msg = "查询成功"
}

func (this *VipCenterController) QueryPrivilege() {
	name := this.GetString("name")
	this.Data["name"] = name
	this.TplName = "front/vipdetail.html"
}

func (this *VipCenterController) ChangeTip() {
	account := this.GetStrings("account")
	o := orm.NewOrm()
	var mt MemberTotal
	_ = o.QueryTable(new(MemberTotal)).Filter("Account", account).One(&mt)
	if mt.Tip == 0 {
		_, _ = o.QueryTable(new(MemberTotal)).Filter("Account", account).Update(orm.Params{"Tip": 1})
	}
	this.Ctx.Output.JSON("", false, false)
}

func (this *VipCenterController) Mission() {
	this.AllowCross()
	var msg string
	var code int
	data := make(map[string]interface{})
	account := this.GetString("account")
	level := this.GetString("level")
	defer retjson(&this.Controller, &msg, &code, &data)
	if account == "" {
		msg = "会员账号不能为空"
		return
	}
	if level == "" {
		msg = "VIP等级不能为空"
		return
	}
	o := orm.NewOrm()
	//任务
	timenow := time.Now().Format("2006-01-02 15:04:05")
	//查找开启的活动
	var ms []Mission
	_, e := o.QueryTable(new(Mission)).Filter("StartTime__lte", timenow).Filter("EndTime__gte", timenow).All(&ms)
	if e != nil {
		beego.Error("get Mission fault", e)
		msg = "系统异常(1),请刷新后尝试"
		return
	}
	var missionids []int64
	for _, v := range ms {
		missionids = append(missionids, v.Id)
	}
	if len(missionids) < 1 {
		missionids = append(missionids, 0)
	}
	//查询开启活动的详情
	var missiondetails []MissionDetail
	_, err := o.QueryTable(new(MissionDetail)).Filter("MissionId__in", missionids).Filter("MinLevel__lte", level).Filter("MaxLevel__gte", level).All(&missiondetails)
	if err != nil {
		beego.Error("get missdiondetails fault", err)
		msg = "系统异常(3),请刷新后尝试"
		return
	}
	//只查询匹配的活动
	var mds []int64
	for _, v := range missiondetails {
		mds = append(mds, v.MissionId)
	}
	var mss []Mission
	if len(mds) >= 1 {
		_, err0 := o.QueryTable(new(Mission)).Filter("Id__in", mds).All(&mss)
		if err0 != nil {
			beego.Error("query Mission fault(1)", err)
		}
	}

	//前台展示数据
	type result struct {
		MissionId             int64
		MissionDescribe       string
		MissionDetailId       int64
		MissionDetailDescribe string
		MinLevel              int64
		MaxLevel              int64
		Integral              int64
		Remark                string
		CreateDate            time.Time
		Status                int
	}
	var array []result
	var missionReview MissionReview
	for _, v := range mss {
		r := result{}
		for _, j := range missiondetails {
			if j.MissionId == v.Id {
				r.MissionId = v.Id
				r.MissionDetailId = j.Id
				r.MinLevel = j.MinLevel
				r.MaxLevel = j.MaxLevel
				r.Integral = v.Integral
				r.MissionDescribe = v.Describe
				r.MissionDetailDescribe = j.Award
				//查询会员最新的中奖记录
				err1 := o.QueryTable(new(MissionReview)).Filter("Account", account).Filter("MissionId", v.Id).Filter("MissionDetailId", j.Id).OrderBy("-CreateDate").One(&missionReview)
				if err1 != nil {
					r.Status = -1
				} else {
					r.Status = missionReview.Status
					r.Remark = missionReview.Remark
					r.CreateDate = missionReview.CreateDate
				}
				array = append(array, r)
				break
			}
		}
	}
	code = 1
	msg = "查询成功"
	data["result"] = array
}

func (this *VipCenterController) CreateMissionReview() {
	this.AllowCross()
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	account := this.GetString("account")
	missionId, _ := this.GetInt64("missionid")
	missionDetailId, _ := this.GetInt64("missioindetailid")
	minLevel, _ := this.GetInt64("minlevel")
	maxLevel, _ := this.GetInt64("maxlevel")
	integral, _ := this.GetInt64("integral")
	if account == "" {
		msg = "会员账号不能为空"
		return
	}
	if missionId == 0 || missionDetailId == 0 || minLevel == 0 || maxLevel == 0 {
		msg = "系统异常请刷新后重试"
		return
	}
	o := orm.NewOrm()
	exist := o.QueryTable(new(MissionReview)).Filter("Account", account).Filter("MissionId", missionId).Filter("MissionDetailId", missionDetailId).Filter("MinLevel", minLevel).Filter("MaxLevel", maxLevel).Filter("Status__in", 1, 0).Exist()
	if exist {
		msg = "您的申请已经提交的了,请耐心等待专员的审核"
		return
	}
	var mr MissionReview
	mr.Account = account
	mr.MissionId = missionId
	mr.MissionDetailId = missionDetailId
	mr.MinLevel = minLevel
	mr.MaxLevel = maxLevel
	mr.Integral = integral
	_, err1 := mr.Create()
	if err1 != nil {
		beego.Error("create MissionReview fault", err1)
		msg = "系统异常请刷新后重试(2)"
		return
	}
	code = 1
	msg = "已提交成功,专员审核通过后会为您派奖"
}

func (this *VipCenterController) Remark() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	content := strings.TrimSpace(this.GetString("content"))
	mr := MissionReview{Id: id}
	mr.Remark = content
	_, e := mr.Update("Remark")
	if e != nil {
		beego.Error("update MissionReview Remark fault", e)
		msg = "备注更新失败"
		return
	}
	code = 1
	msg = "备注更新成功"
}

func retjson(c *beego.Controller, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	ret["data"] = data
	c.Data["json"] = &ret
	c.ServeJSON()
}

func Retjson(ctx *context.Context, msg *string, code *int, data ...interface{}) {
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	if len(data) > 0 {
		ret["data"] = data[0]
	}
	ctx.Output.JSON(ret, false, false)
}
