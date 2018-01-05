package controllers

type TaskController struct {
	BaseController
}

func (this *TaskController) Center() {
	this.viewLayoutTitle("任务队列", "task/center", "page")
}

func (this *TaskController) Node() {
	this.viewLayoutTitle("节点进度", "task/node", "page")
}