package controllers

type ModulesController struct {
	BaseController
}

// 添加模块组
func (this *ModulesController) Add() {
	this.viewLayoutTitle("添加模块组", "modules/form", "page")
}

// 模块组列表
func (this *ModulesController) List() {
	this.viewLayoutTitle("模块组列表", "modules/list", "page")
}

