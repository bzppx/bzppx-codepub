package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	UserID string
	User   map[string]string
}

type JsonResponse struct {
	Code     int                    `json:"code"`
	Message  interface{}            `json:"message"`
	Data     interface{}            `json:"data"`
	Redirect map[string]interface{} `json:"redirect"`
}

// prepare
func (this *BaseController) Prepare() {
	if this.isLogin() && this.inList(beego.AppConfig.String("guest_access_list")) {
		this.Redirect("/main/index", 302)
		this.StopRun()
	}
	if this.isRoot() {
		//root
		if !this.inList(beego.AppConfig.String("root_access_list")) {
			this.viewError("此页面无权访问")
			this.StopRun()
		}
	} else if this.isAdmin() {
		//admin
		if !this.inList(beego.AppConfig.String("admin_access_list")) {
			this.viewError("此页面无权访问")
			this.StopRun()
		}
	} else if this.isLogin() {
		//user
		if !this.inList(beego.AppConfig.String("user_access_list")) {
			this.viewError("此页面无权访问")
			this.StopRun()
		}
	} else {
		//guest
		if this.inList(beego.AppConfig.String("guest_access_list")) {
			return
		}
		this.Redirect("/login/index", 302)
		this.StopRun()
	}
	user := this.GetSession("author").(map[string]string)
	this.User = user
	this.UserID = user["user_id"]
	this.Data["loginUser"] = user
	this.Data["TimeNowYear"] = time.Now().Format("2006")
	this.Layout = "layout/default.html"
}

// check is login
func (this *BaseController) isRoot() bool {
	if !this.isLogin() {
		return false
	}
	user := this.GetSession("author")
	//session 失效
	if user == nil {
		return false
	}
	u := user.(map[string]string)
	return u["role"] == fmt.Sprintf("%s", models.USER_ROLE_ROOT)
}

// check is login
func (this *BaseController) isAdmin() bool {
	if !this.isLogin() {
		return false
	}
	user := this.GetSession("author")
	//session 失效
	if user == nil {
		return false
	}
	u := user.(map[string]string)
	return u["role"] == fmt.Sprintf("%d", models.USER_ROLE_ADMIN) || u["role"] == fmt.Sprintf("%d", models.USER_ROLE_ROOT)
}

// check is login
func (this *BaseController) isLogin() bool {
	passport := beego.AppConfig.String("author.passport")
	cookie := this.Ctx.GetCookie(passport)
	//cookie 失效
	if cookie == "" {
		return false
	}
	user := this.GetSession("author")
	//session 失效
	if user == nil {
		return false
	}
	encrypt := new(utils.Encrypt)
	cookieValue, _ := encrypt.Base64Decode(cookie)
	identifyList := strings.Split(cookieValue, "@")
	if cookieValue == "" || len(identifyList) != 2 {
		return false
	}
	name := identifyList[0]
	identify := identifyList[1]
	userValue := user.(map[string]string)

	//对比cookie 和 session name
	if name != userValue["username"] {
		return false
	}
	//对比客户端UAG and IP
	if identify != utils.NewEncrypt().Md5Encode(this.Ctx.Request.UserAgent()+this.getClientIp()+userValue["password"]) {
		return false
	}
	//success
	return true
}
func (this *BaseController) inList(listString string) bool {
	controllerName, actionName := this.GetControllerAndAction()
	controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	methodName := strings.ToLower(actionName)
	if listString == "*/*" {
		return true
	}
	for _, v := range strings.Split(listString, ";") {
		data := strings.Split(v, "/")
		if len(data) != 2 {
			continue
		}
		c := strings.ToLower(data[0])
		m := strings.ToLower(data[1])
		if c == controllerName {
			if m == "*" {

				return true
			}
			for _, mm := range strings.Split(m, ",") {
				if methodName == mm {
					return true
				}
			}
		}

	}
	return false
}

// view layout title
func (this *BaseController) viewLayoutTitle(title, viewName, layout string) {
	this.Layout = "layout/" + layout + ".html"
	this.TplName = viewName + ".html"
	this.Data["title"] = title
	this.Render()
}

// view layout
func (this *BaseController) viewLayout(viewName, layout string) {
	this.Layout = "layout/" + layout + ".html"
	this.TplName = viewName + ".html"
	this.Data["title"] = ""
	this.Render()
}

// view
func (this *BaseController) view(viewName string) {
	this.Layout = "layout/default.html"
	this.TplName = viewName + ".html"
	this.Data["title"] = ""
	this.Render()
}

// error view
func (this *BaseController) viewError(errorMessage string, data ...interface{}) {
	this.Layout = "layout/default.html"
	errorType := "500"
	if len(data) > 0 {
		errorType = data[0].(string)
	}
	this.TplName = "error/" + errorType + ".html"
	this.Data["title"] = "system error"
	this.Data["errorMessage"] = errorMessage
	this.Render()
}

// view title
func (this *BaseController) viewTitle(title, viewName string) {
	this.Layout = "layout/default.html"
	this.TplName = viewName + ".html"
	this.Data["title"] = title
	this.Render()
}

// return json success
func (this *BaseController) jsonSuccess(message interface{}, data ...interface{}) {
	url := ""
	sleep := 500
	var _data interface{}
	if len(data) > 0 {
		_data = data[0]
	}
	if len(data) > 1 {
		url = data[1].(string)
	}
	if len(data) > 2 {
		sleep = data[2].(int)
	}
	this.Data["json"] = JsonResponse{
		Code:    1,
		Message: message,
		Data:    _data,
		Redirect: map[string]interface{}{
			"url":   url,
			"sleep": sleep,
		},
	}

	j, err := json.MarshalIndent(this.Data["json"], "", "\t")
	if err != nil {
		this.Abort(err.Error())
	} else {
		this.Abort(string(j))
	}
}

// return json error
func (this *BaseController) jsonError(message interface{}, data ...interface{}) {
	url := ""
	sleep := 500
	var _data interface{}
	if len(data) > 0 {
		_data = data[0]
	}
	if len(data) > 1 {
		url = data[1].(string)
	}
	if len(data) > 2 {
		sleep = data[2].(int)
	}
	this.Data["json"] = JsonResponse{
		Code:    0,
		Message: message,
		Data:    _data,
		Redirect: map[string]interface{}{
			"url":   url,
			"sleep": sleep,
		},
	}
	j, err := json.MarshalIndent(this.Data["json"], "", " \t")
	if err != nil {
		this.Abort(err.Error())
	} else {
		this.Abort(string(j))
	}
}

// get client ip
func (this *BaseController) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}
