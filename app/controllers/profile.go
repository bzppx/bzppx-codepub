package controllers

type ProfileController struct {
	BaseController
}

// 个人资料
func (this *ProfileController) My() {
	this.viewLayoutTitle("我的资料", "profile/form", "page")
}


