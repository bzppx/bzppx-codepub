package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"encoding/json"
	"image/color"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
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
		apiAuthId := strings.TrimSpace(this.GetString("api_auth_id"))
		captcha := strings.TrimSpace(this.GetString("captcha"))
		captchaSession := this.GetSession("captcha")
		this.SetSession("captcha", "")
		if captchaSession == nil || captcha == "" || captchaSession != strings.ToLower(captcha) {
			this.jsonError("验证码错误!")
		}

		user := make(map[string]string)
		var err error
		//外部接口登录
		if apiAuthId != "" {
			apiAuth, err := models.ApiAuthModel.GetApiAuthByApiAuthId(apiAuthId)
			if err != nil {
				this.jsonError("获取授权数据失败！")
			}
			request, err := http.PostForm(apiAuth["url"], url.Values{"username": {name}, "password": {password}})
			if err != nil {
				this.jsonError("请求接口失败！")
			}
			defer request.Body.Close()
			responseJson, err := ioutil.ReadAll(request.Body)
			if err != nil {
				this.jsonError("请求接口失败！")
			}
			var response map[string]interface{}
			err = json.Unmarshal(responseJson, &response)
			if err != nil {
				this.jsonError("请求数据错误！")
			}
			if response["msg"].(string) != "" {
				this.jsonError(response["msg"].(string))
			}

			//判断返回的uid的类型，如果不是string转为string
			var uid string
			switch reflect.TypeOf(response["uid"]).String() {
			case "int":
				uid = utils.NewConvert().IntToTenString(response["uid"].(int))
			case "int64":
				uid = utils.NewConvert().IntToString(response["uid"].(int64), 10)
			case "float32":
				uid = utils.NewConvert().FloatToString(response["uid"].(float64), 'f', 0, 32)
			case "float64":
				uid = utils.NewConvert().FloatToString(response["uid"].(float64), 'f', 0, 64)
			case "string":
				uid = response["uid"].(string)
			}
			reg := regexp.MustCompile("^[a-zA-Z0-9_]+$")
			ok := reg.MatchString(uid)
			if !ok {
				this.jsonError("返回数据的uid只能为数字字母和下划线组合错误！")
			}

			userData := make(map[string]interface{})
			userData["given_name"], ok = response["given_name"]
			if !ok {
				this.jsonError("请求数据错误！")
			}
			userData["email"], ok = response["email"]
			if !ok {
				userData["email"] = ""
			}
			userData["mobile"], ok = response["mobile"]
			if !ok {
				userData["mobile"] = ""
			}
			userData["api_auth_id"] = apiAuthId
			userData["username"] = apiAuth["key"] + "_" + uid
			userData["password"] = ""
			userData["last_ip"] = this.getClientIp()
			userData["role"] = "1"
			userData["is_delete"] = "0"
			strTime := utils.NewConvert().IntToString(time.Now().Unix(), 10)
			userData["last_time"] = strTime
			userData["update_time"] = strTime

			user, err = models.UserModel.GetUserByName(userData["username"].(string), apiAuthId)
			if err != nil {
				this.jsonError("获取用户信息错误！")
			}
			//判断user表是否有此账号，没有添加
			if len(user) == 0 {
				userData["create_time"] = strTime
				_, err = models.UserModel.Insert(userData)
				if err != nil {
					this.jsonError("添加用户信息错误！")
				}
				for index, data := range userData {
					user[index] = data.(string)
				}
			} else {
				_, err = models.UserModel.UpdateUserByUsername(userData)
				if err != nil {
					this.jsonError("更新用户信息错误！")
				}
			}

			//转换name
			name = user["username"]
			password = ""
		} else {
			user, err = userModel.GetUserByName(name, "0")
			if err != nil {
				this.jsonError("账号不存在！")
				return
			}
			if len(user) == 0 {
				this.jsonError("账号错误!")
			}

			password = userModel.EncodePassword(password)
			if user["password"] != password {
				this.jsonError("账号或密码错误!")
			}
		}

		//保存 session
		this.SetSession("author", user)
		//保存 cookie
		encrypt := new(utils.Encrypt)
		identify := encrypt.Md5Encode(this.Ctx.Request.UserAgent() + this.getClientIp() + password)
		passportValue := encrypt.Base64Encode(name + "@" + identify)
		passport := beego.AppConfig.String("author.passport")
		//fmt.Println("set cookie " + passportValue)
		this.Ctx.SetCookie(passport, passportValue, 3600)

		userModel.Update(user["user_id"], map[string]interface{}{
			"last_time": time.Now().Unix(),
			"last_ip":   this.getClientIp(),
		})

		this.Ctx.Request.PostForm.Del("password")

		this.InfoLog("登录成功")
		this.jsonSuccess("登录成功", "", "/main/index", 500)
	} else {
		apiLoginList, err := models.ApiAuthModel.GetAllShowApiAuth()
		if err != nil {
			this.viewError("加载API认证信息失败")
		}
		this.Data["apiLoginList"] = apiLoginList
		this.viewLayoutTitle("Login", "login/login", "login")
	}
}

//logout
func (this *LoginController) Logout() {
	this.InfoLog("退出成功")
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
