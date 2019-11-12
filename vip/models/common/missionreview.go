package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MissionReview struct {
	Id              int64     `auto`                              // 自增主键
	CreateDate      time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate      time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator         int64     // 创建人Id
	Modifior        int64     // 更新人Id
	Version         int       // 版本
	Account         string    // 会员账号
	MissionId       int64     // 任务Id
	MissionDetailId int64     // 活动详情ID
	MinLevel        int64     // 任务最小等级
	MaxLevel        int64     // 任务最大等级
	Integral        int64     // 会员积分
	Remark          string    // 备注
	DeliveredTime   time.Time `orm:"auto_now;type(datetime)"` // 审核时间
	Status          int       // 审核状态 0:审核中，1:通过，2:未通过
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(MissionReview))
}

func (model *MissionReview) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *MissionReview) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	//model.Version =
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *MissionReview) Paginate(page int, limit int, account string, timeStart string, timeEnd string, status int, missionId int64) (list []MissionReview, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MissionReview))
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
	if status != 3 {
		cond = cond.And("Status", status)
	}
	if missionId != 0 {
		cond = cond.And("MissionId", missionId)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-CreateDate")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
