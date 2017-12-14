package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"strings"

	"github.com/astaxie/beego"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) Index() {
	if this.isLogin() {
		this.Redirect("/main/index", 301)
		return
	}
	if this.IsAjax() {

		userModel := models.User{}
		name := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		// captcha := strings.TrimSpace(this.GetString("captcha"))
		// captchaSession := this.GetSession("captcha")
		// if captchaSession == nil || captcha == "" || captchaSession != strings.ToLower(captcha) {
		// 	this.SetSession("captcha", "")
		// 	this.jsonError("验证码错误!")
		// }
		// this.SetSession("captcha", "")
		user, err := userModel.GetUserByName(name)
		if err != nil {
			this.jsonError(err)
			return
		}
		if len(user) == 0 {
			this.jsonError("账号错误!")
		}
		encrypt := new(utils.Encrypt)
		password = userModel.EncodePassword(password)

		if user["password"] != password {
			this.jsonError("账号或密码错误!")
		}
		//加载权限列表

		//保存 session
		this.SetSession("author", user)
		//保存 cookie
		identify := encrypt.Md5Encode(this.Ctx.Request.UserAgent() + this.getClientIp() + password)
		passportValue := encrypt.Base64Encode(name + "@" + identify)
		passport := beego.AppConfig.String("author.passport")
		//fmt.Println("set cookie " + passportValue)
		this.Ctx.SetCookie(passport, passportValue, 3600)

		this.jsonSuccess("登录成功", "", "/main/index.html")
	} else {
		this.viewLayoutTitle("Login", "login/login", "login")
	}
}
