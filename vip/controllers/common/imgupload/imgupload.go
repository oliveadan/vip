package imgupload

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"phagego/frameweb-v2/controllers/sysmanage"
	. "phagego/phage-vip4-web/models/common"
)

type IndexImguploadController struct {
	sysmanage.BaseController
}

func (this *IndexImguploadController) Prepare() {
	this.EnableXSRF = false
}

func (this *IndexImguploadController) UplodImg() {
	var code int
	var msg string
	defer sysmanage.Retjson(this.Ctx, &msg, &code)
	imgsrc := this.GetString("imgsrc")
	noticeimg := VipAttribute{}
	noticeimg.Code = "noticeimg"
	noticeimg.Value = imgsrc
	o := orm.NewOrm()
	o.QueryTable(new(VipAttribute)).Filter("Code", "noticeimg").Delete()
	var nc VipAttribute
	nc.Code = "noticeimg"
	nc.Value = imgsrc
	_, err := nc.Create()
	if err != nil {
		beego.Error("create VipAttribute error", err)
		msg = "上传图片失败"
		return
	}
	code = 1
	msg = "上传成功"
}
