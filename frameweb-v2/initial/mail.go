package initial

import (
	"github.com/astaxie/beego"
	"phagego/frameweb-v2/utils"
)

func InitMailConf() {
	sender := beego.AppConfig.String("mailsender")
	host := beego.AppConfig.String("mailhost")
	port, _ := beego.AppConfig.Int("mailport")
	username := beego.AppConfig.String("mailuser")
	password := beego.AppConfig.String("mailpass")

	utils.InitMail(sender, host, port, username, password)
}
