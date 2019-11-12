package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Mission struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64     // 创建人Id
	Modifior   int64     // 更新人Id
	Version    int       // 版本
	StartTime  time.Time // 任务开始时间
	EndTime    time.Time // 任务结束时间
	Describe   string    // 任务描述
	CountEnble int       // 是否可计算 1：可计算，0：不可计算
	SumEnable  int       // 是否累计计算：1可以，0:不可以
	Integral   int64     // 任务积分
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Mission))
}

func GetMissions() []Mission {
	o := orm.NewOrm()
	var ms []Mission
	o.QueryTable(new(Mission)).Filter("Id__gt", 0).All(&ms)
	return ms
}

func (model *Mission) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *Mission) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *Mission) Paginate(page int, limit int) (list []Mission, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	cond := orm.NewCondition()
	o := orm.NewOrm()
	qs := o.QueryTable(new(Mission))
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs = qs.SetCond(cond)
	qs.All(&list)
	total, _ = qs.Count()
	return
}
