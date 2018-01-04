package controllers

import (
	"bzppx-codepub/app/models"
	"strings"
)

type ProfileController struct {
	BaseController
}

// 个人资料
func (this *ProfileController) My() {
	this.Data["user"], _ = models.UserModel.GetUserByUserId(this.UserID)
	this.viewLayoutTitle("我的资料", "profile/form", "page")
}

// 修改密码
func (this *ProfileController) Password() {
	this.Data["user"], _ = models.UserModel.GetUserByUserId(this.UserID)
	this.viewLayoutTitle("修改密码", "profile/form-pwd", "page")
}

// 个人资料保存
func (this *ProfileController) Save() {

	givenName := strings.Trim(this.GetString("given_name", ""), "")
	email := strings.Trim(this.GetString("email", "") ,"")
	mobile := strings.Trim(this.GetString("mobile", ""), "")

	if givenName == "" {
		this.jsonError("姓名不能为空！")
	}
	if email == "" {
		this.jsonError("邮箱不能为空！")
	}
	if mobile == "" {
		this.jsonError("手机号不能为空！")
	}

	_, err := models.UserModel.Update(this.Data["loginUser"].(map[string]string)["user_id"], map[string]interface{}{
		"given_name": givenName,
		"email":      email,
		"mobile":     mobile,
	})

	if err != nil {
		this.RecordLog("修改个人资料失败："+err.Error())
		this.jsonError("修改失败")
	} else {
		this.RecordLog("修改个人资料成功")
		this.jsonSuccess("我的资料修改成功", nil, "/profile/my", 3000)
	}
}

// 修改密码保存
func (this *ProfileController) SavePassword() {

	pwd := strings.Trim(this.GetString("pwd", ""), "")
	pwdNew := strings.Trim(this.GetString("pwd_new", ""), "")
	pwdConfirm := strings.Trim(this.GetString("pwd_confirm", ""), "")

	if (pwd == "") || (pwdNew == "") || (pwdConfirm == "") {
		this.jsonError("密码不能为空！")
	}

	p := models.UserModel.EncodePassword(pwd)
	if p != this.User["password"] {
		this.jsonError("当前密码错误")
	}
	if pwdConfirm != pwdNew {
		this.jsonError("确认密码和新密码不一致")
	}

	_, err := models.UserModel.Update(this.Data["loginUser"].(map[string]string)["user_id"], map[string]interface{}{
		"password": models.UserModel.EncodePassword(pwdNew),
	})

	// 阻止日志记录 password
	this.Ctx.Request.PostForm.Del("pwd")
	this.Ctx.Request.PostForm.Del("pwd_new")
	this.Ctx.Request.PostForm.Del("pwd_confirm")

	if err != nil {
		this.RecordLog("修改密码失败："+err.Error())
		this.jsonError("修改密码失败")
	} else {this.RecordLog("修改密码成功")
		this.RecordLog("修改密码成功")
		this.jsonSuccess("修改密码成功", nil, "/profile/my", 3000)
	}
}
