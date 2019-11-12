package common

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type MemberTotal struct {
	Id              int64     `auto`                              // 自增主键
	CreateDate      time.Time `orm:"auto_now_add;type(datetime)"` // 创建时间
	ModifyDate      time.Time `orm:"auto_now;type(datetime)"`     // 更新时间
	Creator         int64     // 创建人Id
	Modifior        int64     // 更新人Id
	Version         int       // 版本
	Account         string    // 会员账号
	LevelUpTime     time.Time `orm:"auto_now;type(datetime)"` // 等级晋升时间
	GetGiftTime     time.Time // 时间奖励领取时间
	Level           int       // vip等级
	Bet             int64     // 总投注额
	MissionIntegral int64     // 任务积分
	TotalLevelGift  int64     // 晋级金总额
	TotalLuckyGift  int64     // 好运金总额
	KeepEnable      int       // 是否保级成功   0:保级成功，1:保级失败
	Tip             int       // 标记是否已提示 0:未提示，1:已提示
	TimeGiftSum     float64   // 时间奖励总和
	ActivityStatus  int       // 活动状态       0:正常 1：异常
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(MemberTotal))
}

func (model *MemberTotal) Create() (int64, error) {
	model.CreateDate = time.Now()
	model.ModifyDate = time.Now()
	model.Version = 0
	o := orm.NewOrm()
	return o.Insert(model)
}

func (model *MemberTotal) Update(cols ...string) (int64, error) {
	if cols != nil {
		cols = append(cols, "ModifyDate", "Modifior")
	}
	model.ModifyDate = time.Now()
	o := orm.NewOrm()
	return o.Update(model, cols...)
}

func (model *MemberTotal) Paginate(page int, limit int, account string, level string, integral string, keep string, order string) (list []MemberTotal, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(MemberTotal))
	cond := orm.NewCondition()
	if account != "" {
		cond = cond.And("Account__contains", account)
	}
	if level != "" {
		cond = cond.And("Level", level)
	}
	if integral != "" {
		cond = cond.And("Bet", integral)
	}
	if keep != "" {
		cond = cond.And("KeepEnable", keep)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	if order != "" {
		qs = qs.OrderBy(order)
	} else {
		qs = qs.OrderBy("-Level", "-Bet")
	}
	qs.All(&list)
	total, _ = qs.Count()
	return
}
