package controllers

import (
	"bzppx-codepub/app/models"
)

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	this.Data["isAdmin"] = (this.isAdmin() || this.isRoot())
	this.viewLayoutTitle("CodePub POWVEREDBY BZPPX", "main/index", "main")
}
func (this *MainController) Default() {
	var err error
	this.Data["notices"], err = models.NoticeModel.GetNoticesByLimit(0, 5)
	if err != nil {
		this.ErrorLog("获取最新公告失败：" + err.Error())
		this.viewError("获取最新公告失败")
	}
	this.viewLayoutTitle("首页", "main/default", "page")
}
func (this *MainController) Tpl() {
	typ := this.GetString("type")
	this.viewLayoutTitle("模板", "main/tpl-"+typ, "page")
}
