package controllers

import (
	"strings"
	"bzppx-codepub/app/models"
	"time"
	"bzppx-codepub/app/utils"
)

type ModulesController struct {
	BaseController
}

// 添加模块组
func (this *ModulesController) Add() {
	this.viewLayoutTitle("添加模块组", "modules/form", "page")
}

// 保存模块组
func (this *ModulesController) Save() {

	name := strings.Trim(this.GetString("name", ""), "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if name == "" {
		this.jsonError("模块组名称不能为空！")
	}

	moduleGroup, err := models.ModulesModel.HasModulesName(name)
	if err != nil {
		this.jsonError("添加模块组失败！")
	}
	if moduleGroup {
		this.jsonError("模块组名称已存在！")
	}

	modulesValue := map[string]interface{}{
		"name": name,
		"comment": comment,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	modulesId, err := models.ModulesModel.Insert(modulesValue)
	if err != nil {
		this.RecordLog("添加模块组失败: "+err.Error())
		this.jsonError("添加模块组失败！")
	}else {
		this.RecordLog("添加模块组 "+utils.NewConvert().IntToString(modulesId, 10)+" 成功")
		this.jsonSuccess("添加模块组成功", nil, "/modules/list")
	}
}

// 模块组列表
func (this *ModulesController) List() {

	page, _:= this.GetInt("page", 1)
	keyword := this.GetString("keyword", "")

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var moduleGroups []map[string]string
	if (keyword != "") {
		count, err = models.ModulesModel.CountModuleGroupsByKeyword(keyword)
		moduleGroups, err = models.ModulesModel.GetModuleGroupsByKeywordAndLimit(keyword, limit, number)
	}else {
		count, err = models.ModulesModel.CountModuleGroups()
		moduleGroups, err = models.ModulesModel.GetModuleGroupsByLimit(limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/modules/list")
	}

	this.Data["moduleGroups"] = moduleGroups
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)
	this.viewLayoutTitle("模块组列表", "modules/list", "page")
}

// 修改
func (this *ModulesController) Edit() {

	modulesId := this.GetString("modules_id", "")
	if modulesId == "" {
		this.viewError("模块组不存在", "/modules/list")
	}

	moduleGroup, err := models.ModulesModel.GetModuleGroupByModulesId(modulesId)
	if err != nil {
		this.viewError("模块组不存在", "/modules/list")
	}

	this.Data["moduleGroup"] = moduleGroup
	this.viewLayoutTitle("修改模块组", "modules/form", "page")
}

// 修改保存
func (this *ModulesController) Modify() {

	modulesId := this.GetString("modules_id", "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if modulesId == "" {
		this.jsonError("模块组不存在！")
	}

	modules, err := models.ModulesModel.GetModuleGroupByModulesId(modulesId)
	if err != nil {
		this.jsonError("模块组不存在！")
	}
	if len(modules) == 0 {
		this.jsonError("模块组不存在！")
	}

	modulesValue := map[string]interface{}{
		"comment": comment,
		"update_time": time.Now().Unix(),
	}

	_, err = models.ModulesModel.Update(modulesId, modulesValue)
	if err != nil {
		this.RecordLog("修改模块组 "+modulesId+" 失败: "+err.Error())
		this.jsonError("修改模块组失败！")
	}else {
		this.RecordLog("修改模块组 "+modulesId+" 成功")
		this.jsonSuccess("修改模块组成功", nil, "/modules/list")
	}
}

// 删除
func (this *ModulesController) Delete() {

	modulesId := this.GetString("modules_id", "")

	if modulesId == "" {
		this.jsonError("没有选择模块组！")
	}

	modules, err := models.ModulesModel.GetModuleGroupByModulesId(modulesId)
	if err != nil {
		this.jsonError("模块组不存在！")
	}
	if len(modules) == 0 {
		this.jsonError("模块组不存在！")
	}

	// todo 判断模块组下的模块是否需要一起删除

	modulesValue := map[string]interface{}{
		"is_delete": models.MODULES_DELETE,
		"update_time": time.Now().Unix(),
	}

	_, err = models.ModulesModel.Update(modulesId, modulesValue)
	if err != nil {
		this.RecordLog("删除模块组 "+modulesId+" 失败: "+err.Error())
		this.jsonError("删除模块组失败！")
	}

	this.RecordLog("删除模块组 "+modulesId+" 成功")
	this.jsonSuccess("删除模块组成功", nil, "/modules/list")
}