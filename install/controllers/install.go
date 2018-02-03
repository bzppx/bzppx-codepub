package controllers

import (
	"io/ioutil"
	"runtime"
	"os"
	"strings"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"bzppx-codepub/install/storage"
	"bzppx-codepub/app/utils"
)

type InstallController struct {
	BaseController
}

// 安装首页
func (this *InstallController) Index() {
	this.viewLayoutTitle("codepub-安装", "install/index", "install")
}

// 许可协议
func (this *InstallController) License() {

	if this.isPost() {
		license_agree := this.GetString("license_agree", "")
		if license_agree == "" || license_agree == "0"{
			this.jsonError("请先同意协议后再继续")
		}
		storage.Data.License = storage.License_Agree
		this.jsonSuccess("", nil, "/install/env")
	}else {
		bytes, _ := ioutil.ReadFile("../LICENSE")
		license := string(bytes)
		this.Data["license"] = license

		this.Data["license_agree"] = storage.Data.License
		this.viewLayoutTitle("codepub-安装-License", "install/license", "install")
	}
}

// 环境检测
func (this *InstallController) Env() {

	if this.isPost() {
		if storage.Data.Env == storage.Env_NotAccess {
			this.jsonError("环境检测未通过")
		}
		storage.Data.Env = storage.Env_Access
		this.jsonSuccess("", nil, "/install/config")
	}
	storage.Data.Env = storage.Env_Access
	//获取服务器信息
	host := this.Ctx.Input.Host()
	osSys := runtime.GOOS
	installDir, _ := os.Getwd()
	installDir = strings.Replace(installDir, "install", "", 1)
	server := map[string]string{
		"host": host,
		"sys": osSys,
		"install_dir": installDir,
	}

	// 环境检测
	vm, _ := mem.VirtualMemory()
	vmTotal := vm.Total/1024/1024
	cpuCount, _ := cpu.Counts(true)
	memData := map[string]interface{}{
		"name": "内存",
		"require": "512M",
		"value": strconv.FormatInt(int64(vmTotal), 10)+"M",
		"result": "1",
	}
	if int(vmTotal) < 512 {
		storage.Data.Env = storage.Env_NotAccess
		memData["result"] = "0"
	}
	cpuData := map[string]interface{}{
		"name": "CPU",
		"require": "1核",
		"value": strconv.Itoa(cpuCount)+"核",
		"result": "1",
	}
	if cpuCount < 1 {
		storage.Data.Env = storage.Env_NotAccess
		cpuData["result"] = "0"
	}
	envData := []map[string]interface{}{}
	envData = append(envData, memData)
	envData = append(envData, cpuData)

	// 目录权限检测
	fileTool := utils.NewFile()
	confDir := map[string]string{
		"path": "conf",
		"require": "读/写",
		"result": "1",
	}
	err := fileTool.IsWriterReadable(installDir+confDir["path"]+"/common.conf")
	if err != nil {
		storage.Data.Env = storage.Env_NotAccess
		confDir["result"] = "0"
	}

	docsDir := map[string]string{
		"path": "docs/databases",
		"require": "读/写",
		"result": "1",
	}
	err = fileTool.IsWriterReadable(installDir+docsDir["path"]+"/table.sql")
	if err != nil {
		storage.Data.Env = storage.Env_NotAccess
		docsDir["result"] = "0"
	}
	dirData := []map[string]string{}
	dirData = append(dirData, confDir)
	dirData = append(dirData, docsDir)

	this.Data["server"] = server
	this.Data["envData"] = envData
	this.Data["dirData"] = dirData
	this.viewLayoutTitle("codepub-安装-环境检测", "install/env", "install")
}

// 系统配置
func (this *InstallController) Config() {

	if this.isPost() {
		addr := this.GetString("addr", "")
		port, _ := this.GetInt32("port", 0)

		if addr == "" {
			this.jsonError("addr 不能为空，默认请填写 0.0.0.0")
		}
		if port == 0 {
			this.jsonError("启动端口不能为空")
		}
		if port > int32(65535) {
			this.jsonError("端口超出范围")
		}

		storage.Data.SystemConf = map[string]interface{}{
			"addr": addr,
			"port": port,
		}
		storage.Data.System = storage.Sys_Access
		this.jsonSuccess("", nil, "/install/database")
	}

	sysConf := storage.Data.SystemConf
	this.Data["sysConf"] = sysConf
	this.viewLayoutTitle("codepub-安装-系统配置", "install/config", "install")
}

