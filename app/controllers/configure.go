package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"log"
	"strings"
	"time"
)

type ConfigureController struct {
	BaseController
}

//封版配置
func (this *ConfigureController) Block() {

	block, err := models.ConfigureModel.GetBlock()
	if err != nil {
		this.viewError("获取封版配置错误", "configure/block")
	}

	var startTime int64
	var endTime int64
	if block["block_start_time"] == "" || block["block_start_time"] == "0" {
		startTime = time.Now().Unix()
	} else {
		startTime = utils.NewConvert().StringToInt64(block["block_start_time"])
	}
	if block["block_end_time"] == "" || block["block_end_time"] == "0" {
		endTime = time.Now().Unix()
	} else {
		endTime = utils.NewConvert().StringToInt64(block["block_end_time"])
	}
	timePattern := "2006-01-02 15:04"
	block["block_start_time"] = time.Unix(startTime, 0).Format(timePattern)
	block["block_end_time"] = time.Unix(endTime, 0).Format(timePattern)

	this.Data["block"] = block
	this.viewLayoutTitle("封版配置", "configure/block", "page")
}

//添加封版配置
func (this *ConfigureController) AddBlock() {

	blockMessage := strings.Trim(this.GetString("block_message", ""), "")
	blockIsEnable := strings.Trim(this.GetString("block_is_enable", ""), "")
	blockStartTime := strings.Trim(this.GetString("block_start_time", ""), "")
	blockEndTime := strings.Trim(this.GetString("block_end_time", ""), "")

	if blockIsEnable == "0" && blockMessage == "" {
		this.jsonError("封版提示文本不能为空")
	}
	if blockIsEnable == "" {
		this.jsonError("请选择封版开关")
	}
	if blockEndTime == "" {
		this.jsonError("请选择封版开始时间")
	}
	if blockStartTime == "" {
		this.jsonError("请选择封版结束时间")
	}

	if blockIsEnable == "0" {
		timePattern := "2006-01-02 15:04"
		loc, _ := time.LoadLocation("Local")
		startTime, err := time.ParseInLocation(timePattern, blockStartTime, loc)
		if err != nil {
			this.jsonError("开始时间格式错误")
		}
		endTime, err := time.ParseInLocation(timePattern, blockEndTime, loc)
		if err != nil {
			this.jsonError("结束时间格式错误")
		}
		start := startTime.Unix()
		end := endTime.Unix()
		if end < start {
			this.jsonError("开始时间应小于结束时间")
		}
		blockStartTime = utils.NewConvert().IntToString(start, 10)
		blockEndTime = utils.NewConvert().IntToString(end, 10)
	} else {
		blockStartTime = "0"
		blockEndTime = "0"
	}

	blockValue := make([]map[string]interface{}, 4)
	blockValue[0] = map[string]interface{}{
		"key":         "block_message",
		"value":       blockMessage,
		"update_time": time.Now().Unix(),
	}
	blockValue[1] = map[string]interface{}{
		"key":         "block_is_enable",
		"value":       blockIsEnable,
		"update_time": time.Now().Unix(),
	}
	blockValue[2] = map[string]interface{}{
		"key":         "block_start_time",
		"value":       blockStartTime,
		"update_time": time.Now().Unix(),
	}
	blockValue[3] = map[string]interface{}{
		"key":         "block_end_time",
		"value":       blockEndTime,
		"update_time": time.Now().Unix(),
	}

	err := models.ConfigureModel.InsertBlock(blockValue)
	if err != nil {
		log.Println(err.Error())
		this.RecordLog("封版信息修改失败：" + err.Error())
		this.jsonError("封版信息修改失败！")
	} else {
		this.RecordLog("封版信息修改成功")
		this.jsonSuccess("封版信息修改成功", nil, "/configure/block")
	}
}

func (this *ConfigureController) Email() {
	email, err := models.ConfigureModel.GetEmail()
	if err != nil {
		this.viewError("获取邮件配置错误", "configure/email")
	}
	this.Data["email"] = email
	this.viewLayoutTitle("邮件配置", "configure/email", "page")
}

func (this *ConfigureController) AddEmailConfig() {
	emailHost := strings.Trim(this.GetString("email_host", ""), "")
	emailPort := strings.Trim(this.GetString("email_port", ""), "")
	emailUsername := strings.Trim(this.GetString("email_username", ""), "")
	emailPassword := strings.Trim(this.GetString("email_password", ""), "")
	emailFrom := strings.Trim(this.GetString("email_from", ""), "")
	emailIsSsl := strings.Trim(this.GetString("email_is_ssl", ""), "")
	if emailHost == "" {
		this.jsonError("邮箱smtp地址不能为空")
	}
	if emailPort == "" {
		this.jsonError("邮箱smtp端口不能为空")
	}
	if emailUsername == "" {
		this.jsonError("邮箱用户名不能为空")
	}
	if emailPassword == "" {
		this.jsonError("邮箱密码不能为空")
	}
	if emailIsSsl == "" {
		this.jsonError("请选择是否使用ssl")
	}

	blockValue := make([]map[string]interface{}, 6)
	blockValue[0] = map[string]interface{}{
		"key":         "email_host",
		"value":       emailHost,
		"update_time": time.Now().Unix(),
	}
	blockValue[1] = map[string]interface{}{
		"key":         "email_port",
		"value":       emailPort,
		"update_time": time.Now().Unix(),
	}
	blockValue[2] = map[string]interface{}{
		"key":         "email_username",
		"value":       emailUsername,
		"update_time": time.Now().Unix(),
	}
	blockValue[3] = map[string]interface{}{
		"key":         "email_password",
		"value":       emailPassword,
		"update_time": time.Now().Unix(),
	}
	blockValue[4] = map[string]interface{}{
		"key":         "email_from",
		"value":       emailFrom,
		"update_time": time.Now().Unix(),
	}
	blockValue[5] = map[string]interface{}{
		"key":         "email_is_ssl",
		"value":       emailIsSsl,
		"update_time": time.Now().Unix(),
	}

	err := models.ConfigureModel.InsertEmailConfig(blockValue)
	if err != nil {
		log.Println(err.Error())
		this.RecordLog("邮箱信息修改失败：" + err.Error())
		this.jsonError("邮箱信息修改失败！")
	} else {
		this.RecordLog("邮箱信息修改成功")
		this.jsonSuccess("邮箱信息修改成功", nil, "/configure/email")
	}
}

func (this *ConfigureController) AddCcList() {
	emailCcList := strings.Trim(this.GetString("email_cc_list", ""), "")
	if emailCcList == "" {
		this.jsonError("邮件抄送人列表不能为空")
	}
	ccList := make(map[string]interface{})
	ccList["value"] = emailCcList
	err := models.ConfigureModel.InsertCcList(ccList)
	if err != nil {
		log.Println(err.Error())
		this.RecordLog("发版抄送列表修改失败：" + err.Error())
		this.jsonError("发版抄送列表修改失败！")
	} else {
		this.RecordLog("发版抄送列表修改成功")
		this.jsonSuccess("发版抄送列表修改成功", nil, "/configure/email")
	}
}
