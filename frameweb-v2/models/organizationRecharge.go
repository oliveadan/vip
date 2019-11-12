package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

type OrganizationRecharge struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	RefId      int64                                         // 组织ID
	OrderNo    string    `orm:"unique"`                      // 订单号
	Amount     int                                           // 充值金额
	Status     int                                           // 状态 0：待支付，1：充值成功，2：充值失败
	Remark     string                                        // 备注
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(OrganizationRecharge))
}

func (model *OrganizationRecharge) Paginate(page int, limit int, refId int64) (list []OrganizationRecharge, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(OrganizationRecharge))
	if refId != 0 {
		qs = qs.Filter("RefId", refId)
	}
	qs = qs.OrderBy("-Id")
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs.All(&list)
	total, _ = qs.Count()
	return
}
