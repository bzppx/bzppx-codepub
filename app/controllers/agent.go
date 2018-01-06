package controllers

// 用于测试

import (
	"github.com/astaxie/beego"
	"bzppx-codepub/app/container"
)

type AgentController struct {
	beego.Controller
}

func (this *AgentController) Publish() {
	a := container.AgentMessage{
		Ip: "127.0.0.1",
		Port: "9091",
		Args: map[string]interface{}{
			"a": 2,
			"b": 8,
		},
	}
	container.Worker.StartPublish(a)

	this.Abort(string("okok"))
}