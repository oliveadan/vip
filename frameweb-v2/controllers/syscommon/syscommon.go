package syscommon

import (
	"fmt"
	"net/url"
	"os"
	"phagego/common/utils"
	fu "phagego/frameweb-v2/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type SyscommonController struct {
	beego.Controller
}

// 注意，本接口 code = 0 时才是上传成功
func (c *SyscommonController) Upload() {
	var code int
	var msg string
	var uploadName string
	f, h, err := c.GetFile("file")
	defer f.Close()
	if err != nil {
		beego.Error("Syscommon upload file get file error", err)
		code = 1
		msg = "上传失败，请重试(1)"
	} else {
		fname := url.QueryEscape(h.Filename)
		suffix := utils.SubString(fname, len(fname), strings.LastIndex(fname, ".")-len(fname))
		uploadPath := fmt.Sprintf("upload/%d/%s/%d/", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
		if flag, _ := utils.PathExists(uploadPath); !flag {
			if err2 := os.MkdirAll(uploadPath, 0644); err2 != nil {
				beego.Error("Syscommon upload file get file error", err2)
				code = 2
				msg = "上传失败，请重试(2)"
			}
		}

		if code == 0 {
			uploadName = uploadPath + strconv.FormatInt(time.Now().UnixNano(), 10) + suffix
			err3 := c.SaveToFile("file", uploadName)
			if err3 != nil {
				beego.Error("Syscommon upload file save file error2", err3)
				code = 3
				msg = "上传失败，请重试(3)"
			} else {
				msg = "上传成功"
			}
		}
	}
	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	ret["data"] = map[string]string{"src": "/" + uploadName}
	c.Data["json"] = ret
	c.ServeJSON()
}

// 暂不要使用
// 发送验证码 和 校验验证码
// 发送验证码 参数：type=send, to=邮箱
// 校验验证码 参数：type=verify, to=邮箱, code=验证码
func (c *SyscommonController) MailVerify() {
	var code int
	var msg string
	t := c.GetString("type")
	to := c.GetString("to")
	verifyCode := c.GetString("code")
	if to == "" {
		msg = "收件邮箱不能为空"
	} else {
		if t == "send" {
			vc := utils.RandStringLower(4)
			ms := fu.MailSender{To: []string{to},
				Subject: "Phage系统验证码",
				Body: vc}

			err := ms.Send()
			if err != nil {
				msg = "验证码发送失败"
			} else {
				c.SetSession("mailverifycode"+to, vc)
				code = 1
			}
		} else if t == "verify" {
			vc := c.GetSession("mailverifycode"+to)
			if vc == nil {
				msg = "验证失败"
			} else if vc.(string) == verifyCode {
				code = 1
				msg = "验证成功"
			} else {
				msg = "验证失败"
			}
		}
	}

	ret := make(map[string]interface{})
	ret["code"] = code
	ret["msg"] = msg
	c.Data["json"] = ret
	c.ServeJSON()
}

func (c *SyscommonController) HealthCheck() {
	c.Ctx.ResponseWriter.Write([]byte("1"))
}
