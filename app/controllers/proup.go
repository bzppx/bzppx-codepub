package controllers

import (
	"strings"
	"bzppx-codepub/app/models"
	"time"
	"bzppx-codepub/app/utils"
)

type GroupController struct {
	BaseController
}

// 添加项目组
func (this *GroupController) Add() {
	this.viewLayoutTitle("添加项目组", "group/form", "page")
}

// 保存项目组
func (this *GroupController) Save() {

	name := strings.Trim(this.GetString("name", ""), "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if name == "" {
		this.jsonError("项目组名称不能为空！")
	}

	group, err := models.GroupModel.HasGroupName(name)
	if err != nil {
		this.ErrorLog("查找项目组 "+name+" 失败："+err.Error())
		this.jsonError("添加项目组失败！")
	}
	if group {
		this.jsonError("项目组名称已存在！")
	}

	groupValue := map[string]interface{}{
		"name": name,
		"comment": comment,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	groupId, err := models.GroupModel.Insert(groupValue)
	if err != nil {
		this.ErrorLog("添加项目组失败: "+err.Error())
		this.jsonError("添加项目组失败！")
	}else {
		this.InfoLog("添加项目组 "+utils.NewConvert().IntToString(groupId, 10)+" 成功")
		this.jsonSuccess("添加项目组成功", nil, "/group/list")
	}
}

// 项目组列表
func (this *GroupController) List() {

	page, _:= this.GetInt("page", 1)
	keyword := strings.Trim(this.GetString("keyword", ""), "")

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var groups []map[string]string
	if (keyword != "") {
		count, err = models.GroupModel.CountGroupsByKeyword(keyword)
		groups, err = models.GroupModel.GetGroupsByKeywordAndLimit(keyword, limit, number)
	}else {
		count, err = models.GroupModel.CountGroups()
		groups, err = models.GroupModel.GetGroupsByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("查找项目组列表失败："+err.Error())
		this.viewError("查找项目组列表失败", "/group/list")
	}

	this.Data["groups"] = groups
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)
	this.viewLayoutTitle("项目组列表", "group/list", "page")
}

// 修改
func (this *GroupController) Edit() {

	groupId := this.GetString("group_id", "")
	if groupId == "" {
		this.viewError("项目组不存在", "/group/list")
	}

	group, err := models.GroupModel.GetGroupByGroupId(groupId)
	if err != nil {
		this.ErrorLog("查找项目组 "+groupId+" 失败："+err.Error())
		this.viewError("项目组不存在", "/group/list")
	}

	this.Data["group"] = group
	this.viewLayoutTitle("修改项目组", "group/form", "page")
}

// 修改保存
func (this *GroupController) Modify() {

	groupId := this.GetString("group_id", "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if groupId == "" {
		this.jsonError("项目组不存在！")
	}

	group, err := models.GroupModel.GetGroupByGroupId(groupId)
	if err != nil {
		this.ErrorLog("查找项目组 "+groupId+" 失败："+err.Error())
		this.jsonError("项目组不存在！")
	}
	if len(group) == 0 {
		this.jsonError("项目组不存在！")
	}

	groupValue := map[string]interface{}{
		"comment": comment,
		"update_time": time.Now().Unix(),
	}

	_, err = models.GroupModel.Update(groupId, groupValue)
	if err != nil {
		this.ErrorLog("修改项目组 "+groupId+" 失败: "+err.Error())
		this.jsonError("修改项目组失败！")
	}else {
		this.InfoLog("修改项目组 "+groupId+" 成功")
		this.jsonSuccess("修改项目组成功", nil, "/group/list")
	}
}

// 删除
func (this *GroupController) Delete() {

	groupId := this.GetString("group_id", "")

	if groupId == "" {
		this.jsonError("没有选择项目组！")
	}

	group, err := models.GroupModel.GetGroupByGroupId(groupId)
	if err != nil {
		this.ErrorLog("查找项目组 "+groupId+" 失败："+err.Error())
		this.jsonError("项目组不存在！")
	}
	if len(group) == 0 {
		this.jsonError("项目组不存在！")
	}

	// todo 判断项目组下的项目是否需要一起删除

	groupValue := map[string]interface{}{
		"is_delete": models.GROUP_DELETE,
		"update_time": time.Now().Unix(),
	}

	_, err = models.GroupModel.Update(groupId, groupValue)
	if err != nil {
		this.ErrorLog("删除项目组 "+groupId+" 失败: "+err.Error())
		this.jsonError("删除项目组失败！")
	}

	this.InfoLog("删除项目组 "+groupId+" 成功")
	this.jsonSuccess("删除项目组成功", nil, "/group/list")
}