package main

import (
	_ "sq/api/app/routers"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego"
)

func main() {

	initLog()
	beego.Run()
}

func initLog() {
	logConfigs, err := beego.AppConfig.GetSection("log")
	if err != nil {
		panic("conf log error!")
	}
	for adapter, config := range logConfigs {
		beego.SetLogger(adapter, config)
	}
	beego.SetLogFuncCall(true)
}
