package controllers

import (
	"bzppx-codepub/app/models"
	"strings"
	"time"
	"bzppx-codepub/app/utils"
)

type NoticeController struct {
	BaseController
}

// 添加公告
func (this *NoticeController) Add() {
	this.viewLayoutTitle("添加公告", "notice/form", "page")
}

// 保存公告
func (this *NoticeController) Save() {

	title := strings.Trim(this.GetString("title", ""), "")
	content := strings.Trim(this.GetString("content", ""), "")

	if title == "" {
		this.jsonError("标题不能为空！")
	}
	if content == "" {
		this.jsonError("内容不能为空！")
	}

	noticeValue := map[string]interface{}{
		"title": title,
		"content": content,
		"user_id": this.UserID,
		"username": this.User["user_name"],
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	noticeId, err := models.NoticeModel.Insert(noticeValue)
	if err != nil {
		this.ErrorLog("添加公告插入数据失败: "+err.Error())
		this.jsonError("添加公告失败！")
	}

	this.InfoLog("添加公告 "+utils.NewConvert().IntToString(noticeId, 10)+" 成功")
	this.jsonSuccess("添加公告成功", nil, "/notice/list")
}

// 公告列表
func (this *NoticeController) List() {

	page, _:= this.GetInt("page", 1)
	keyword := strings.Trim(this.GetString("keyword", ""), "")

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var notices []map[string]string
	if (keyword != "") {
		count, err = models.NoticeModel.CountNoticesByKeyword(keyword)
		notices, err = models.NoticeModel.GetNoticesByKeywordAndLimit(keyword, limit, number)
	}else {
		count, err = models.NoticeModel.CountNotices()
		notices, err = models.NoticeModel.GetNoticesByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("公告列表查找公告数据失败: "+err.Error())
		this.viewError("查找公告失败", "/notice/list")
	}

	this.Data["notices"] = notices
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)

	this.viewLayoutTitle("公告列表", "notice/list", "page")
}

// 修改公告
func (this *NoticeController) Edit() {

	noticeId := this.GetString("notice_id", "")
	if noticeId == "" {
		this.viewError("公告不存在", "/notice/list")
	}

	notice, err := models.NoticeModel.GetNoticeByNoticeId(noticeId)
	if err != nil {
		this.viewError("公告不存在", "/notice/list")
	}
	if len(notice) == 0 {
		this.viewError("公告不存在", "/notice/list")
	}

	this.Data["notice"] = notice
	this.viewLayoutTitle("修改公告", "notice/form", "page")
}

// 修改保存公告
func (this *NoticeController) Modify() {
	
	title := strings.Trim(this.GetString("title", ""), "")
	content := strings.Trim(this.GetString("content", ""), "")
	noticeId := strings.Trim(this.GetString("notice_Id", ""), "")

	if title == "" {
		this.jsonError("标题不能为空！")
	}
	if content == "" {
		this.jsonError("内容不能为空！")
	}
	
	noticeValue := map[string]interface{}{
		"title": title,
		"content": content,
		"update_time": time.Now().Unix(),
	}

	_, err := models.NoticeModel.Update(noticeId, noticeValue)
	if err != nil {
		this.ErrorLog("修改公告 "+noticeId+" 失败: "+err.Error())
		this.jsonError("修改公告失败！")
	}else {
		this.InfoLog("修改公告 "+noticeId+" 成功")
		this.jsonSuccess("修改公告成功", nil, "/notice/list")
	}
}

// 公告详细信息
func (this *NoticeController) Info() {
	
	noticeId := this.GetString("notice_id", "")
	if noticeId == "" {
		this.viewError("公告不存在", "/notice/list")
	}
	
	notice, err := models.NoticeModel.GetNoticeByNoticeId(noticeId)
	if err != nil {
		this.ErrorLog("查找公告 "+noticeId+" 失败: "+err.Error())
		this.viewError("公告不存在", "/notice/list")
	}
	if len(notice) == 0 {
		this.viewError("公告不存在", "/notice/list")
	}

	this.Data["notice"] = notice
	this.viewLayoutTitle("公告详细信息", "notice/info", "page")
}

// 删除节点
func (this *NoticeController) Delete() {

	noticeId := this.GetString("notice_id", "")
	
	if noticeId == "" {
		this.jsonError("没有选择公告！")
	}
	
	notice, err := models.NoticeModel.GetNoticeByNoticeId(noticeId)
	if err != nil {
		this.ErrorLog("查找公告 "+noticeId+" 失败: "+err.Error())
		this.jsonError("公告不存在！")
	}
	if len(notice) == 0 {
		this.jsonError("公告不存在！")
	}

	// 删除公告
	noticeValue := map[string]interface{}{
		"is_delete": models.NOTICE_DELETE,
		"update_time": time.Now().Unix(),
	}
	_, err = models.NoticeModel.Update(noticeId, noticeValue)
	if err != nil {
		this.ErrorLog("删除公告 "+noticeId+" 失败: "+err.Error())
		this.jsonError("删除公告失败！")
	}


	this.InfoLog("删除公告 "+noticeId+" 成功")
	this.jsonSuccess("删除公告成功", nil, "/notice/list")
}