package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MissionResult struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	Account    string    // 会员账号
	MissionId  int64     // 关联任务Id
	Prize      string    // 奖品
	Enable     int       // 是否领取 1:已领取，0:未领取
	Status     int       // 审核状态 0:审核中，1:通过，2:未通过
	GetTime    time.Time `orm:"auto_now;type(datetime)"` // 领取时间
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(MissionResult))
}

func (model *MissionResult) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}
func (model *MissionResult) Upate(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model)
}

func (model *MissionResult) Paginate(page int, limit int, id int64, account string, starttime string, endtime string, status int) (list []MissionResult, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MissionResult))
	cond := orm.NewCondition()
	if account != "" {
		cond = cond.And("Account", account)
	}
	if starttime != "" {
		cond = cond.And("CreateDate__gte", starttime)
	}
	if endtime != "" {
		cond = cond.And("CreateDate__lte", endtime)
	}
	if status != 3 {
		cond = cond.And("Status", status)
	}
	cond = cond.And("MissionId", id)
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-CreateDate")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
