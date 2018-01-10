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
	viewModuleIds := this.GetString("module_id", "")
	viewTaskIds := this.GetString("task_id", "")
	var moduleIds []string
	taskIds, err := models.TaskLogModel.GetExcutingTaskIdByTaskLog()
	if err != nil {
		this.jsonError("获取taskId失败！")
	}
	if viewTaskIds != "" {
		taskIds = append(taskIds, strings.Split(viewTaskIds, ",")...)
	}

	if !this.isAdmin() && !this.isRoot() {
		moduleUsers, err := models.UserModuleModel.GetUserModuleByUserId(this.UserID)
		if err != nil {
			this.jsonError("获取权限失败！")
		}
		if len(moduleUsers) == 0 {
			this.jsonError("此账号没有被授予任何模块权限！")
		}
		moduleIds = utils.NewArray().ArrayColumn(moduleUsers, "module_id")
	}
	if len(moduleIds) > 0 && viewModuleIds != "" {
		moduleIds = append(moduleIds, strings.Split(viewModuleIds, ",")...)
	}

	moduleIds = utils.NewArray().ArrayUnique(moduleIds)
	taskIds = utils.NewArray().ArrayUnique(taskIds)
	tasks, err := models.TaskModel.GetTaskByModuleIdsAndTaskIds(moduleIds, taskIds)
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

	moduleIds = utils.NewArray().ArrayColumn(tasks, "module_id")
	moduleIds = utils.NewArray().ArrayUnique(moduleIds)
	modules, err := models.ModuleModel.GetModuleByModuleIds(moduleIds)
	if err != nil {
		this.jsonError("获取模块信息失败！")
	}
	modulesValue := utils.NewArray().ChangeKey(modules, "module_id")

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
	data["module"] = modulesValue
	data["task"] = taskValue
	data["result"] = result
	this.jsonSuccess("查询正在执行的任务成功！", data, "")
}
