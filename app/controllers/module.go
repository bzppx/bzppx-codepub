package controllers

type ModuleController struct {
	BaseController
}

// 添加模块
func (this *ModuleController) Add() {

	this.viewLayoutTitle("添加模块", "module/form", "page")
}

// 模块列表
func (this *ModuleController) List() {
	this.viewLayoutTitle("模块列表", "module/list", "page")
}
