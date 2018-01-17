package container

import (
	"bzppx-codepub/app/remotes"
	"log"
	"time"
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
				err := remotes.Task.Publish(agentMsg.Ip, agentMsg.Port, agentMsg.Args)
				if err != nil {
					log.Println(err.Error())
				}else {
					w.SendGetStatusChan(agentMsg)
				}
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
				for {
					isFinish, err := remotes.Task.GetResults(agentMsg.Ip, agentMsg.Port, agentMsg.Args)
					if err != nil {
						log.Println(err.Error())
					}
					if isFinish {
						break
					}
					time.Sleep(1 * time.Second)
				}
			}(agentMsg)
		}
	}
}

// 初始化
func init()  {
	go Worker.StartPublish()
	go Worker.StartGetStatus()
}