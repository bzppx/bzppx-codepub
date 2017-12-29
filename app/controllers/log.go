package controllers

import "bzppx-codepub/app/models"

type LogController struct {
	BaseController
}

// 行为日志列表
func (this *LogController) Action() {
	
	page, _:= this.GetInt("page", 1)
	keyword := this.GetString("keyword", "")
	
	number := 15
	limit := (page - 1) * number
	var err error
	var count int64
	var logActions []map[string]string
	if (keyword != "") {
		count, err = models.LogModel.CountLogsByKeyword(keyword)
		logActions, err = models.LogModel.GetLogsByKeywordAndLimit(keyword, limit, number)
	}else {
		count, err = models.LogModel.CountLogs()
		logActions, err = models.LogModel.GetLogsByLimit(limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/log/action")
	}
	
	this.Data["logActions"] = logActions
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)
	this.viewLayoutTitle("行为日志", "log/action", "page")
}


// 行为日志详细信息
func (this *LogController) ActionInfo() {

	logId := this.GetString("log_id", "")
	if logId == "" {
		this.viewError("日志不存在", "/log/action")
	}

	logAction, err := models.LogModel.GetLogByLogId(logId)
	if err != nil {
		this.viewError("日志不存在", "/log/action")
	}

	this.Data["logAction"] = logAction

	this.viewLayoutTitle("日志详细信息", "log/action-info", "page")
}

// 任务日志列表
func (this *LogController) List() {
	this.viewLayoutTitle("任务日志", "log/task", "page")
}
