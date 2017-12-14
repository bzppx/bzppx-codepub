package controllers

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	this.viewLayoutTitle("CodePub POWVEREDBY BZPPX", "main/index", "main")
}
func (this *MainController) Default() {
	this.viewLayoutTitle("首页", "main/default", "page")
}
