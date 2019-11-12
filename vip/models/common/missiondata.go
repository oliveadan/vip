package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MissionDate struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	MissionId  int64     // 任务Id
	Account    string    // 会员账号
	Data       int64     // 会员数据
	Period     string    // 期数（导入时间）
	Enable     int       // 是否已计算
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(MissionDate))
}

func (model *MissionDate) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *MissionDate) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model)
}

func (model *MissionDate) Paginate(page int, limit int, id int64, starttime string, endtime string, account string) (list []MissionDate, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MissionDate))
	cond := orm.NewCondition()
	if starttime != "" {
		cond = cond.And("Period__gte", starttime)
	}
	if endtime != "" {
		cond = cond.And("Period__lte", endtime)
	}
	if account != "" {
		cond = cond.And("Account", account)
	}
	cond = cond.And("MissionId", id)
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Period", "Enable", "-Data")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
