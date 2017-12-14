package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"image/color"
	"image/png"
	"strings"
	"time"

	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
)

var (
	cap = captcha.New()
)

func init() {
	bs, _ := utils.NewEncrypt().Base64DecodeBytes(fontData)
	cap.AddFontFromBytes(bs)
}

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
		captcha := strings.TrimSpace(this.GetString("captcha"))
		captchaSession := this.GetSession("captcha")
		this.SetSession("captcha", "")
		if captchaSession == nil || captcha == "" || captchaSession != strings.ToLower(captcha) {
			this.jsonError("验证码错误!")
		}
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

		//保存 session
		this.SetSession("author", user)
		//保存 cookie
		identify := encrypt.Md5Encode(this.Ctx.Request.UserAgent() + this.getClientIp() + password)
		passportValue := encrypt.Base64Encode(name + "@" + identify)
		passport := beego.AppConfig.String("author.passport")
		//fmt.Println("set cookie " + passportValue)
		this.Ctx.SetCookie(passport, passportValue, 3600)

		userModel.Update(user["user_id"], map[string]interface{}{
			"last_time": time.Now().Unix(),
			"last_ip":   this.getClientIp(),
		})

		this.jsonSuccess("登录成功", "", "/main/index")
	} else {
		this.viewLayoutTitle("Login", "login/login", "login")
	}
}

//logout
func (this *LoginController) Logout() {
	passport := beego.AppConfig.String("author.passport")
	this.Ctx.SetCookie(passport, "")
	this.SetSession("author", "")
	this.Redirect("/login/index", 302)
	this.StopRun()
}

func (this *LoginController) Captcha() {
	cap.SetSize(80, 28)
	cap.SetDisturbance(captcha.UPPER)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{22, 22, 22, 00})
	img, str := cap.Create(4, captcha.ALL)
	this.SetSession("captcha", strings.ToLower(str))
	png.Encode(this.Ctx.ResponseWriter, img)
}
