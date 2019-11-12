package role

import (
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/frameweb-v2/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

func validate(role *Role) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(role.Name, "errmsg").Message("角色名必输")
	valid.MaxSize(role.Name, 50, "errmsg").Message("角色名最长50位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type RoleIndexController struct {
	sysmanage.BaseController
}

func (this *RoleIndexController) Get() {
	var roleList []Role
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	qs.All(&roleList)
	// 返回值
	this.Data["dataList"] = roleList
	this.TplName = "sysmanage/role/index.html"
}

func (this *RoleIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	id, _ := this.GetInt64("id")
	role := Role{Id: id}
	o := orm.NewOrm()
	err := o.Read(&role)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	} else if role.IsSystem == 1 {
		msg = "系统内置角色，不能删除"
		return
	}
	// 先删除角色权限关联
	o.Begin()
	if _, err := o.QueryTable(new(RolePermission)).Filter("RoleId", role.Id).Delete(); err != nil {
		o.Rollback()
		beego.Error("Delete role error 1", err)
		msg = "删除失败"
		return
	}

	if _, err := o.Delete(&Role{Id: id}); err != nil {
		o.Rollback()
		beego.Error("Delete role error 2", err)
		msg = "删除失败"
		return
	}
	o.Commit()
	code = 1
	msg = "删除成功"
}

type RoleAddController struct {
	sysmanage.BaseController
}

func (this *RoleAddController) Get() {
	this.Data["permissionList"] = GetPermissionList()
	this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "sysmanage/role/add.html"
}

func (this *RoleAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	role := Role{}
	if err := this.ParseForm(&role); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&role); hasError {
		msg = errMsg
		return
	}
	role.Creator = this.LoginAdminId
	role.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&role); err != nil {
		o.Rollback()
		msg = "添加失败"
		beego.Error("Insert role error 1", err)
	} else {
		permissions := this.GetStrings("permissions")
		rolePermissions := make([]RolePermission, 0)
		for _, v := range permissions {
			permissionId, _ := strconv.ParseInt(v, 10, 64)
			ar := RolePermission{RoleId: role.Id, PermissionId: permissionId}
			rolePermissions = append(rolePermissions, ar)
		}
		if _, err := o.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
			o.Rollback()
			msg = "添加失败"
			beego.Error("Insert role error 2", err)
			return
		}
		o.Commit()
		code = 1
		msg = "添加成功"
	}
}

type RoleEditController struct {
	sysmanage.BaseController
}

func (this *RoleEditController) Get() {
	id, _ := this.GetInt64("id")
	o := orm.NewOrm()
	role := Role{Id: id}

	err := o.Read(&role)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		this.Redirect(beego.URLFor("RoleIndexController.get"), 302)
	} else {
		// 当前角色包含的权限
		var rpList orm.ParamsList
		o.QueryTable(new(RolePermission)).Filter("RoleId", id).ValuesFlat(&rpList, "PermissionId")
		rpMap := make(map[int64]bool)
		for _, v := range rpList {
			rpId, ok := v.(int64)
			if ok {
				rpMap[rpId] = true
			}
		}
		this.Data["data"] = role
		this.Data["rolePermissionMap"] = rpMap
		this.Data["permissionList"] = GetPermissionList()
		this.Data["xsrfdata"] = template.HTML(this.XSRFFormHTML())
		this.TplName = "sysmanage/role/edit.html"
	}
}

func (this *RoleEditController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	role := Role{}
	if err := this.ParseForm(&role); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Name", "Description", "Enabled", "ModifyDate"}
	role.Modifior = this.LoginAdminId
	o := orm.NewOrm()
	o.Begin()
	if _, err := o.Update(&role, cols...); err != nil {
		o.Rollback()
		msg = "更新失败"
		beego.Error("Update role error 1", err)
	} else {
		// 删除旧权限
		if _, err := o.QueryTable(new(RolePermission)).Filter("RoleId", role.Id).Delete(); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update role error 2", err)
			return
		}
		// 重新插入新权限
		permissions := this.GetStrings("permissions")
		rolePermissions := make([]RolePermission, 0)
		for _, v := range permissions {
			permissionId, _ := strconv.ParseInt(v, 10, 64)
			ar := RolePermission{RoleId: role.Id, PermissionId: permissionId}
			rolePermissions = append(rolePermissions, ar)
		}
		if _, err := o.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update role error 3", err)
			return
		}
		o.Commit()
		code = 1
		msg = "更新成功"
	}
}
