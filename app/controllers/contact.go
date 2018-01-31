package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"regexp"
	"strings"
	"time"
)

type ContactController struct {
	BaseController
}

//列表
func (this *ContactController) List() {
	var err error
	var contacts []map[string]string
	contacts, err = models.ContactModel.GetAllContact()

	if err != nil {
		this.viewError(err.Error(), "contact/list")
	}

	this.Data["contacts"] = contacts
	this.viewLayoutTitle("联系人列表", "contact/list", "page")
}

//添加页面
func (this *ContactController) Add() {
	this.viewLayoutTitle("添加联系人", "contact/form", "page")
}

//修改页面
func (this *ContactController) Update() {
	contactId := strings.Trim(this.GetString("contact_id", ""), "")
	contact, err := models.ContactModel.GetContactByContactId(contactId)
	if err != nil {
		this.viewError(err.Error(), "contact/list")
	}
	if len(contact) == 0 {
		this.viewError("未查到联系人信息", "contact/list")
	}

	this.Data["contact"] = contact
	this.viewLayoutTitle("修改联系人信息", "contact/form", "page")
}

//修改操作
func (this *ContactController) Modify() {
	contactId := strings.Trim(this.GetString("contact_id", ""), "")
	name := strings.Trim(this.GetString("name", ""), "")
	telephone := strings.Trim(this.GetString("telephone", ""), "")
	mobile := strings.Trim(this.GetString("mobile", ""), "")
	position := strings.Trim(this.GetString("position", ""), "")

	if contactId == "" {
		this.jsonError("参数错误！")
	}
	if name == "" {
		this.jsonError("联系人姓名不能为空！")
	}
	if telephone == "" && mobile == "" {
		this.jsonError("手机号和座机号必须有一个！")
	}
	res, err := regexp.MatchString(`^(?:\d{3}-?\d{8}|\d{4}-?\d{7})$`, telephone)
	if err != nil {
		this.jsonError("正则匹配座机电话号失败！")
	}
	if !res && telephone != "" {
		this.jsonError("座机电话号不正确！")
	}
	res, err = regexp.MatchString(`^(((13[0-9]{1})|(15[0-9]{1})|(18[0-9]{1})|(14[0-9]{1})|(17[0-9]{1}))+\d{8})$`, mobile)
	if err != nil {
		this.jsonError("正则匹配手机号失败！")
	}
	if !res && mobile != "" {
		this.jsonError("手机号不正确！")
	}
	if position == "" {
		this.jsonError("联系人职位不能为空！")
	}

	contact := map[string]interface{}{
		"name":        name,
		"telephone":   telephone,
		"mobile":      mobile,
		"position":    position,
		"update_time": time.Now().Unix(),
	}

	_, err = models.ContactModel.UpdateByContactId(contact, contactId)
	if err != nil {
		this.ErrorLog("修改联系人 " + contactId + " 失败: " + err.Error())
		this.jsonError("修改联系人失败！")
	}

	this.InfoLog("修改联系人 " + contactId + " 成功")
	this.jsonSuccess("修改联系人成功", nil, "/contact/list")
}

//添加操作
func (this *ContactController) Save() {
	name := strings.Trim(this.GetString("name", ""), "")
	telephone := strings.Trim(this.GetString("telephone", ""), "")
	mobile := strings.Trim(this.GetString("mobile", ""), "")
	position := strings.Trim(this.GetString("position", ""), "")

	if name == "" {
		this.jsonError("联系人姓名不能为空！")
	}
	if telephone == "" && mobile == "" {
		this.jsonError("手机号和座机号必须有一个！")
	}
	res, err := regexp.MatchString(`^(?:\d{3}-?\d{8}|\d{4}-?\d{7})$`, telephone)
	if err != nil {
		this.jsonError("正则匹配座机电话号失败！")
	}
	if !res && telephone != "" {
		this.jsonError("座机电话号不正确！")
	}
	res, err = regexp.MatchString(`^(((13[0-9]{1})|(15[0-9]{1})|(18[0-9]{1})|(14[0-9]{1})|(17[0-9]{1}))+\d{8})$`, mobile)
	if err != nil {
		this.jsonError("正则匹配手机号失败！")
	}
	if !res && mobile != "" {
		this.jsonError("手机号不正确！")
	}
	if position == "" {
		this.jsonError("联系人职位不能为空！")
	}

	timeNow := time.Now().Unix()
	contact := map[string]interface{}{
		"name":        name,
		"telephone":   telephone,
		"mobile":      mobile,
		"position":    position,
		"create_time": timeNow,
		"update_time": timeNow,
	}
	contactId, err := models.ContactModel.Insert(contact)
	if err != nil {
		this.ErrorLog("添加联系人失败: " + err.Error())
		this.jsonError("添加联系人失败！")
	}

	this.InfoLog("添加联系人 " + utils.NewConvert().IntToString(contactId, 10) + " 成功")
	this.jsonSuccess("添加联系人成功", nil, "/contact/list")
}

//删除操作
func (this *ContactController) Delete() {
	contactId := strings.Trim(this.GetString("contact_id", ""), "")
	_, err := models.ContactModel.DeleteByContactId(contactId)
	if err != nil {
		this.ErrorLog("删除联系人" + contactId + "失败: " + err.Error())
		this.jsonError("删除联系人失败！")
	}

	this.InfoLog("删除联系人 " + contactId + " 成功")
	this.jsonSuccess("删除联系人成功", nil, "/contact/list")
}
