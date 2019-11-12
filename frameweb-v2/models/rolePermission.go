package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type RolePermission struct {
	Id           int64
	PermissionId int64
	RoleId       int64
}

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysqlpre"), new(RolePermission))
}

func (model *RolePermission) TableUnique() [][]string {
	return [][]string{
		[]string{"PermissionId", "RoleId"},
	}
}
