package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "phagego/frameweb-v2/utils"
	"phagego/phage-vip4-web/models/common"
	"phagego/phage-vip4-web/utils"
)

var static_front = beego.AppConfig.String("staticfront")

func initTemplateFunc() {
	beego.AddFuncMap("getSiteConfigCodeMap", GetSiteConfigCodeMap)
	beego.AddFuncMap("getCountEnAble", utils.GetCountEnAble)
	beego.AddFuncMap("getMissionDescribe", utils.GetMissionDescribe)
	beego.AddFuncMap("getMissionDetail", utils.GetMissionDetail)
	beego.AddFuncMap("getSumTimeGift", utils.GetSumTimeGift)
	beego.AddFuncMap("getNoticeImg", GetNoticeImg)
	beego.AddFuncMap("static_front", getStaticFront)
}

func getStaticFront() string {
	return static_front
}

func GetSiteConfigCodeMap() map[string]string {
	m := map[string]string{
		Scname:            "站点名称",
		utils.Scofficial:  "官网网址",
		utils.Scranking:   "排行榜网址",
		utils.Scregister:  "官网注册",
		utils.Sccust:      "在线客服",
		utils.Scfqa:       "博彩责任",
		utils.Scpromotion: "优惠活动",
		utils.Scnotice:    "网站公告"}
	return m
}

//获取公告图片
func GetNoticeImg() string {
	o := orm.NewOrm()
	var va common.VipAttribute
	err := o.QueryTable(new(common.VipAttribute)).Filter("Code", "noticeimg").One(&va)
	if err != nil {
		return ""
	}
	return va.Value
}
