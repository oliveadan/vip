package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"html/template"
	"net/http"
	project "phagego/phage-vip4-web/initial"
	_ "phagego/phage-vip4-web/routers"
)

func main() {
	// 初始化项目数据
	project.InitDbProjectData()

	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "pr6FTlKXhAEaYdee5cEeGeJdWgAo7DEn63XWTW5g"
	beego.BConfig.WebConfig.XSRFExpire = 3600
	beego.ErrorHandler("404", page_not_found)
	beego.ErrorHandler("401", page_note_permission)
	beego.SetStaticPath("/upload", "upload")
	beego.Run()
}

func page_not_found(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/404.html")
	data := make(map[string]interface{})
	t.Execute(rw, data)
}

func page_note_permission(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("401.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/401.html")
	data := make(map[string]interface{})
	t.Execute(rw, data)
}
