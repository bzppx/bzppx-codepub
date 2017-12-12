package controllers

type LoginController struct {
	BaseController
}

func (this *LoginController) Index() {
	this.viewLayoutTitle("Login", "login/login", "login")
}
