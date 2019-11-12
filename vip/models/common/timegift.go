package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type TimeGift struct {
	Id          int64     `auto`                              // 自增主键
	CreateDate  time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate  time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator     int64     // 创建人Id
	Modifior    int64     // 更新人Id
	Version     int       // 版本
	GiftName    string    // 礼物名称
	MinMoney    int64     // 最小金额
	MaxMoney    int64     // 最大金额
	GiftContent string    // 礼物内容
	GiftLevel   int64     // 礼物等级
	Category    int       // 类别 0：元 1：角
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(TimeGift))
}

func (model *TimeGift) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *TimeGift) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *TimeGift) Paginate(page int, limit int) (list []TimeGift, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(TimeGift))
	cond := orm.NewCondition()
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
