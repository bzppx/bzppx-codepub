package container

import (
	"bzppx-codepub/app/remotes"
	"bzppx-codepub/app/models"
	"errors"
)

var Worker = NewWorker()

type worker struct {
	publishChan chan AgentMessage
	statusChan chan AgentMessage
}

type AgentMessage struct {
	Ip string
	Port string
	Args map[string]interface{}
}

func NewWorker() *worker {
	return &worker{
		publishChan: make(chan AgentMessage, 1),
		statusChan: make(chan AgentMessage, 1),
	}
}

func (w *worker) SendPublishChan(agentMsg AgentMessage)  {
	w.publishChan <- agentMsg
}

func (w *worker) SendGetStatusChan(agentMsg AgentMessage)  {
	w.statusChan <- agentMsg
}

// 初始化 agent 执行 worker
func (w *worker) StartPublish() {
	for {
		select {
		case agentMsg := <-w.publishChan:
			go func(agentMsg AgentMessage) {
				remotes.Git.Publish(agentMsg.Ip, agentMsg.Port, agentMsg.Args)
			}(agentMsg)
		}
	}
}

// 初始化获取 agent 状态 worker
func (w *worker) StartGetStatus() {
	for {
		select {
		case agentMsg := <-w.statusChan:
			go func(agentMsg AgentMessage) {
				remotes.Git.GetResults(agentMsg.Ip, agentMsg.Port, agentMsg.Args)
			}(agentMsg)
		}
	}
}

// 根据 taskLogId 获取 AgentMessage
func (w *worker) GetAgentMessageByTaskLogId(taskLogId string) (agentMessage AgentMessage, err error) {

	taskLog, err := models.TaskLogModel.GetTaskLogByTaskLogId(taskLogId)
	if err != nil {
		return
	}
	if len(taskLog) == 0 {
		return agentMessage, errors.New("task log "+taskLogId+" not exist")
	}
	taskId := taskLog["task_id"]
	nodeId := taskLog["node_id"]

	// task info
	task, err := models.TaskModel.GetTaskByTaskId(taskId)
	if err != nil {
		return
	}
	if len(task) == 0 {
		return agentMessage, errors.New("task "+taskId+" not exist")
	}
	moduleId := task["module_id"]

	// module info
	module, err := models.ModuleModel.GetModuleByModuleId(moduleId)
	if err != nil {
		return
	}
	if len(module) == 0 {
		return agentMessage, errors.New("module "+moduleId+" not exist")
	}

	// node info
	node, err := models.NodeModel.GetNodeByNodeId(nodeId)
	if err != nil {
		return
	}
	if len(node) == 0 {
		return agentMessage, errors.New("node "+nodeId+" not exist")
	}

	agentMessage = AgentMessage{
		Ip: node["ip"],
		Port: node["port"],
		Args:map[string]interface{}{
			"url": module["repository_url"],
			"ssh_key": module["ssh_key"],
			"ssh_key_salt": module["ssh_key_salt"],
			"path": module["code_path"],
			"branch": module["branch"],
			"username": module["https_username"],
			"password": module["https_password"],
		},
	}
	return
}

// 初始化
func init()  {
	go Worker.StartPublish()
	go Worker.StartGetStatus()
}