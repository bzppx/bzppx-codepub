package controllers

type MainController struct {
	BaseController
}

// return json data
func (this *MainController) Index() {
	this.viewLayoutTitle("CodePub POWVEREDBY BZPPX", "main/index", "main")
}
