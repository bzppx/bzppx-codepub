package controllers

import (
	"bzppx-codepub/app/models"
	"strings"
	"time"
	"bzppx-codepub/app/utils"
	"log"
)

type ModuleController struct {
	BaseController
}

// 添加模块
func (this *ModuleController) Add() {

	moduleGroups, err := models.ModulesModel.GetModuleGroups();
	if err != nil {
		this.viewError("获取模块组错误", "module/list")
	}

	this.Data["moduleGroups"] = moduleGroups
	this.viewLayoutTitle("添加模块", "module/form", "page")
}

// 保存模块
func (this *ModuleController) Save() {

	name := strings.Trim(this.GetString("name", ""), "")
	modulesId := strings.Trim(this.GetString("modules_id", ""), "")
	pullType := strings.Trim(this.GetString("pull_type", "http"), "")
	repositoryUrl := strings.Trim(this.GetString("repository_url", ""), "")
	sshKey := strings.Trim(this.GetString("ssh_key", ""), "")
	sshKeySalt := strings.Trim(this.GetString("ssh_key_salt", ""), "")
	httpsUsername := strings.Trim(this.GetString("https_username", ""), "")
	httpsPassword := strings.Trim(this.GetString("https_password", ""), "")
	branch := strings.Trim(this.GetString("branch", ""), "")
	codePath := strings.Trim(this.GetString("code_path", ""), "")
	codeDirUser := strings.Trim(this.GetString("code_dir_user", ""), "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if name == "" {
		this.jsonError("模块名不能为空！")
	}
	if modulesId == "" {
		this.jsonError("没有选择模块组！")
	}
	if repositoryUrl == "" {
		this.jsonError("代码仓库地址不能为空！")
	}
	if pullType == "http" {
		if repositoryUrl[0:7] != "http://" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else if(pullType == "https") {
		if repositoryUrl[0:8] != "https://" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else if(pullType == "ssh") {
		if repositoryUrl[0:4] != "git@" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else {
		this.jsonError("无效的代码拉取方式！")
	}

	if pullType == "http" || pullType == "https" {
		if httpsUsername == "" {
			this.jsonError("用户名不能为空！")
		}
		if httpsPassword == "" {
			this.jsonError("密码不能为空！")
		}
	}
	if pullType == "ssh" {
		if sshKey == "" {
			this.jsonError("ssh key 不能为空！")
		}
	}
	if branch == "" {
		this.jsonError("代码分支不能为空！")
	}
	if codePath == "" {
		this.jsonError("代码发布路径不能为空！")
	}
	if codeDirUser == "" {
		this.jsonError("目录所属用户不能为空！")
	}

	module, err := models.ModuleModel.GetModuleByName(name)
	if err != nil {
		this.RecordLog("添加模块失败: "+err.Error())
		this.jsonError("添加模块失败！")
	}
	if len(module) > 0 {
		this.jsonError("该模块名已存在！")
	}

	moduleValue := map[string]interface{}{
		"name": name,
		"user_id": this.UserID,
		"modules_id": modulesId,
		"repository_url": repositoryUrl,
		"branch": branch,
		"ssh_key": sshKey,
		"ssh_key_salt": sshKeySalt,
		"https_username": httpsUsername,
		"https_password": httpsPassword,
		"code_path": codePath,
		"code_dir_user": codeDirUser,
		"comment": comment,
		"pre_command": "",
		"post_command": "",
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	moduleId, err := models.ModuleModel.Insert(moduleValue)
	if err != nil {
		log.Println(err.Error())
		this.RecordLog("添加模块失败: "+err.Error())
		this.jsonError("添加模块失败！")
	}else {
		this.RecordLog("添加模块 "+utils.NewConvert().IntToString(moduleId, 10)+" 成功")
		this.jsonSuccess("添加模块成功", nil, "/module/list")
	}
}

// 模块列表
func (this *ModuleController) List() {

	page, _:= this.GetInt("page", 1)
	modulesId := this.GetString("modules_id", "")
	keyword := this.GetString("keyword", "")
	keywords := map[string]string {
		"modules_id": modulesId,
		"keyword": keyword,
	}

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var modules []map[string]string
	if (len(keywords) > 0) {
		count, err = models.ModuleModel.CountModulesByKeywords(keywords)
		modules, err = models.ModuleModel.GetModulesByKeywordsAndLimit(keywords, limit, number)
	}else {
		count, err = models.ModuleModel.CountModules()
		modules, err = models.ModuleModel.GetModulesByLimit(limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/module/list")
	}

	moduleGroups, err := models.ModulesModel.GetModuleGroups();
	if err != nil {
		this.viewError("获取模块组错误", "module/list")
	}

	this.Data["modules"] = modules
	this.Data["keywords"] = keywords
	this.Data["moduleGroups"] = moduleGroups
	this.SetPaginator(number, count)

	this.viewLayoutTitle("模块列表", "module/list", "page")
}

// 修改模块
func (this *ModuleController) Edit() {

	moduleId := this.GetString("module_id", "")
	if moduleId == "" {
		this.viewError("模块不存在", "/module/list")
	}

	module, err := models.ModuleModel.GetModuleByModuleId(moduleId)
	if err != nil {
		this.viewError("模块不存在", "/module/list")
	}

	repUrl := module["repository_url"]
	pullType := "http"
	if repUrl[0:4] == "ssh@" {
		pullType = "ssh"
	}else if repUrl[0:8] == "https://" {
		pullType = "https"
	}else {
		pullType = "http"
	}

	moduleGroups, err := models.ModulesModel.GetModuleGroups();
	if err != nil {
		this.viewError("获取模块组错误", "/module/list")
	}

	this.Data["module"] = module
	this.Data["pullType"] = pullType
	this.Data["moduleGroups"] = moduleGroups
	this.viewLayoutTitle("修改模块", "module/form", "page")
}