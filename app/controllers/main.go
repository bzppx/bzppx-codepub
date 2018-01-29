package controllers

import (
	"bzppx-codepub/app/models"
	"encoding/json"
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
	var groupCount int64
	if (this.isAdmin() || this.isRoot()) {
		projectCount, err = models.ProjectModel.CountProjects()
		groupCount, err = models.GroupModel.CountGroups()
	}else {
		userProjects, err := models.UserProjectModel.GetUserProjectByUserId(this.UserID)
		if err != nil {
			this.ErrorLog("获取我的项目组数据失败: " + err.Error())
			this.viewError("获取数据失败")
		}
		projectIds := []string{}
		for _, userProject := range userProjects {
			projectIds = append(projectIds, userProject["project_id"])
		}
		projectCount = int64(len(projectIds))
		groupCount, err = models.ProjectModel.CountGroupByProjectIds(projectIds)
		if err != nil {
			this.ErrorLog("获取我的项目组数据失败: " + err.Error())
			this.viewError("获取数据失败")
		}
	}
	if err != nil {
		this.ErrorLog("获取我的项目组数据失败：" + err.Error())
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
	projectPublishCountRank := []map[string]string{}
	projectCountIds, err := models.TaskModel.GetProjectIdsOrderByCountProject()
	if err != nil {
		this.ErrorLog("获取项目总排行数据失败：" + err.Error())
		this.viewError("获取数据失败")
	}
	if len(projectCountIds) > 0 {
		projectIds := []string{}
		for _, projectCountId := range projectCountIds {
			projectIds = append(projectIds, projectCountId["project_id"])
		}
		projects, err := models.ProjectModel.GetProjectByProjectIds(projectIds)
		if err != nil {
			this.ErrorLog("获取项目总排行数据失败：" + err.Error())
			this.viewError("获取数据失败")
		}
		for _, projectCountId := range projectCountIds {
			projectPublishCount := map[string]string{
				"project_name": "",
				"total": projectCountId["total"],
			}
			for _, project := range projects {
				if projectCountId["project_id"] == project["project_id"] {
					projectPublishCount["project_name"] = project["name"]
					break
				}
			}
			projectPublishCountRank = append(projectPublishCountRank, projectPublishCount)
		}
	}

	// 获取最新公告
	notices, err := models.NoticeModel.GetNoticesByLimit(0, 6)
	if err != nil {
		this.ErrorLog("获取最新公告失败：" + err.Error())
		this.viewError("获取最新公告失败")
	}

	jsonProjectPublishCountRank, _ := json.Marshal(projectPublishCountRank)

	this.Data["groupCount"] = groupCount
	this.Data["projectCount"] = projectCount
	this.Data["failedPublish"] = failedPublish
	this.Data["successPublish"] = successPublish
	this.Data["projectPublishCountRank"] = string(jsonProjectPublishCountRank)
	this.Data["notices"] = notices
	this.viewLayoutTitle("首页", "main/default", "index")
}

func (this *MainController) Tpl() {
	typ := this.GetString("type")
	this.viewLayoutTitle("模板", "main/tpl-"+typ, "index")
}
