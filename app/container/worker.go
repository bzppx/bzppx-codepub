package container

import (
	"bzppx-codepub/app/remotes"
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

func (w *worker) StartPublish(agentMsg AgentMessage)  {
	w.publishChan <- agentMsg
}

func (w *worker) StartGetStatus(agentMsg AgentMessage)  {
	w.statusChan <- agentMsg
}

// 初始化 agent 执行 worker
func (w *worker) InitPublish() {
	for {
		select {
		case agentMsg := <-w.publishChan:
			go func(agentMsg AgentMessage) {
				address := agentMsg.Ip + ":"+agentMsg.Port
				remotes.Git.GetCommitId(address, agentMsg.Args)
				remotes.Git.Publish(address, agentMsg.Args)
			}(agentMsg)
		}
	}
}

// 初始化获取 agent 状态 worker
func (w *worker) InitStatus() {
	for {
		select {
		case agentMsg := <-w.statusChan:
			go func(agentMsg AgentMessage) {
				address := agentMsg.Ip + ":" + agentMsg.Port
				remotes.Git.GetResults(address, agentMsg.Args)
			}(agentMsg)
		}
	}
}

// 初始化
func init()  {
	go Worker.InitPublish()
	go Worker.InitStatus()
}