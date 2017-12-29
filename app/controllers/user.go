package controllers

import (
	"strings"
	"bzppx-codepub/app/models"
	"time"
	"bzppx-codepub/app/utils"
)

type UserController struct {
	BaseController
}

// 添加用户
func (this *UserController) Add() {
	this.viewLayoutTitle("添加用户", "user/form", "page")
}

// 保存用户
func (this *UserController) Save() {

	username := strings.Trim(this.GetString("username", ""), "")
	givenName := strings.Trim(this.GetString("given_name", ""), "")
	password := strings.Trim(this.GetString("password", ""), "")
	email := strings.Trim(this.GetString("email", ""), "")
	mobile := strings.Trim(this.GetString("mobile", ""), "")

	if username == "" {
		this.jsonError("用户名不能为空！")
	}
	if givenName == "" {
		this.jsonError("姓名不能为空！")
	}
	if password == "" {
		this.jsonError("密码不能为空！")
	}
	if email == "" {
		this.jsonError("邮箱不能为空！")
	}
	if mobile == "" {
		this.jsonError("手机号不能为空！")
	}

	user, err := models.UserModel.GetUserByName(username)
	if err != nil {
		this.jsonError("添加用户失败！")
	}
	if len(user) > 0 {
		this.jsonError("该用户名已存在！")
	}

	userValue := map[string]interface{}{
		"username": username,
		"given_name": givenName,
		"password": models.UserModel.EncodePassword(password),
		"email": email,
		"mobile": mobile,
		"role": models.USER_ROLE_USER,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	userId, err := models.UserModel.Insert(userValue)
	if err != nil {
		this.RecordLog("添加用户失败: "+err.Error())
		this.jsonError("添加用户失败！")
	}else {
		this.RecordLog("添加用户 "+utils.NewConvert().IntToString(userId, 10)+" 成功")
		this.jsonSuccess("添加用户成功", nil, "/user/list")
	}
}

// 用户列表
func (this *UserController) List() {

	page, _:= this.GetInt("page", 1)
	keyword := this.GetString("keyword", "")

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var users []map[string]string
	if (keyword != "") {
		count, err = models.UserModel.CountUsersByKeyword(keyword)
		users, err = models.UserModel.GetUsersByKeywordAndLimit(keyword, limit, number)
	}else {
		count, err = models.UserModel.CountUsers()
		users, err = models.UserModel.GetUsersByLimit(limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/user/list")
	}

	this.Data["users"] = users
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)
	this.viewLayoutTitle("用户列表", "user/list", "page")
}

// 修改
func (this *UserController) Edit() {

	userId := this.GetString("user_id", "")
	if userId == "" {
		this.viewError("用户不存在", "/user/list")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.viewError("用户不存在", "/user/list")
	}

	this.Data["user"] = user
	this.viewLayoutTitle("修改用户", "user/form", "page")
}

// 修改保存
func (this *UserController) Modify() {

	userId := strings.Trim(this.GetString("user_id", ""), "")
	givenName := strings.Trim(this.GetString("given_name", ""), "")
	email := strings.Trim(this.GetString("email", ""), "")
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

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.jsonError("用户不存在！")
	}
	if len(user) == 0 {
		this.jsonError("用户不存在！")
	}

	userValue := map[string]interface{}{
		"given_name": givenName,
		"email": email,
		"mobile": mobile,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	_, err = models.UserModel.Update(userId, userValue)
	if err != nil {
		this.RecordLog("修改用户 "+userId+" 失败: "+err.Error())
		this.jsonError("修改用户失败！")
	}else {
		this.RecordLog("修改用户 "+userId+" 成功")
		this.jsonSuccess("修改用户成功", nil, "/user/list")
	}
}

func (this *UserController) Delete() {

	userId := this.GetString("user_id", "")

	if userId == "" {
		this.jsonError("没有选择用户！")
	}

	user, err := models.UserModel.GetUserByUserId(userId)
	if err != nil {
		this.jsonError("用户不存在！")
	}
	if len(user) == 0 {
		this.jsonError("用户不存在！")
	}

	userValue := map[string]interface{}{
		"is_delete": models.USER_DELETE,
		"update_time": time.Now().Unix(),
	}

	_, err = models.UserModel.Update(userId, userValue)
	if err != nil {
		this.RecordLog("删除用户 "+userId+" 失败: "+err.Error())
		this.jsonError("删除用户失败！")
	}

	this.RecordLog("删除用户 "+userId+" 成功")
	this.jsonSuccess("删除用户成功", nil, "/user/list")
}