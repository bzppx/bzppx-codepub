package main

import (
	_ "bzppx-codepub/app/routers"
	_ "bzppx-codepub/app/models"
	"github.com/astaxie/beego"
	"log"
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/container"
)

func main() {

	initLog()
	initTask()
	beego.Run()
}

func initLog() {
	logConfigs, err := beego.AppConfig.GetSection("log")
	if err != nil {
		log.Println(err.Error())
	}
	for adapter, config := range logConfigs {
		beego.SetLogger(adapter, config)
	}
	beego.SetLogFuncCall(true)
}

func initTask() {

	// 获取所有的创建的 task log
	taskLogs, err := models.TaskLogModel.GetTaskLogByStatus(models.TASKLOG_STATUS_CREATE)
	if err != nil {
		panic(err)
	}
	if len(taskLogs) > 0 {
		for _, taskLog := range taskLogs {
			agentMessage, err := container.Worker.GetAgentMessageByTaskLogId(taskLog["task_log_id"])
			if err != nil {
				log.Println(err.Error())
			}
			container.Worker.SendPublishChan(agentMessage)
		}
	}
}
