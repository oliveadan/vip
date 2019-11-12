package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	frame "phagego/frameweb-v2/initial"
	. "phagego/frameweb-v2/models"
)

// 本项目数据初始化
func InitDbProjectData() {
	// 初始化系统基础数据
	frame.InitDbFrameData()
	// 初始化项目数据
	if beego.AppConfig.DefaultInt("dbautocreate", 0) == 0 {
		return
	}
	beego.Info("Init project data")
	permisions := []Permission{
		//vip设置
		{Id: 100, Pid: 0, Enabled: 1, Display: 1, Description: "vip设置", Url: "", Name: "vip设置", Icon: "#xe614;", Sort: 100},
		{Id: 101, Pid: 100, Enabled: 1, Display: 1, Description: "vip等级", Url: "LevelController.Get", Name: "vip等级", Icon: "", Sort: 100},
		{Id: 102, Pid: 100, Enabled: 1, Display: 0, Description: "添加vip等级", Url: "LevelAddController.Get", Name: "添加vip等级", Icon: "", Sort: 100},
		{Id: 103, Pid: 100, Enabled: 1, Display: 0, Description: "修改vip等级", Url: "LevelEditController.Get", Name: "修改vip等级", Icon: "", Sort: 100},
		{Id: 104, Pid: 100, Enabled: 1, Display: 0, Description: "删除vip等级", Url: "LevelController.Delone", Name: "删除vip等级", Icon: "", Sort: 100},
		//周期分类
		{Id: 110, Pid: 100, Enabled: 1, Display: 1, Description: "周期分类", Url: "PeriodIndexController.get", Name: "周期分类", Icon: "", Sort: 100},
		{Id: 111, Pid: 100, Enabled: 1, Display: 0, Description: "添加周期分类名称", Url: "PeriodAddController.get", Name: "添加周期分类名称", Icon: "", Sort: 100},
		{Id: 112, Pid: 100, Enabled: 1, Display: 0, Description: "修改周期分类", Url: "PeriodEditController.get", Name: "修改周期分类", Icon: "", Sort: 100},
		{Id: 113, Pid: 100, Enabled: 1, Display: 0, Description: "删除周期分类", Url: "PeriodIndexController.Delone", Name: "删除周期分类", Icon: "", Sort: 100},
		//会员投注
		{Id: 120, Pid: 100, Enabled: 1, Display: 1, Description: "会员投注", Url: "MembersingleIndexController.get", Name: "会员投注", Icon: "", Sort: 100},
		{Id: 121, Pid: 100, Enabled: 1, Display: 0, Description: "添加会员投注", Url: "MembersingleAddController.get", Name: "添加会员投注", Icon: "", Sort: 100},
		{Id: 122, Pid: 100, Enabled: 1, Display: 0, Description: "修改会员投注", Url: "MembersingleEditController.get", Name: "修改会员投注", Icon: "", Sort: 100},
		{Id: 123, Pid: 100, Enabled: 1, Display: 0, Description: "删除会员投注", Url: "MembersingleIndexController.Delone", Name: "删除会员投注", Icon: "", Sort: 100},
		{Id: 124, Pid: 100, Enabled: 1, Display: 0, Description: "导入会员投注", Url: "MembersingleIndexController.Import", Name: "导入会员投注", Icon: "", Sort: 100},
		{Id: 125, Pid: 100, Enabled: 1, Display: 0, Description: "删除一期会员投注", Url: "MembersingleIndexController.DelBatch", Name: "删除一期会员投注", Icon: "", Sort: 100},
		{Id: 126, Pid: 100, Enabled: 1, Display: 0, Description: "计算本期彩金", Url: "MembersingleIndexController.CountGift", Name: "计算本期彩金", Icon: "", Sort: 100},
		{Id: 127, Pid: 100, Enabled: 1, Display: 0, Description: "导出会员投注", Url: "MembersingleIndexController.Export", Name: "导出会员投注", Icon: "", Sort: 100},
		//会员统计
		{Id: 130, Pid: 100, Enabled: 1, Display: 1, Description: "会员统计", Url: "MemberTotalIndexController.get", Name: "会员统计", Icon: "", Sort: 100},
		{Id: 131, Pid: 100, Enabled: 1, Display: 0, Description: "删除所有会员统计", Url: "MemberTotalIndexController.Delbatch", Name: "删除所有会员统计", Icon: "", Sort: 100},
		{Id: 132, Pid: 100, Enabled: 1, Display: 0, Description: "导出会员统计", Url: "MemberTotalIndexController.Export", Name: "导出会员统计", Icon: "", Sort: 100},
		{Id: 133, Pid: 100, Enabled: 1, Display: 0, Description: "修改会员统计", Url: "MemberTotalEditController.get", Name: "修改会员统计", Icon: "", Sort: 100},
		{Id: 134, Pid: 100, Enabled: 1, Display: 0, Description: "保级计算", Url: "MemberTotalIndexController.Count", Name: "保级计算", Icon: "", Sort: 100},
		{Id: 135, Pid: 100, Enabled: 1, Display: 0, Description: "导入会员统计", Url: "MemberTotalIndexController.Import", Name: "导入会员统计", Icon: "", Sort: 100},
		{Id: 136, Pid: 100, Enabled: 1, Display: 0, Description: "更改活动状态", Url: "MemberTotalIndexController.ChangeActivityStatus", Name: "更改活动状态", Icon: "", Sort: 100},
		//好运金配置
		{Id: 140, Pid: 100, Enabled: 1, Display: 0, Description: "好运金", Url: "LuckyController.Get", Name: "好运金", Icon: "", Sort: 99},
		{Id: 141, Pid: 100, Enabled: 1, Display: 0, Description: "添加好运金配置", Url: "LuckyAddController.Get", Name: "添加好运金配置", Icon: "", Sort: 99},
		{Id: 142, Pid: 100, Enabled: 1, Display: 0, Description: "修改好运金配置", Url: "LuckyEditController.Get", Name: "修改好运金配置", Icon: "", Sort: 99},
		{Id: 143, Pid: 100, Enabled: 1, Display: 0, Description: "删除好运金配置", Url: "LuckyController.Delone", Name: "删除好运金配置", Icon: "", Sort: 99},
		//中奖记录
		{Id: 149, Pid: 0, Enabled: 1, Display: 1, Description: "中奖记录管理", Url: "", Name: "中奖记录管理", Icon: "#xe614", Sort: 100},
		{Id: 150, Pid: 149, Enabled: 1, Display: 1, Description: "中奖记录", Url: "RewardLogIndexController.get", Name: "中奖记录", Icon: "", Sort: 100},
		{Id: 151, Pid: 149, Enabled: 1, Display: 0, Description: "删除中奖记录", Url: "RewardLogIndexController.Delone", Name: "删除中奖记录", Icon: "", Sort: 100},
		{Id: 152, Pid: 149, Enabled: 1, Display: 0, Description: "删除所有中奖记录", Url: "RewardLogIndexController.Delbatch", Name: "删除所有中奖记录", Icon: "", Sort: 100},
		{Id: 153, Pid: 149, Enabled: 1, Display: 0, Description: "派送中奖记录", Url: "RewardLogIndexController.Delivered", Name: "派送中奖记录", Icon: "", Sort: 100},
		{Id: 154, Pid: 149, Enabled: 1, Display: 0, Description: "批量派送中奖记录", Url: "RewardLogIndexController.Deliveredbatch", Name: "批量派送中奖记录", Icon: "", Sort: 100},
		//任务管理
		{Id: 170, Pid: 0, Enabled: 1, Display: 1, Description: "任务管理", Url: "", Name: "任务管理", Icon: "#xe614;", Sort: 100},
		{Id: 171, Pid: 170, Enabled: 1, Display: 1, Description: "任务配置", Url: "IndexMissionController.get", Name: "任务配置", Icon: "", Sort: 100},
		{Id: 172, Pid: 170, Enabled: 1, Display: 0, Description: "删除任务配置", Url: "IndexMissionController.Delone", Name: "删除任务配置", Icon: "", Sort: 100},
		{Id: 173, Pid: 170, Enabled: 1, Display: 0, Description: "添加任务配置", Url: "AddMissionController.get", Name: "添加任务配置", Icon: "", Sort: 100},
		{Id: 174, Pid: 170, Enabled: 1, Display: 0, Description: "修改任务配置", Url: "EditMissionController.get", Name: "修改任务配置", Icon: "", Sort: 100},
		//任务详情
		{Id: 180, Pid: 170, Enabled: 1, Display: 0, Description: "任务详情", Url: "IndexMissionDetailController.get", Name: "任务详情", Icon: "", Sort: 100},
		{Id: 181, Pid: 170, Enabled: 1, Display: 0, Description: "删除任务详情", Url: "IndexMissionDetailController.Delone", Name: "删除任务详情", Icon: "", Sort: 100},
		{Id: 182, Pid: 170, Enabled: 1, Display: 0, Description: "添加任务详情", Url: "AddMissionDetailController.get", Name: "添加任务详情", Icon: "", Sort: 100},
		{Id: 183, Pid: 170, Enabled: 1, Display: 0, Description: "修改任务详情", Url: "EditMissionDetailController.get", Name: "修改任务详情", Icon: "", Sort: 100},
		//任务审核
		{Id: 190, Pid: 170, Enabled: 1, Display: 1, Description: "任务审核", Url: "IndexMissionReview.get", Name: "任务审核", Icon: ";", Sort: 100},
		{Id: 191, Pid: 170, Enabled: 1, Display: 0, Description: "删除任务审核", Url: "IndexMissionReview.Delone", Name: "删除任务审核", Icon: "", Sort: 100},
		{Id: 192, Pid: 170, Enabled: 1, Display: 0, Description: "改变审核状态", Url: "IndexMissionReview.Review", Name: "改变审核状态", Icon: "", Sort: 100},
		{Id: 193, Pid: 170, Enabled: 1, Display: 0, Description: "批量审核", Url: "IndexMissionReview.ReviewBatch", Name: "批量审核", Icon: "", Sort: 100},
		//任务详情
		{Id: 200, Pid: 170, Enabled: 1, Display: 0, Description: "任务数据", Url: "IndexMissionDateController.get", Name: "任务数据", Icon: "", Sort: 100},
		{Id: 201, Pid: 170, Enabled: 1, Display: 0, Description: "删除任务数据", Url: "IndexMissionDateController.Delone", Name: "删除任务数据", Icon: "", Sort: 100},
		{Id: 202, Pid: 170, Enabled: 1, Display: 0, Description: "导入任务数据", Url: "IndexMissionDateController.Import", Name: "导入任务数据", Icon: "", Sort: 100},
		{Id: 203, Pid: 170, Enabled: 1, Display: 0, Description: "进行计算", Url: "IndexMissionDateController.Count", Name: "进行计算", Icon: "", Sort: 100},
		{Id: 204, Pid: 170, Enabled: 1, Display: 0, Description: "批量删除任务数据", Url: "IndexMissionDateController.Delbatch", Name: "批量删除任务数据", Icon: "", Sort: 100},
		//任务计算结果
		{Id: 210, Pid: 170, Enabled: 1, Display: 0, Description: "任务计算结果", Url: "IndexMissionResultController.get", Name: "任务计算结果", Icon: "", Sort: 100},
		{Id: 211, Pid: 170, Enabled: 1, Display: 0, Description: "删除任务计算结果", Url: "IndexMissionResultController.Delone", Name: "删除任务计算结果", Icon: "", Sort: 100},
		{Id: 212, Pid: 170, Enabled: 1, Display: 0, Description: "审核计算结果", Url: "IndexMissionResultController.Review", Name: "审核计算结果", Icon: "", Sort: 100},
		{Id: 213, Pid: 170, Enabled: 1, Display: 0, Description: "批量删除计算结果", Url: "IndexMissionResultController.Delbatch", Name: "批量删除计算结果", Icon: "", Sort: 100},
		{Id: 214, Pid: 170, Enabled: 1, Display: 0, Description: "批量标记计算结果", Url: "IndexMissionResultController.Reviewbatch", Name: "批量标记计算结果", Icon: "", Sort: 100},
		//时间奖励
		{Id: 220, Pid: 170, Enabled: 1, Display: 1, Description: "时间奖励", Url: "IndexTimeGiftController.get", Name: "时间奖励", Icon: "", Sort: 100},
		{Id: 221, Pid: 170, Enabled: 1, Display: 0, Description: "删除时间奖励", Url: "IndexTimeGiftController.Delone", Name: "删除时间奖励", Icon: "", Sort: 100},
		{Id: 222, Pid: 170, Enabled: 1, Display: 0, Description: "修复时间奖励配置", Url: "IndexTimeGiftController.ModifyAttr", Name: "修复时间奖励配置", Icon: "", Sort: 100},
		{Id: 223, Pid: 170, Enabled: 1, Display: 0, Description: "添加时间奖励", Url: "AddTimeGiftController.get", Name: "添加时间奖励", Icon: "", Sort: 100},
		{Id: 224, Pid: 170, Enabled: 1, Display: 0, Description: "修改时间奖励", Url: "EditTimeGiftController.get", Name: "修改时间奖励", Icon: "", Sort: 100},
		{Id: 225, Pid: 170, Enabled: 1, Display: 0, Description: "修改时间奖励状态", Url: "IndexTimeGiftController.ModifyStatus", Name: "修改时间奖励状态", Icon: "", Sort: 100},
		//公告上传图片
		{Id: 230, Pid: 170, Enabled: 1, Display: 0, Description: "上传公告图片", Url: "IndexImguploadController.UplodImg", Name: "上传公告图片", Icon: "", Sort: 100},
	}
	rolePermissions := []RolePermission{
		{Id: 200, RoleId: 2, PermissionId: 100},
		{Id: 201, RoleId: 2, PermissionId: 101},
		{Id: 202, RoleId: 2, PermissionId: 102},
		{Id: 203, RoleId: 2, PermissionId: 103},
		{Id: 204, RoleId: 2, PermissionId: 104},
		{Id: 205, RoleId: 2, PermissionId: 110},
		{Id: 206, RoleId: 2, PermissionId: 111},
		{Id: 207, RoleId: 2, PermissionId: 112},
		{Id: 208, RoleId: 2, PermissionId: 113},
		{Id: 209, RoleId: 2, PermissionId: 120},
		{Id: 210, RoleId: 2, PermissionId: 121},
		{Id: 211, RoleId: 2, PermissionId: 122},
		{Id: 212, RoleId: 2, PermissionId: 123},
		{Id: 213, RoleId: 2, PermissionId: 124},
		{Id: 214, RoleId: 2, PermissionId: 125},
		{Id: 215, RoleId: 2, PermissionId: 126},
		{Id: 216, RoleId: 2, PermissionId: 127},
		{Id: 220, RoleId: 2, PermissionId: 130},
		{Id: 221, RoleId: 2, PermissionId: 131},
		{Id: 222, RoleId: 2, PermissionId: 132},
		{Id: 223, RoleId: 2, PermissionId: 133},
		{Id: 224, RoleId: 2, PermissionId: 134},
		{Id: 225, RoleId: 2, PermissionId: 135},
		{Id: 226, RoleId: 2, PermissionId: 136},
		{Id: 240, RoleId: 2, PermissionId: 140},
		{Id: 241, RoleId: 2, PermissionId: 141},
		{Id: 252, RoleId: 2, PermissionId: 142},
		{Id: 253, RoleId: 2, PermissionId: 143},
		{Id: 260, RoleId: 2, PermissionId: 149},
		{Id: 261, RoleId: 2, PermissionId: 150},
		{Id: 262, RoleId: 2, PermissionId: 151},
		{Id: 263, RoleId: 2, PermissionId: 152},
		{Id: 264, RoleId: 2, PermissionId: 153},
		{Id: 265, RoleId: 2, PermissionId: 154},
		//客服权限
		{Id: 268, RoleId: 3, PermissionId: 1},
		{Id: 269, RoleId: 3, PermissionId: 3},
		{Id: 270, RoleId: 3, PermissionId: 149},
		{Id: 271, RoleId: 3, PermissionId: 150},
		{Id: 272, RoleId: 3, PermissionId: 153},
		{Id: 273, RoleId: 3, PermissionId: 154},
		//
		{Id: 290, RoleId: 2, PermissionId: 170},
		{Id: 291, RoleId: 2, PermissionId: 171},
		{Id: 292, RoleId: 2, PermissionId: 172},
		{Id: 293, RoleId: 2, PermissionId: 173},
		{Id: 294, RoleId: 2, PermissionId: 174},
		{Id: 300, RoleId: 2, PermissionId: 180},
		{Id: 300, RoleId: 2, PermissionId: 180},
		{Id: 301, RoleId: 2, PermissionId: 181},
		{Id: 302, RoleId: 2, PermissionId: 182},
		{Id: 303, RoleId: 2, PermissionId: 183},
		{Id: 310, RoleId: 2, PermissionId: 190},
		{Id: 311, RoleId: 2, PermissionId: 191},
		{Id: 312, RoleId: 2, PermissionId: 192},
		{Id: 313, RoleId: 2, PermissionId: 193},
		{Id: 320, RoleId: 2, PermissionId: 200},
		{Id: 321, RoleId: 2, PermissionId: 201},
		{Id: 322, RoleId: 2, PermissionId: 202},
		{Id: 323, RoleId: 2, PermissionId: 203},
		{Id: 324, RoleId: 2, PermissionId: 204},
		{Id: 330, RoleId: 2, PermissionId: 210},
		{Id: 331, RoleId: 2, PermissionId: 211},
		{Id: 332, RoleId: 2, PermissionId: 212},
		{Id: 333, RoleId: 2, PermissionId: 213},
		{Id: 334, RoleId: 2, PermissionId: 214},
		//时间奖励
		{Id: 340, RoleId: 2, PermissionId: 220},
		{Id: 341, RoleId: 2, PermissionId: 221},
		{Id: 342, RoleId: 2, PermissionId: 222},
		{Id: 343, RoleId: 2, PermissionId: 223},
		{Id: 344, RoleId: 2, PermissionId: 224},
		{Id: 345, RoleId: 2, PermissionId: 225},
		//上传公告图片
		{Id: 350, RoleId: 2, PermissionId: 230},
	}
	o := orm.NewOrm()
	for _, v := range permisions {
		if _, _, err := o.ReadOrCreate(&v, "Id"); err != nil {
			beego.Error("InitProjectData Permission error", err)
		}
	}
	for _, v := range rolePermissions {
		if _, _, err := o.ReadOrCreate(&v, "Id"); err != nil {
			beego.Error("InitProjectData RolePermission error", err)
		}
	}
}
