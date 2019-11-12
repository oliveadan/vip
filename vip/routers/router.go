package routers

import (
	"github.com/astaxie/beego"
	_ "phagego/frameweb-v2/routers"
	"phagego/phage-vip4-web/controllers/common/imgupload"
	"phagego/phage-vip4-web/controllers/common/level"
	"phagego/phage-vip4-web/controllers/common/lucky"
	"phagego/phage-vip4-web/controllers/common/membersingle"
	"phagego/phage-vip4-web/controllers/common/membertotal"
	"phagego/phage-vip4-web/controllers/common/mission"
	"phagego/phage-vip4-web/controllers/common/missiondata"
	"phagego/phage-vip4-web/controllers/common/missiondetail"
	"phagego/phage-vip4-web/controllers/common/missionresult"
	"phagego/phage-vip4-web/controllers/common/missionreview"
	"phagego/phage-vip4-web/controllers/common/period"
	"phagego/phage-vip4-web/controllers/common/timegift"
	"phagego/phage-vip4-web/controllers/front"
	"phagego/phage-vip4-web/controllers/rewardlog"
)

func init() {
	// 后台管理系统
	var adminRouter string = beego.AppConfig.String("adminrouter")

	//前端
	beego.Router("/", &front.FrontIndexController{})
	beego.Router("/query", &front.FrontIndexController{}, "post:Query")
	beego.Router("/getgift", &front.FrontIndexController{}, "post:GetGift")
	beego.Router("/vipcenterindex", &front.VipCenterController{})
	beego.Router("/queryprivilege", &front.VipCenterController{}, "get:QueryPrivilege")
	beego.Router("/changetip", &front.VipCenterController{}, "post:ChangeTip")
	beego.Router("/createmissionreview", &front.VipCenterController{}, "post:CreateMissionReview")
	beego.Router("/remark", &front.VipCenterController{}, "post:Remark")
	beego.Router("/api/mission", &front.VipCenterController{}, "post:Mission")
	beego.Router("/gettimegift", &front.VipCenterController{}, "post:GetTimeGift")
	beego.Router("/querytimegift", &front.VipCenterController{}, "post:QueryTimeGift")
	//公告图片上传
	beego.Router(adminRouter+"/upnoticeimg", &imgupload.IndexImguploadController{}, "post:UplodImg")
	//vip等级
	beego.Router(adminRouter+"/level/index", &level.LevelController{})
	beego.Router(adminRouter+"/level/add", &level.LevelAddController{})
	beego.Router(adminRouter+"/level/edit", &level.LevelEditController{})
	beego.Router(adminRouter+"/level/delone", &level.LevelController{}, "post:Delone")
	//好运金
	beego.Router(adminRouter+"/lucky/index", &lucky.LuckyController{})
	beego.Router(adminRouter+"/lucky/add", &lucky.LuckyAddController{})
	beego.Router(adminRouter+"/lucky/edit", &lucky.LuckyEditController{})
	beego.Router(adminRouter+"/lucky/delone", &lucky.LuckyController{}, "post:Delone")
	//周期分类
	beego.Router(adminRouter+"/period/index", &period.PeriodIndexController{})
	beego.Router(adminRouter+"/period/add", &period.PeriodAddController{})
	beego.Router(adminRouter+"/period/edit", &period.PeriodEditController{})
	beego.Router(adminRouter+"/period/delone", &period.PeriodIndexController{}, "post:Delone")
	//单期投注
	beego.Router(adminRouter+"/membersingle/index", &membersingle.MembersingleIndexController{})
	beego.Router(adminRouter+"/membersingle/add", &membersingle.MembersingleAddController{})
	beego.Router(adminRouter+"/membersingle/edit", &membersingle.MembersingleEditController{})
	beego.Router(adminRouter+"/membersingle/delone", &membersingle.MembersingleIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/membersingle/delbatch", &membersingle.MembersingleIndexController{}, "post:DelBatch")
	beego.Router(adminRouter+"/membersingle/import", &membersingle.MembersingleIndexController{}, "post:Import")
	beego.Router(adminRouter+"/membersingle/countgift", &membersingle.MembersingleIndexController{}, "post:CountGift")
	beego.Router(adminRouter+"/membersingle/export", &membersingle.MembersingleIndexController{}, "post:Export")
	//会员统计
	beego.Router(adminRouter+"/membertotal/index", &membertotal.MemberTotalIndexController{})
	beego.Router(adminRouter+"/membertotal/delbatch", &membertotal.MemberTotalIndexController{}, "post:Delbatch")
	beego.Router(adminRouter+"/membertotal/export", &membertotal.MemberTotalIndexController{}, "post:Export")
	beego.Router(adminRouter+"/membertotal/edit", &membertotal.MemberTotalEditController{})
	beego.Router(adminRouter+"/membertotal/count", &membertotal.MemberTotalIndexController{}, "post:Count")
	beego.Router(adminRouter+"/membertotal/import", &membertotal.MemberTotalIndexController{}, "post:Import")
	beego.Router(adminRouter+"/membertotal/changeactivitystatus", &membertotal.MemberTotalIndexController{}, "post:ChangeActivityStatus")
	//中奖记录
	beego.Router(adminRouter+"/rewardlog/index", &rewardlog.RewardLogIndexController{})
	beego.Router(adminRouter+"/rewardlog/delone", &rewardlog.RewardLogIndexController{}, "post:Delone")
	beego.Router(adminRouter+"/rewardlog/Delbatch", &rewardlog.RewardLogIndexController{}, "post:Delbatch")
	beego.Router(adminRouter+"/rewardlog/export", &rewardlog.RewardLogIndexController{}, "post:Export")
	beego.Router(adminRouter+"/rewardlog/delivered", &rewardlog.RewardLogIndexController{}, "post:Delivered")
	beego.Router(adminRouter+"/rewardlog/deliveredbatch", &rewardlog.RewardLogIndexController{}, "post:Deliveredbatch")
	//任务管理
	beego.Router(adminRouter+"/mission/index", &mission.IndexMissionController{})
	beego.Router(adminRouter+"/mission/delone", &mission.IndexMissionController{}, "post:Delone")
	beego.Router(adminRouter+"/mission/Add", &mission.AddMissionController{})
	beego.Router(adminRouter+"/mission/edit", &mission.EditMissionController{})
	//任务详情
	beego.Router(adminRouter+"/missiondetail/index", &missiondetail.IndexMissionDetailController{})
	beego.Router(adminRouter+"/missiondetail/delone", &missiondetail.IndexMissionDetailController{}, "post:Delone")
	beego.Router(adminRouter+"/missiondetail/add", &missiondetail.AddMissionDetailController{})
	beego.Router(adminRouter+"/missiondetail/edit", &missiondetail.EditMissionDetailController{})
	//任务审核
	beego.Router(adminRouter+"/missionreview/index", &missionreview.IndexMissionReview{})
	beego.Router(adminRouter+"/missionreview/delone", &missionreview.IndexMissionReview{}, "post:Delone")
	beego.Router(adminRouter+"/missionreview/review", &missionreview.IndexMissionReview{}, "post:Review")
	beego.Router(adminRouter+"/missionreview/reviewbatch", &missionreview.IndexMissionReview{}, "post:ReviewBatch")
	//任务数据
	beego.Router(adminRouter+"/missiondata/index", &missiondata.IndexMissionDateController{})
	beego.Router(adminRouter+"/missiondata/delone", &missiondata.IndexMissionDateController{}, "post:Delone")
	beego.Router(adminRouter+"/missiondata/import", &missiondata.IndexMissionDateController{}, "post:Import")
	beego.Router(adminRouter+"/missiondata/count", &missiondata.IndexMissionDateController{}, "post:Count")
	beego.Router(adminRouter+"/missiondata/Delbatch", &missiondata.IndexMissionDateController{}, "post:Delbatch")
	//任务计算结果
	beego.Router(adminRouter+"/missionresult/index", &missionresult.IndexMissionResultController{})
	beego.Router(adminRouter+"/missionresult/delone", &missionresult.IndexMissionResultController{}, "post:Delone")
	beego.Router(adminRouter+"/missionresult/review", &missionresult.IndexMissionResultController{}, "post:Review")
	beego.Router(adminRouter+"/missionresult/delbatch", &missionresult.IndexMissionResultController{}, "post:Delbatch")
	beego.Router(adminRouter+"/missionresult/reviewbatch", &missionresult.IndexMissionResultController{}, "post:Reviewbatch")
	//时间奖励
	beego.Router(adminRouter+"/timegift/index", &timegift.IndexTimeGiftController{})
	beego.Router(adminRouter+"/timegift/delone", &timegift.IndexTimeGiftController{}, "post:Delone")
	beego.Router(adminRouter+"/timegift/modifyattr", &timegift.IndexTimeGiftController{}, "post:ModifyAttr")
	beego.Router(adminRouter+"/timegift/modifystatus", &timegift.IndexTimeGiftController{}, "post:ModifyStatus")
	beego.Router(adminRouter+"/timegift/add", &timegift.AddTimeGiftController{})
	beego.Router(adminRouter+"/timegift/edit", &timegift.EditTimeGiftController{})
}
