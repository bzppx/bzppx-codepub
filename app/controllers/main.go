package controllers

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	this.Data["isAdmin"] = (this.isAdmin() || this.isRoot())
	this.viewLayoutTitle("CodePub POWVEREDBY BZPPX", "main/index", "main")
}
func (this *MainController) Default() {
	this.viewLayoutTitle("首页", "main/default", "page")
}
func (this *MainController) Tpl() {
	typ := this.GetString("type")
	this.viewLayoutTitle("模板", "main/tpl-"+typ, "page")
}
