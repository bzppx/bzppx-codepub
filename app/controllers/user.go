package controllers

import (
	"strings"
	"bzppx-codepub/app/models"
	"log"
	"time"
)

type UserController struct {
	BaseController
}

// 添加用户
func (this *UserController) Add() {
	this.viewLayoutTitle("新增用户", "user/form", "page")
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
		//todo logger
		log.Println(userId)
		this.jsonError("添加用户失败！")
	}else {
		this.jsonSuccess("添加用户成功", nil, "/user/list")
	}
}

// 用户列表
func (this *UserController) List() {
	this.viewLayoutTitle("用户列表", "user/list", "page")
}

func (this *UserController) Default() {
	this.viewLayoutTitle("首页", "main/default", "page")
}
