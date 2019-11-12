package utils

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"phagego/phage-vip4-web/models/common"
	"strconv"
)

func GetCountEnAble() map[int]string {
	m := map[int]string{
		0: "否",
		1: "是",
	}
	return m
}

func GetSumTimeGift(account string) string {
	o := orm.NewOrm()
	var rl []common.RewardLog
	_, _ = o.QueryTable(new(common.RewardLog)).Filter("Account", account).Filter("Category", 1).All(&rl)
	var sum float64
	for _, v := range rl {
		s, _ := strconv.ParseFloat(v.GiftContent, 32)
		sum += s
	}
	return fmt.Sprintf("%1.2f", sum)
}

func GetMissionDescribe(id int64) string {
	o := orm.NewOrm()
	var mission common.Mission
	_ = o.QueryTable(new(common.Mission)).Filter("Id", id).One(&mission)
	return mission.Describe
}

func GetMissionDetail(id int64) string {
	o := orm.NewOrm()
	var missionDetail common.MissionDetail
	_ = o.QueryTable(new(common.MissionDetail)).Filter("Id", id).One(&missionDetail)
	return missionDetail.Award
}
