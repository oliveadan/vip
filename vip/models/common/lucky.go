package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Lucky struct {
	Id          int64     `auto`                              // 自增主键
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator     int64     // 创建人Id
	Modifior    int64     // 更新人Id
	Version     int       // 版本
	MinVipLevel int       // 最小vip等级
	MaxVipLevel int       // 最大vip等级
	MonthBet    int64     // 天投注
	LuckyGift   int64     // 好运金
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Lucky))
}

func (model *Lucky) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *Lucky) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *Lucky) Paginate(page int, limit int) (list []Lucky, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lucky))
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("MinVipLevel")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
