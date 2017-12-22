package controllers

import (
	"bzppx-codepub/app/models"
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
	_, err := models.UserModel.Update(this.Data["loginUser"].(map[string]string)["user_id"], map[string]interface{}{
		"given_name": this.GetString("given_name", ""),
		"email":      this.GetString("email", ""),
		"mobile":     this.GetString("mobile", ""),
	})
	if err != nil {
		this.jsonError("修改失败")
	} else {
		this.jsonSuccess("我的资料修改成功", nil, "/profile/my", 3000)
	}
}

// 修改密码保存
func (this *ProfileController) SavePassword() {

	p := models.UserModel.EncodePassword(this.GetString("pwd", ""))
	if p != this.User["password"] {
		this.jsonError("当前密码错误")
	}
	pnew := this.GetString("pwd_new", "")
	if pnew == "" {
		this.jsonError("新密码不能为空")
	}
	_, err := models.UserModel.Update(this.Data["loginUser"].(map[string]string)["user_id"], map[string]interface{}{
		"password": models.UserModel.EncodePassword(pnew),
	})
	if err != nil {
		this.jsonError("修改失败")
	} else {
		this.jsonSuccess("修改成功", nil, "/profile/my", 3000)
	}
}
