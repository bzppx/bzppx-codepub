package main

import (
	_ "bzppx-codepub/app/routers"
	_ "bzppx-codepub/app/models"
	"bzppx-codepub/app/container"
	"github.com/astaxie/beego"
	"log"
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
	go func() {
		err := recover()
		if err != nil {
			beego.Error(err)
		}
		container.NewMonitor().MonitorCreateStatus()
	}()
}
