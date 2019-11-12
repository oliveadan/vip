package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

type Organization struct {
	Id         int64     `auto`                              // 自增主键
	CreateDate time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator    int64                                         // 创建人Id
	Modifior   int64                                         // 更新人Id
	Version    int                                           // 版本
	Name       string                                        // 名称
	BindDomain string                                        // 绑定域名
	ExpireTime time.Time                                     // 过期
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Organization))
}

func (model *Organization) Paginate(page int, limit int) (list []Organization, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Organization))
	qs = qs.OrderBy("-Id")
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs.All(&list)
	total, _ = qs.Count()
	return
}
