package controllers

import (
	"bzppx-codepub/app/models"
)

type NodeController struct {
	BaseController
}

// 添加节点
func (this *NodeController) Add() {
	this.viewLayoutTitle("添加节点", "node/form", "page")
}

// 节点列表
func (this *NodeController) List() {

	page, _ := this.GetInt("page", 1)
	ip := this.GetString("ip", "")
	keywords := map[string]string{
		"ip": ip,
	}

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var nodes []map[string]string
	if keywords["ip"] != "" {
		count, err = models.NodeModel.CountNodesByKeywords(keywords)
		nodes, err = models.NodeModel.GetNodesByKeywordsAndLimit(keywords, limit, number)
	} else {
		count, err = models.NodeModel.CountNodes()
		nodes, err = models.NodeModel.GetNodesByLimit(limit, number)
	}
	if err != nil {
		this.viewError(err.Error(), "/node/list")
	}
	this.Data["nodes"] = nodes
	this.Data["ip"] = ip
	this.SetPaginator(number, count)

	this.viewLayoutTitle("节点列表", "node/list", "page")
}
