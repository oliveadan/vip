package permission

import (
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/frameweb-v2/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func validate(permission *Permission) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(permission.Name, "errmsg").Message("菜单名必输")
	valid.MaxSize(permission.Name, 30, "errmsg").Message("菜单名最长30位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type PermissionIndexController struct {
	sysmanage.BaseController
}

func (this *PermissionIndexController) Get() {
	var permissionList []Permission
	o := orm.NewOrm()
	qs := o.QueryTable(new(Permission))
	qs.All(&permissionList)
	// 返回值
	this.Data["dataList"] = permissionList
	this.TplName = "sysmanage/permission/index.html"
}

func (this *PermissionIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	permission := Permission{Id: id}
	o := orm.NewOrm()
	err := o.Read(&permission)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&permission)
	if err1 != nil {
		beego.Error("Delete permission error", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type PermissionAddController struct {
	sysmanage.BaseController
}

func (this *PermissionAddController) Get() {
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/permission/add.html"
}

func (this *PermissionAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	permission := Permission{}
	if err := this.ParseForm(&permission); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&permission); hasError {
		msg = errMsg
		return
	}
	permission.Creator = this.LoginAdminId
	permission.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&permission); err != nil {
		msg = "添加失败"
		beego.Error("Insert permission error", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type PermissionEditController struct {
	sysmanage.BaseController
}

func (this *PermissionEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	permission := Permission{Id: id}

	err := o.Read(&permission)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("PermissionIndexController.get"), 302)
	} else {
		this.Data["data"] = permission
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "sysmanage/permission/edit.html"
	}
}

func (this *PermissionEditController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	permission := Permission{}
	if err := this.ParseForm(&permission); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&permission); hasError {
		msg = errMsg
		return
	}
	permission.Modifior = this.LoginAdminId
	cols := []string{"Name", "Description", "Enabled", "Pid", "Url", "Icon", "Sort", "Display", "ModifyDate"}
	o := orm.NewOrm()
	if num, err := o.Update(&permission, cols...); err != nil {
		msg = "更新失败"
		beego.Error("Update permission error", err)
	} else if num == 0 {
		msg = "更新失败"
	} else {
		code = 1
		msg = "更新成功"
	}
}
