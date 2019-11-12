package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Level struct {
	Id              int64     `auto`                              // 自增主键
	CreateDate      time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate      time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator         int64     // 创建人Id
	Modifior        int64     // 更新人Id
	Version         int       // 版本
	VipLevel        int       // vip等级
	VipName         string    // vip名称
	Bgimg           string    // 背景图片
	Colorimg        string    // 彩色图片
	Wbimg           string    // 黑白图片
	TotalBet        int64     // 累计投注
	LevelGift       int64     // 晋级礼金
	MonthBet        int64     // 每月打码量
	MonthGift       int64     // 每月好运金
	KeepLevelAmount int64     // 保级金额
	KeepLevelDown   int64     // 倒退至等级
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(Level))
}

func (model *Level) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *Level) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *Level) Paginate(page int, limit int) (list []Level, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Level))
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	qs = qs.OrderBy("VipLevel")
	qs.All(&list)
	total, _ = qs.Count()
	return
}
