package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"strings"
	"time"
)

type TaskController struct {
	BaseController
}

func (this *TaskController) Center() {
	this.viewLayoutTitle("任务队列", "task/center", "page")
}

func (this *TaskController) Node() {
	this.viewLayoutTitle("节点进度", "task/node", "page")
}

func (this *TaskController) GetExcutingTask() {
	//页面传的数据
	viewProjectIds := this.GetString("project_id", "")
	viewTaskIds := this.GetString("task_id", "")
	var projectIds []string
	taskIds, err := models.TaskLogModel.GetExcutingTaskIdByTaskLog()
	if err != nil {
		this.jsonError("获取taskId失败！")
	}
	if viewTaskIds != "" {
		taskIds = append(taskIds, strings.Split(viewTaskIds, ",")...)
	}

	if !this.isAdmin() && !this.isRoot() {
		projectUsers, err := models.UserProjectModel.GetUserProjectByUserId(this.UserID)
		if err != nil {
			this.jsonError("获取权限失败！")
		}
		if len(projectUsers) == 0 {
			this.jsonError("此账号没有被授予任何项目权限！")
		}
		projectIds = utils.NewArray().ArrayColumn(projectUsers, "project_id")
	}
	if len(projectIds) > 0 && viewProjectIds != "" {
		projectIds = append(projectIds, strings.Split(viewProjectIds, ",")...)
	}

	projectIds = utils.NewArray().ArrayUnique(projectIds)
	taskIds = utils.NewArray().ArrayUnique(taskIds)
	tasks, err := models.TaskModel.GetTaskByProjectIdsAndTaskIds(projectIds, taskIds)
	if err != nil {
		this.jsonError("获取任务失败！")
	}
	taskValue := make(map[string]interface{})
	result := make(map[string]map[string]int)
	timePattern := "2006-01-02 15:04:05"
	//初始化task数据
	for _, task := range tasks {
		result[task["task_id"]] = make(map[string]int)
		result[task["task_id"]]["finish"] = 0
		result[task["task_id"]]["doing"] = 0
		result[task["task_id"]]["result"] = models.TASKLOG_SUCCESS
		createTime := utils.NewConvert().StringToInt64(task["create_time"])
		task["create_time"] = time.Unix(createTime, 0).Format(timePattern)
		taskValue[task["task_id"]] = task
	}
	excutingTaskIds := utils.NewArray().ArrayColumn(tasks, "task_id")

	projectIds = utils.NewArray().ArrayColumn(tasks, "project_id")
	projectIds = utils.NewArray().ArrayUnique(projectIds)
	projects, err := models.ProjectModel.GetProjectByProjectIds(projectIds)
	if err != nil {
		this.jsonError("获取项目信息失败！")
	}
	projectsValue := utils.NewArray().ChangeKey(projects, "project_id")

	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskIds(excutingTaskIds)
	if err != nil {
		this.jsonError("通过任务ID获取任务日志失败！")
	}
	for _, taskLog := range taskLogs {
		if taskLog["status"] == "2" {
			result[taskLog["task_id"]]["finish"]++
		} else {
			result[taskLog["task_id"]]["doing"]++
		}
		if taskLog["status"] == "2" && taskLog["is_success"] == "0" {
			result[taskLog["task_id"]]["result"] = models.TASKLOG_FAILED
		}
	}
	if err != nil {
		this.jsonError("获取任务日志失败！")
	}

	data := make(map[string]interface{})
	data["project"] = projectsValue
	data["task"] = taskValue
	data["result"] = result
	this.jsonSuccess("查询正在执行的任务成功！", data, "")
}

func (this *TaskController) Task() {
	page, _ := this.GetInt("page", 1)
	userId := this.GetString("user_id", "")
	projectId := this.GetString("project_id", "")
	if projectId == "" {
		this.viewError("项目ID参数出错！")
	}
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil || len(project) == 0 {
		this.viewError("获取项目信息失败！")
	}

	var count int64
	var tasks []map[string]string
	number := 20
	limit := (page - 1) * number
	if userId != "" {
		count, err = models.TaskModel.CountTaskByProjectIdAndUserId(projectId, userId)
		tasks, err = models.TaskModel.GetTaskByProjectIdAndUserId(projectId, userId, limit, number)
	} else {
		count, err = models.TaskModel.CountTaskByProjectId(projectId)
		tasks, err = models.TaskModel.GetTaskByProjectId(projectId, limit, number)
	}
	if err != nil {
		this.viewError("获取任务信息出错！")
	}

	this.Data["tasks"] = tasks
	this.Data["userId"] = userId
	this.Data["project"] = project
	this.SetPaginator(number, count)
	this.viewLayoutTitle("项目任务信息", "task/task", "page")
}

func (this *TaskController) TaskLog() {
	taskId := this.GetString("task_id", "")
	this.Data["taskId"] = taskId
	this.viewLayoutTitle("项目任务信息", "task/task-log", "page")
}

func (this *TaskController) GetTaskLog() {
	taskId := this.GetString("task_id", "")
	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskId(taskId)
	if err != nil {
		this.jsonError("获取任务日志出错！")
	}

	//获取nodeIds
	nodeIds := make([]string, len(taskLogs))
	timePattern := "2006-01-02 15:04:05"
	for index, taskLog := range taskLogs {
		nodeIds[index] = taskLog["node_id"]
		updateTime := utils.NewConvert().StringToInt64(taskLog["update_time"])
		taskLog["update_time"] = time.Unix(updateTime, 0).Format(timePattern)
		createTime := utils.NewConvert().StringToInt64(taskLog["create_time"])
		taskLog["create_time"] = time.Unix(createTime, 0).Format(timePattern)
	}

	nodesValue, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.jsonError("获取节点信息出错！")
	}
	nodes := utils.NewArray().ChangeKey(nodesValue, "node_id")

	data := make(map[string]interface{})
	data["nodes"] = nodes
	data["taskLogs"] = taskLogs
	this.jsonSuccess("查询任务日志成功！", data, "")
}
