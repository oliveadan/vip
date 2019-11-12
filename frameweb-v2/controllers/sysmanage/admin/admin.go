package admin

import (
	"fmt"
	"html/template"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/frameweb-v2/models"
	. "phagego/frameweb-v2/utils"
	. "phagego/common/utils"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
	"strings"
	"github.com/astaxie/beego/utils/pagination"
)

func validate(admin *Admin) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(admin.Username, "errmsg").Message("用户名必输")
	valid.AlphaDash(admin.Username, "errmsg").Message("用户名必须为字母和数字")
	valid.MaxSize(admin.Username, 30, "errmsg").Message("用户名最长30位")
	valid.Required(admin.Name, "errmsg").Message("名称必输")
	valid.MaxSize(admin.Name, 30, "errmsg").Message("名称最长30位")
	valid.MaxSize(admin.Password, 32, "errmsg").Message("密码不符合规范")
	if admin.Email != "" {
		valid.Email(admin.Email, "errmsg").Message("邮箱格式不正确")
	}
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type AdminIndexController struct {
	sysmanage.BaseController
}

func (c *AdminIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := beego.AppConfig.Int("pagelimit")
	list, total := new(Admin).Paginate(page, limit, c.LoginAdminRefId, param1)
	pagination.SetPaginator(c.Ctx, limit, total)
	// 返回值
	c.Data["dataList"] = &list
	// 查询条件
	c.Data["condArr"] = map[string]interface{}{"param1": param1}
	c.TplName = "sysmanage/admin/index.html"
}

func (c *AdminIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		beego.Error("Delete Admin error", err)
		return
	}
	o := orm.NewOrm()
	// 验证数据权限
	if c.LoginAdminRefId != 0 {
		if exists := o.QueryTable(new(Admin)).Filter("Id", id).Filter("RefId", c.LoginAdminRefId).Exist(); !exists {
			msg = "非法操作"
			return
		}
	}
	o.Begin()
	// 删除管理员角色关联
	if _, err := o.QueryTable(new(AdminRole)).Filter("AdminId", id).Delete(); err != nil {
		o.Rollback()
		beego.Error("Delete admin error 1", err)
		msg = "删除失败"
		return
	}
	if _, err := o.Delete(&Admin{Id: id}); err != nil {
		o.Rollback()
		beego.Error("Delete admin error 2", err)
		msg = "删除失败"
	} else {
		o.Commit()
		code = 1
		msg = "删除成功"
	}
}

func (c *AdminIndexController) LoginVerify() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	t := c.GetString("type")
	verifyCode := c.GetString("code")

	if err != nil {
		msg = "数据错误"
		beego.Error("LoginVerify Admin error", err)
		return
	}
	o := orm.NewOrm()
	model := Admin{Id: id}
	if err := o.Read(&model); err != nil {
		beego.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	if c.LoginAdminRefId != 0 && c.LoginAdminRefId != model.RefId {
		msg = "非法操作"
		return
	}
	if model.LoginVerify == 0 {
		if model.Email == "" {
			msg = "邮箱未配置，请先配置邮箱"
			return
		}
		if t == "send" {
			if err := SendMailVerifyCode(model.Email); err != nil {
				msg = "验证码发送失败"
				return
			} else {
				code = 1
				msg = "验证码发送成功，请查看收件箱或垃圾箱"
				return
			}
		} else {
			if isVerify := VerifyMailVerifyCode(model.Email, verifyCode); !isVerify {
				msg = "验证失败，请重试"
				return
			}
		}
		// 启用登录验证
		model.LoginVerify = 1
	} else {
		model.LoginVerify = 0
	}

	if _, err := o.Update(&model, "LoginVerify"); err != nil {
		beego.Error("Update admin error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
	}
}

func (c *AdminIndexController) Locked() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		beego.Error("Locked Admin error", err)
		return
	}
	o := orm.NewOrm()
	model := Admin{Id: id}
	if err := o.Read(&model); err != nil {
		beego.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	if c.LoginAdminRefId != 0 && c.LoginAdminRefId != model.RefId {
		msg = "非法操作"
		return
	}
	if model.Locked == 1 {
		model.Locked = 0
		model.LoginFailureCount = 0
	} else {
		model.Locked = 1
		model.LockedDate = time.Now()
	}

	if _, err := o.Update(&model, "Locked", "LockedDate"); err != nil {
		beego.Error("Update admin error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
		if model.Locked == 1 { // 如果是锁定，则一并清楚登录token，强制用户退出
			DelCache(fmt.Sprintf("loginAdminId%d", id))
		}
	}
}

type AdminAddController struct {
	sysmanage.BaseController
}

func (c *AdminAddController) Get() {
	refId, _ := c.GetInt64("refId", 0)
	if c.LoginAdminRefId!=0 {
		refId = c.LoginAdminRefId
	}
	c.Data["refId"] = refId
	c.Data["isOrg"] = c.LoginAdminRefId!=0
	c.Data["prefix"] = toLetters(int(c.LoginAdminRefId))
	c.Data["roleList"] = GetRoleList(c.LoginAdminRefId!=0)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "sysmanage/admin/add.html"
}

func (c *AdminAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	admin := Admin{}
	if err := c.ParseForm(&admin); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&admin); hasError {
		msg = errMsg
		return
	} else if admin.Password == "" {
		msg = "密码不能为空"
		return
	} else if admin.Password != c.GetString("repassword") {
		msg = "两次输入的密码不一致"
		return
	}
	roles := c.GetStrings("roles")
	if len(roles) == 0 {
		msg = "请选择所属权限组"
		return
	}
	o := orm.NewOrm()
	if c.LoginAdminRefId != 0 {
		admin.RefId = c.LoginAdminRefId
		if count, err := o.QueryTable(new(Role)).Filter("IsOrg", 1).Filter("Id__in", roles).Count(); err != nil || int(count) != len(roles) {
			msg = "权限获取异常，请刷新后重试"
			return
		}
	}
	if admin.RefId != 0 {
		// 添加用户名前缀
		admin.Username = toLetters(int(admin.RefId)) + "_" + admin.Username
	}
	salt := GetGuid()
	admin.Password = Md5(admin.Password, Pubsalt, salt)
	admin.Salt = salt
	admin.Creator = c.LoginAdminId
	admin.Modifior = c.LoginAdminId
	o.Begin()
	if created, _, err := o.ReadOrCreate(&admin, "Username"); err != nil {
		o.Rollback()
		msg = "添加失败"
		beego.Error("Insert admin error 1", err)
	} else if created {
		adminRoles := make([]AdminRole, 0)
		for _, v := range roles {
			roleId, _ := strconv.ParseInt(v, 10, 64)
			ar := AdminRole{AdminId: admin.Id, RoleId: roleId}
			adminRoles = append(adminRoles, ar)
		}
		if _, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
			o.Rollback()
			msg = "添加失败"
			beego.Error("Insert admin error 3", err)
			return
		}
		o.Commit()
		code = 1
		msg = "添加成功"
	} else {
		msg = "账号已存在"
	}
}

