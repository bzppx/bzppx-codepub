package controllers

type UserController struct {
	BaseController
}

func (this *UserController) Create() {

	this.viewLayoutTitle("新增用户", "user/form", "page")
}
func (this *UserController) Default() {
	this.viewLayoutTitle("首页", "main/default", "page")
}
