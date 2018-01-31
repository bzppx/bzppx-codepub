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

// 系统配置
func (this *InstallController) Config() {
	this.viewLayoutTitle("安装", "install/config", "install")
}

// 数据库配置
func (this *InstallController) Database() {
	this.viewLayoutTitle("安装", "install/database", "install")
}

// 正在安装
func (this *InstallController) Installing() {
	this.viewLayoutTitle("安装", "install/installing", "install")
}

// 安装完成
func (this *InstallController) Finish() {
	this.viewLayoutTitle("安装", "install/finish", "install")
}