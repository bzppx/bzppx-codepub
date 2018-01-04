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

//封板配置
func (this *ConfigureController) Block() {

	block, err := models.ConfigureModel.GetBlock()
	if err != nil {
		this.viewError("获取模块组错误", "configure/block")
	}
	startTime := utils.NewConvert().StringToInt64(block["block_start_time"])
	endTime := utils.NewConvert().StringToInt64(block["block_end_time"])
	timePattern := "2006-01-02 15:04"
	block["block_start_time"] = time.Unix(startTime, 0).Format(timePattern)
	block["block_end_time"] = time.Unix(endTime, 0).Format(timePattern)

	this.Data["block"] = block
	this.viewLayoutTitle("封板配置", "configure/block", "page")
}

//添加封板配置
func (this *ConfigureController) AddBlock() {

	blockMessage := strings.Trim(this.GetString("block_message", ""), "")
	blockIsEnable := strings.Trim(this.GetString("block_is_enable", ""), "")
	blockStartTime := strings.Trim(this.GetString("block_start_time", ""), "")
	blockEndTime := strings.Trim(this.GetString("block_end_time", ""), "")

	if blockMessage == "" {
		this.jsonError("封板提示文本不能为空")
	}
	if blockIsEnable == "" {
		this.jsonError("请选择封板开关")
	}
	if blockEndTime == "" {
		this.jsonError("请选择封板开始时间")
	}
	if blockStartTime == "" {
		this.jsonError("请选择封板结束时间")
	}

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

	blockValue := make([]map[string]interface{}, 4)
	blockValue[0] = map[string]interface{}{
		"key":         "block_message",
		"value":       blockMessage,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	blockValue[1] = map[string]interface{}{
		"key":         "block_is_enable",
		"value":       blockIsEnable,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	blockValue[2] = map[string]interface{}{
		"key":         "block_start_time",
		"value":       blockStartTime,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	blockValue[3] = map[string]interface{}{
		"key":         "block_end_time",
		"value":       blockEndTime,
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}

	err = models.ConfigureModel.InsertBlock(blockValue)
	if err != nil {
		log.Println(err.Error())
		this.RecordLog("封板信息修改失败：" + err.Error())
		this.jsonError("封板信息修改失败！")
	} else {
		this.RecordLog("封板信息修改成功")
		this.jsonSuccess("封板信息修改成功", nil, "/configure/block")
	}
}
