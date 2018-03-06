package controllers

import (
	"bzppx-codepub/app/models"
	"strings"
)

type LogController struct {
	BaseController
}

// 行为日志列表
func (this *LogController) Action() {

	page, _ := this.GetInt("page", 1)
	level := strings.Trim(this.GetString("level", ""), "")
	message := strings.Trim(this.GetString("message", ""), "")
	username := strings.Trim(this.GetString("username", ""), "")

	number := 15
	limit := (page - 1) * number
	var err error
	var count int64
	var logActions []map[string]string
	if level != "" || message != "" || username != "" {
		count, err = models.LogModel.CountLogsByKeyword(level, message, username)
		logActions, err = models.LogModel.GetLogsByKeywordAndLimit(level, message, username, limit, number)
	} else {
		count, err = models.LogModel.CountLogs()
		logActions, err = models.LogModel.GetLogsByLimit(limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/log/action")
	}

	this.Data["logActions"] = logActions
	this.Data["username"] = username
	this.Data["level"] = level
	this.Data["message"] = message
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
func (this *LogController) Task() {
	page, _ := this.GetInt("page", 1)
	projectName := strings.Trim(this.GetString("project_name", ""), "")
	userName := strings.Trim(this.GetString("user_name", ""), "")

	number := 15
	limit := (page - 1) * number
	var err error
	var count int64
	projectIds := []string{}
	projectData := make(map[string]map[string]string)
	if projectName != "" {
		projects, err := models.ProjectModel.GetProjectsByLikeName(projectName)
		if err != nil {
			this.viewError(err.Error(), "/log/task")
		}

		for _, project := range projects {
			projectIds = append(projectIds, project["project_id"])
			projectData[project["project_id"]] = project
		}
	}

	userIds := []string{}
	userData := make(map[string]map[string]string)
	if userName != "" {
		users, err := models.UserModel.GetUserByLikeName(userName)
		if err != nil {
			this.viewError(err.Error(), "/log/task")
		}

		for _, user := range users {
			userIds = append(userIds, user["user_id"])
			userData[user["user_id"]] = user
		}
	}
	var tasks []map[string]string
	if userName == "" && projectName == "" {
		count, err = models.TaskModel.CountTask()
		tasks, err = models.TaskModel.GetTasksByLimit(limit, number)
	} else {
		count, err = models.TaskModel.CountTaskByUserIdsAndProjectIds(userName, projectName, userIds, projectIds)
		tasks, err = models.TaskModel.GetTasksByUserIdsAndProjectIdsAndLimit(userName, projectName, userIds, projectIds, limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/log/task")
	}

	taskIds := make([]string, len(tasks))
	userIds = make([]string, len(tasks))
	projectIds = make([]string, len(tasks))
	for index, task := range tasks {
		taskIds[index] = task["task_id"]
		userIds[index] = task["user_id"]
		projectIds[index] = task["project_id"]
	}

	if len(userData) <= 0 {
		users, err := models.UserModel.GetUserByUserIds(userIds)
		if err != nil {
			this.viewError(err.Error(), "/log/task")
		}
		for _, user := range users {
			userData[user["user_id"]] = user
		}
	}
	if len(projectData) <= 0 {
		projects, err := models.ProjectModel.GetProjectByProjectIds(projectIds)
		if err != nil {
			this.viewError(err.Error(), "/log/task")
		}
		for _, project := range projects {
			projectData[project["project_id"]] = project
		}
	}

	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskIds(taskIds)
	if err != nil {
		this.viewError(err.Error(), "/log/task")
	}
	
	tasklogsChangeKey := make(map[string]map[int]map[string]string)
	for index, taskLog := range taskLogs {
		_, ok := tasklogsChangeKey[taskLog["task_id"]]
		if ok {
			tasklogsChangeKey[taskLog["task_id"]][index] = taskLog
		} else {
			tasklogsChangeKey[taskLog["task_id"]] = make(map[int]map[string]string)
			tasklogsChangeKey[taskLog["task_id"]][index] = taskLog
		}
	}
	for index, task := range tasks {
		tasks[index]["status"] = "1"
		tasks[index]["project_name"] = projectData[task["project_id"]]["name"]
		tasks[index]["username"] = userData[task["user_id"]]["username"]
		for _, taskLogChangeKey := range tasklogsChangeKey[task["task_id"]] {
			if taskLogChangeKey["status"] != "2" {
				tasks[index]["status"] = "2"
				break
			}
			if taskLogChangeKey["is_success"] == "0" {
				tasks[index]["status"] = "0"
			}
		}
	}

	this.Data["tasks"] = tasks
	this.Data["user_name"] = userName
	this.Data["project_name"] = projectName
	this.SetPaginator(number, count)
	this.viewLayoutTitle("任务日志", "log/task", "page")
}
