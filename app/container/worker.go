package container

import (
	"bzppx-codepub/app/remotes"
	"time"
	"bzppx-codepub/app/models"
	"github.com/astaxie/beego"
)

var Worker = NewWorker()

type worker struct {
	publishChan chan AgentMessage
	statusChan chan AgentMessage
}

type AgentMessage struct {
	Ip string
	Port string
	Token string
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
				defer func() {
					err := recover()
					if err != nil {
						beego.Error(err)
					}
				}()
				err := remotes.Task.Publish(agentMsg.Ip, agentMsg.Port, agentMsg.Token, agentMsg.Args)
				if err != nil {
					beego.Error(err.Error())
					w.UpdateResult(agentMsg.Args["task_log_id"].(string), err.Error())
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
				defer func() {
					err := recover()
					if err != nil {
						beego.Error(err)
					}
				}()
				for {
					isFinish, err := remotes.Task.GetResults(agentMsg.Ip, agentMsg.Port, agentMsg.Token,agentMsg.Args)
					if err != nil {
						beego.Error(err.Error())
						w.UpdateResult(agentMsg.Args["task_log_id"].(string), err.Error())
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

func (t *worker) UpdateResult(taskLogId string, result string) {
	update := map[string]interface{}{
		"result": result,
		"update_time": time.Now().Unix(),
	}
	_, err := models.TaskLogModel.Update(taskLogId, update)
	if err != nil {
		beego.Error("update task_log result error: "+ err.Error())
	}
}

// 初始化
func InitWorker()  {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				beego.Error(err)
			}
		}()
		Worker.StartPublish()
	}()
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				beego.Error(err)
			}
		}()
		Worker.StartGetStatus()
	}()
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				beego.Error(err)
			}
		}()
		NewMonitor().HandleCreateStatusTaskLog()
	}()
}