// 数据库配置
func (this *InstallController) Database() {

	if !this.isPost() {
		this.Data["databaseConf"] = storage.Data.DatabaseConf
		this.viewLayoutTitle("codepub-安装-数据库配置", "install/database", "install")
	}

	host := this.GetString("host", "")
	port := this.GetString("port", "")
	name := this.GetString("name", "")
	user := this.GetString("user", "")
	pass := this.GetString("pass", "")
	connMaxIdle, _:= this.GetInt16("conn_max_idle", 0)
	connMaxConn, _:= this.GetInt16("conn_max_connection", 0)
	adminName := this.GetString("admin_name", "")
	adminPass := this.GetString("admin_pass", "")

	if host == "" {
		this.jsonError("数据库 host 不能为空！")
	}
	if port == "" {
		this.jsonError("数据库端口不能为空！")
	}
	if name == "" {
		this.jsonError("数据库名不能为空！")
	}
	if user == "" {
		this.jsonError("数据库用户名不能为空！")
	}
	if pass == "" {
		this.jsonError("数据库密码不能为空！")
	}
	if connMaxIdle == 0 {
		this.jsonError("数据库连接数不能为0！")
	}
	if connMaxConn == 0 {
		this.jsonError("最大连接数不能为0！")
	}
	if adminName == "" {
		this.jsonError("超级管理员用户名不能为空！")
	}
	if adminPass == "" {
		this.jsonError("超级管理员密码不能为空！")
	}

	storage.Data.DatabaseConf = map[string]interface{}{
		"host": host,
		"port": port,
		"name": name,
		"user": user,
		"pass": pass,
		"conn_max_idle": connMaxIdle,
		"conn_max_connection": connMaxConn,
		"admin_name": adminName,
		"admin_pass": adminPass,
	}
	storage.Data.Database = storage.Database_Access
	this.jsonSuccess("", nil, "/install/ready")
}

// 安装准备
func (this *InstallController) Ready() {

	if this.isPost() {
		if (storage.Data.License != storage.License_Agree) ||
			(storage.Data.Env != storage.Env_Access) ||
			(storage.Data.System != storage.Sys_Access) ||
			(storage.Data.Database != storage.Database_Access) {
			this.jsonError("请先完成安装准备")
		}
		storage.StartInstall()
		this.jsonSuccess("", nil, "/install/end")
	}

	// 协议
	licenseConf := map[string]interface{}{
		"name": "许可协议",
		"value": "同意",
		"result": "1",
		"url": "/install/license",
	}
	if storage.Data.License != storage.License_Agree {
		licenseConf["value"] = "未同意"
		licenseConf["result"] = "0"
	}
	//环境检测
	envConf := map[string]interface{}{
		"name": "环境检测",
		"value": "通过",
		"result": "1",
		"url": "/install/env",
	}
	if storage.Data.Env != storage.Env_Access {
		envConf["value"] = "未通过"
		envConf["result"] = "0"
	}
	//系统配置
	sysConf := map[string]interface{}{
		"name": "系统配置",
		"value": "完成",
		"result": "1",
		"url": "/install/config",
	}
	if storage.Data.System != storage.Sys_Access {
		sysConf["value"] = "未完成"
		sysConf["result"] = "0"
	}
	//数据库配置
	databaseConf := map[string]interface{}{
		"name": "数据库配置",
		"value": "完成",
		"result": "1",
		"url": "/install/database",
	}
	if storage.Data.Database != storage.Database_Access {
		databaseConf["value"] = "未完成"
		databaseConf["result"] = "0"
	}

	readyConf := []map[string]interface{}{}
	readyConf = append(readyConf, licenseConf)
	readyConf = append(readyConf, envConf)
	readyConf = append(readyConf, sysConf)
	readyConf = append(readyConf, databaseConf)

	this.Data["readyConf"] = readyConf
	this.viewLayoutTitle("codepub-安装-准备", "install/ready", "install")
}

// 安装完成
func (this *InstallController) End() {

	if storage.Data.Status == storage.Install_Ready {
		this.Redirect("/install/ready", 302)
	}

	this.viewLayoutTitle("codepub-安装完成", "install/end", "install")
}

// 获取状态
func (this *InstallController) Status() {

	data := map[string]interface{}{
		"status": storage.Data.Status,
		"is_success": storage.Data.IsSuccess,
		"result": storage.Data.Result,
	}

	this.jsonSuccess("ok", data)
}