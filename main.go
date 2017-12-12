package main

import (
	_ "bzppx-codepub/app/routers"
	_ "bzppx-codepub/app/models"
	"github.com/astaxie/beego"
	"log"
)

func main() {

	initLog()
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
