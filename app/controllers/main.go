package controllers

import (
	"bzppx-codepub/app/models"
)

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	this.Data["isAdmin"] = (this.isAdmin() || this.isRoot())
	this.viewLayoutTitle("CodePub POWVEREDBY BZPPX", "main/index", "main")
}

func (this *MainController) Default() {
	var err error
	
	//获取我的项目数
	var projectCount int64
	if (this.isAdmin() || this.isRoot()) {
		projectCount, err = models.ProjectModel.CountProjects()
	}else {
		projectCount, err = models.UserProjectModel.CountProjectByUserId(this.UserID)
	}
	if err != nil {
		this.ErrorLog("获取我的项目数失败：" + err.Error())
		this.viewError("获取数据失败")
	}

	tasks, err := models.TaskModel.GetTasksByUserId(this.UserID)
	if err != nil {
		this.ErrorLog("获取我的发布成功次数失败：" + err.Error())
		this.viewError("获取数据失败")
	}
	taskIds := []string{}
	for _, task := range tasks {
		taskIds = append(taskIds, task["task_id"])
	}
	//获取我的发布失败次数
	failedPublish := 0
	//获取我的发布成功次数
	successPublish := 0
	if len(taskIds) > 0 {
		failedPublish, err = models.TaskLogModel.CountTaskLogByTaskIdsAndIsSuccess(taskIds, models.TASKLOG_FAILED)
		if err != nil {
			this.ErrorLog("获取我的项目数据失败：" + err.Error())
			this.viewError("获取数据失败")
		}
		successPublish = len(taskIds) - failedPublish;
	}

	// 获取项目总排行

	// 获取最新公告
	notices, err := models.NoticeModel.GetNoticesByLimit(0, 6)
	if err != nil {
		this.ErrorLog("获取最新公告失败：" + err.Error())
		this.viewError("获取最新公告失败")
	}

	this.Data["projectCount"] = projectCount
	this.Data["failedPublish"] = failedPublish
	this.Data["successPublish"] = successPublish
	this.Data["notices"] = notices
	this.viewLayoutTitle("首页", "main/default", "index")
}

func (this *MainController) Tpl() {
	typ := this.GetString("type")
	this.viewLayoutTitle("模板", "main/tpl-"+typ, "index")
}
