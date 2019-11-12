package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type RewardLog struct {
	Id            int64     `auto`                              // 自增主键
	CreateDate    time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate    time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator       int64     // 创建人Id
	Modifior      int64     // 更新人Id
	Version       int       // 版本
	Account       string    // 会员账号
	GiftName      string    // 奖品名称
	GiftContent   string    // 奖品内容
	Delivered     int8      // 是否派奖
	Category      int64     // 奖品类别  1:时间礼品
	DeliveredTime time.Time `orm:"null"` // 派奖时间
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(RewardLog))
}

func (model *RewardLog) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *RewardLog) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	//model.Version =
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *RewardLog) Paginate(page int, limit int, account string, timeStart string, timeEnd string, delivered int8) (list []RewardLog, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(RewardLog))
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
	if delivered != 2 {
		cond = cond.And("Delivered", delivered)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("-Id")
	qs.All(&list)
	total, _ = qs.Count()
	return

}
