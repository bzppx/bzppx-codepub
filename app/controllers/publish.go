package controllers

type PublishController struct {
	BaseController
}

// 发布模块
func (this *PublishController) Module() {
	this.viewLayoutTitle("发布模块", "publish/module", "page")
}

// 模块信息
func (this *PublishController) Info() {
	this.viewLayoutTitle("模块信息", "module/info", "page")
}

