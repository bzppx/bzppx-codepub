package controllers

type PublishController struct {
	BaseController
}

func (this *PublishController) Module() {
	this.viewLayoutTitle("我的资料", "publish/module", "page")
}



