package controllers

import (
	"io/ioutil"
)

type InstallController struct {
	BaseController
}

// 安装首页
func (this *InstallController) Index() {
	this.viewLayoutTitle("安装", "install/index", "install")
}

// 许可协议
func (this *InstallController) License() {
	bytes, _ := ioutil.ReadFile("LICENSE")
	license := string(bytes)
	
	this.Data["license"] = license
	this.viewLayoutTitle("安装", "install/license", "install")
}

// 环境检测
func (this *InstallController) Env() {
	this.viewLayoutTitle("安装", "install/env", "install")
}