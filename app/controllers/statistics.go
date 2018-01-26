package controllers

import (
	"bzppx-codepub/app/models"
)

type StatisticsController struct {
	BaseController
}

const (
	TASKLOG_STATUS_FINISH = "2"

	TASKLOG_SUCCESS = "1"
)

func (this *StatisticsController) Statistics() {
	nodeNumber, err := models.NodeModel.CountNodes()
	if err != nil {
		this.viewError(err.Error(), "statistics/statistics")
	}

	nodesGroupNumber, err := models.NodesModel.CountNodeGroups()
	if err != nil {
		this.viewError(err.Error(), "statistics/statistics")
	}

	projectNumber, err := models.ProjectModel.CountProjects()
	if err != nil {
		this.viewError(err.Error(), "statistics/statistics")
	}

	projectGroupNumber, err := models.GroupModel.CountGroups()
	if err != nil {
		this.viewError(err.Error(), "statistics/statistics")
	}

	tasks, err := models.TaskModel.GetAllTask()
	if err != nil {
		this.viewError(err.Error(), "statistics/statistics")
	}
	tasksChangKey := make(map[string]map[string]string, len(tasks))
	for _, task := range tasks {
		task["status"] = "success"
		tasksChangKey[task["task_id"]] = task
	}

	taskLogs, err := models.TaskLogModel.GetAllTaskLog()
	if err != nil {
		this.viewError(err.Error(), "statistics/statistics")
	}

	taskNumber := map[string]int{
		"success": 0,
		"failed":  0,
		"doing":   0,
	}
	taskLogNumber := map[string]int{
		"success": 0,
		"failed":  0,
		"doing":   0,
	}

	//处理 task_log 数量和 task 的状态
	for _, taskLog := range taskLogs {
		if taskLog["status"] == TASKLOG_STATUS_FINISH {
			if taskLog["is_success"] == TASKLOG_SUCCESS {
				taskLogNumber["success"]++
			} else {
				if tasksChangKey[taskLog["task_id"]]["status"] != "doing" {
					tasksChangKey[taskLog["task_id"]]["status"] = "failed"
				}
				taskLogNumber["failed"]++
			}
		} else {
			tasksChangKey[taskLog["task_id"]]["status"] = "doing"
			taskLogNumber["doing"]++
		}
	}

	for _, task := range tasksChangKey {
		switch task["status"] {
		case "success":
			taskNumber["success"]++
		case "failed":
			taskNumber["failed"]++
		case "doing":
			taskNumber["doing"]++
		}
	}

	this.Data["node_number"] = nodeNumber
	this.Data["node_group_number"] = nodesGroupNumber
	this.Data["project_number"] = projectNumber
	this.Data["project_group_number"] = projectGroupNumber
	this.Data["task_log_number"] = taskLogNumber
	this.Data["task_number"] = taskNumber
	this.viewLayoutTitle("统计列表", "statistics/statistics", "page")
}
