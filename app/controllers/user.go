package controllers

type UserController struct {
	BaseController
}

// 添加用户
func (this *UserController) Add() {

	this.viewLayoutTitle("新增用户", "user/form", "page")
}

// 用户列表
func (this *UserController) List() {
	this.viewLayoutTitle("用户列表", "user/list", "page")
}

func (this *UserController) Default() {
	this.viewLayoutTitle("首页", "main/default", "page")
}
