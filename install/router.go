package main

import (
	"bzppx-codepub/install/controllers"
	"net/http"
	"bzppx-codepub/app/utils"
	"github.com/astaxie/beego"
)

func init() {

	beego.AppConfig.Set("sys.name", "codepub-installer")
	beego.BConfig.AppName = beego.AppConfig.String("sys.name")
	beego.BConfig.ServerName = beego.AppConfig.String("sys.name")

	// set static path
	beego.SetStaticPath("/static/", "../static")

	// views path
	beego.BConfig.WebConfig.ViewsPath = "../views/"

	// session
	beego.BConfig.WebConfig.Session.SessionName = "codepubssid"
	beego.BConfig.WebConfig.Session.SessionOn = true

	// router
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.RouterCaseSensitive = false

	// todo add router..
	beego.AutoRouter(&controllers.InstallController{})
	beego.Router("/", &controllers.InstallController{}, "*:Index")
	beego.ErrorHandler("404", http_404)
	beego.ErrorHandler("500", http_500)

	// add template func
	beego.AddFuncMap("dateFormat", utils.NewDate().Format)

}

func http_404(rs http.ResponseWriter, req *http.Request) {
	rs.Write([]byte("404 not found!"))
}

func http_500(rs http.ResponseWriter, req *http.Request) {
	rs.Write([]byte("500 server error!"))
}
