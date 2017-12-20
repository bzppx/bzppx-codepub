package controllers

type ProfileController struct {
	BaseController
}

// 个人资料
func (this *ProfileController) My() {
	this.viewLayoutTitle("我的资料", "profile/form", "page")
}

// 个人资料保存
func (this *ProfileController) Save() {
	this.jsonSuccess("我的资料修改成功", nil, "/profile/my", 3000)
	//this.jsonError("okok", nil)
}
