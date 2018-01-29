package container

import (
	"github.com/astaxie/beego"
	"bzppx-codepub/app/models"
	"time"
)

//监控后台数据库未提交的数据

func NewMonitor() *Monitor {
	return &Monitor{}
}

type Monitor struct {
	
}

// 每 5 s 监控一次是否有没有提交的(暂时取消)
func (m *Monitor) MonitorCreateStatus()  {
	for {
		m.HandleCreateStatusTaskLog()
		time.Sleep(10 * time.Second)
	}
}

// 程序启动监控已创建的任务日志
func (m *Monitor) HandleCreateStatusTaskLog() {
	taskLogs, err := models.TaskLogModel.GetTaskLogByStatus(models.TASKLOG_STATUS_CREATE)
	if err != nil {
		beego.Error(err.Error())
		return
	}
	if len(taskLogs) == 0 {
		return
	}
	nodeIds := []string{}
	taskIds := []string{}
	for _, taskLog := range taskLogs {
		nodeIds = append(nodeIds, taskLog["node_id"])
		taskIds = append(taskIds, taskLog["task_id"])
	}
	tasks, err := models.TaskModel.GetTaskByTaskIds(taskIds)
	if len(tasks) == 0 {
		return
	}
	if err != nil {
		beego.Error(err.Error())
		return
	}
	projectIds := []string{}
	for _, task := range tasks {
		projectIds = append(projectIds, task["project_id"])
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		beego.Error(err.Error())
		return
	}
	projects, err := models.ProjectModel.GetProjectByProjectIds(projectIds)
	if err != nil {
		beego.Error(err.Error())
		return
	}

	for _, taskLog := range taskLogs {
		ip := ""
		port := ""
		projectId := "0"
		sha1Id := ""
		for _, task := range tasks {
			if task["task_id"] == taskLog["task_id"] {
				projectId = task["project_id"]
				sha1Id = task["sha1_id"]
				break
			}
		}
		project := map[string]string{}
		for _, projectItem := range projects {
			if projectItem["project_id"] == projectId {
				project = projectItem
				break
			}
		}
		// sha1d 不为空，为回滚操作
		if sha1Id == "" {
			sha1Id = project["branch"]
		}
		args := map[string]interface{}{
			"task_log_id":  taskLog["task_log_id"],
			"url":          project["repository_url"],
			"ssh_key":      project["ssh_key"],
			"ssh_key_salt": project["ssh_key_salt"],
			"path":         project["code_path"],
			"branch":       sha1Id,
			"username":     project["https_username"],
			"password":     project["https_password"],
			"pre_command":                  project["pre_command"],
			"pre_command_exec_type":        project["pre_command_exec_type"],
			"pre_command_exec_timeout":     project["pre_command_exec_timeout"],
			"post_command":                 project["post_command"],
			"post_command_exec_type":       project["post_command_exec_type"],
			"post_command_exec_timeout":    project["post_command_exec_timeout"],
			"exec_user":                    project["exec_user"],
		}
		for _, node := range nodes {
			if node["node_id"] == taskLog["node_id"] {
				ip = node["ip"]
				port = node["port"]
				break
			}
		}
		agentMessage := AgentMessage{
			Ip:   ip,
			Port: port,
			Args: args,
		}
		Worker.SendPublishChan(agentMessage)
	}
}