func toLetters(i int) string {
	i += 26
	var j = i/26
	var k = i%26
	var s string
	if j>26 {
		s = toLetters(j)
	} else if j==0 {
		return string(rune(k+97))
	} else {
		return string(rune(j+96))+string(rune(k+97))
	}
	return s+string(rune(k+97))
}

type AdminEditController struct {
	sysmanage.BaseController
}

func (c *AdminEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	admin := Admin{Id: id}

	err := o.Read(&admin)

	if c.LoginAdminRefId != 0 && c.LoginAdminRefId != admin.RefId {
		c.Redirect(beego.URLFor("AdminIndexController.get"), 302)
		return
	}
	arMap := make(map[int64]bool)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(beego.URLFor("AdminIndexController.get"), 302)
		return
	} else if c.LoginAdminRefId == 0 || c.LoginAdminRefId == admin.RefId {
		// 当前管理员所属角色
		var arList orm.ParamsList
		o.QueryTable(new(AdminRole)).Filter("AdminId", id).ValuesFlat(&arList, "RoleId")
		for _, v := range arList {
			arId, ok := v.(int64)
			if ok {
				arMap[arId] = true
			}
		}
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	}
	c.Data["data"] = &admin
	c.Data["adminRoleMap"] = arMap
	c.Data["roleList"] = GetRoleList(c.LoginAdminRefId!=0)
	c.TplName = "sysmanage/admin/edit.html"
}

func (c *AdminEditController) Post() {
	var code int
	var msg string
	var reurl = c.URLFor("AdminIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &reurl)
	admin := Admin{}
	if err := c.ParseForm(&admin); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&admin); hasError {
		msg = errMsg
		return
	} else if admin.Password != "" && admin.Password != c.GetString("repassword") {
		msg = "两次输入的密码不一致"
		return
	}
	roles := c.GetStrings("roles")
	o := orm.NewOrm()
	// 验证数据权限
	if c.LoginAdminRefId != 0 {
		if exists := o.QueryTable(new(Admin)).Filter("Id", admin.Id).Filter("RefId", c.LoginAdminRefId).Exist(); !exists {
			msg = "非法操作"
			return
		}
		if count, err := o.QueryTable(new(Role)).Filter("IsOrg", 1).Filter("Id__in", roles).Count(); err != nil || int(count) != len(roles) {
			msg = "权限获取异常，请刷新后重试"
			return
		}
	}
	cols := []string{"Name", "Enabled", "Email", "ModifyDate"}
	isChangePwd := false
	if admin.Password != "" {
		salt := GetGuid()
		admin.Password = Md5(admin.Password, Pubsalt, salt)
		admin.Salt = salt
		cols = append(cols, "Password", "Salt")
		isChangePwd = true
	}
	admin.Modifior = c.LoginAdminId
	o.Begin()
	if _, err := o.Update(&admin, cols...); err != nil {
		o.Rollback()
		msg = "更新失败"
		beego.Error("Update admin error 1", err)
	} else {
		// 删除旧角色
		if _, err := o.QueryTable(new(AdminRole)).Filter("AdminId", admin.Id).Delete(); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update admin error 2", err)
		}
		// 重新插入角色
		adminRoles := make([]AdminRole, 0)
		for _, v := range roles {
			roleId, _ := strconv.ParseInt(v, 10, 64)
			ar := AdminRole{AdminId: admin.Id, RoleId: roleId}
			adminRoles = append(adminRoles, ar)
		}

		if _, err := o.InsertMulti(len(adminRoles), adminRoles); err != nil {
			o.Rollback()
			msg = "更新失败"
			beego.Error("Update admin error 3", err)
		}
		o.Commit()
		// 如修改了密码，则重置登录，让用户必须重新登录
		if isChangePwd {
			DelCache(fmt.Sprintf("loginAdminId%d", admin.Id))
		}

		code = 1
		msg = "更新成功"
	}
}
