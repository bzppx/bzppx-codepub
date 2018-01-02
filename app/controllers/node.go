package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"regexp"
	"strings"
	"time"
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

// 保存节点
func (this *NodeController) Save() {

	ip := strings.Trim(this.GetString("ip", ""), "")
	port, _ := this.GetInt("port", 0)
	comment := strings.Trim(this.GetString("comment", ""), "")

	res, err := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip)
	if err != nil {
		this.jsonError("ip正则匹配失败！")
	}
	if !res {
		this.jsonError("ip不正确！")
	}

	if port <= 0 || port >= 65535 {
		this.jsonError("port不正确！")
	}

	has, err := models.NodeModel.HasNodeByIpAndPort("0", ip, port)
	if err != nil {
		this.jsonError("添加节点失败！")
	}
	if has {
		this.jsonError("节点ip和端口已存在！")
	}

	nodeValue := map[string]interface{}{
		"ip":          ip,
		"port":        port,
		"comment":     comment,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	nodeId, err := models.NodeModel.Insert(nodeValue)
	if err != nil {
		this.RecordLog("添加节点失败: " + err.Error())
		this.jsonError("添加节点失败！")
	} else {
		this.RecordLog("添加节点 " + utils.NewConvert().IntToString(nodeId, 10) + " 成功")
		this.jsonSuccess("添加节点成功", nil, "/node/list")
	}
}

// 修改
func (this *NodeController) Edit() {

	nodeId := this.GetString("node_id", "")
	if nodeId == "" {
		this.viewError("节点不存在", "/node/list")
	}

	node, err := models.NodeModel.GetNodeByNodeId(nodeId)
	if err != nil {
		this.viewError("节点不存在", "/node/list")
	}

	this.Data["node"] = node
	this.viewLayoutTitle("修改节点组", "node/form", "page")
}

// 修改保存
func (this *NodeController) Modify() {

	nodeId := this.GetString("node_id", "")
	ip := strings.Trim(this.GetString("ip", ""), "")
	port, _ := this.GetInt("port", 0)
	comment := strings.Trim(this.GetString("comment", ""), "")

	if nodeId == "" {
		this.viewError("节点不存在", "/node/list")
	}

	res, err := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip)
	if err != nil {
		this.jsonError("ip正则匹配失败！")
	}
	if !res {
		this.jsonError("ip不正确！")
	}

	if port <= 0 || port >= 65535 {
		this.jsonError("port不正确！")
	}

	has, err := models.NodeModel.HasNodeByIpAndPort(nodeId, ip, port)
	if err != nil {
		this.jsonError("添加节点失败！")
	}
	if has {
		this.jsonError("节点ip和端口已存在！")
	}

	nodeValue := map[string]interface{}{
		"ip":          ip,
		"port":        port,
		"comment":     comment,
		"update_time": time.Now().Unix(),
	}

	_, err = models.NodeModel.Update(nodeId, nodeValue)
	if err != nil {
		this.RecordLog("修改节点 " + nodeId + " 失败: " + err.Error())
		this.jsonError("修改节点失败！")
	} else {
		this.RecordLog("修改节点 " + nodeId + " 成功")
		this.jsonSuccess("修改节点成功", nil, "/node/list")
	}
}

func (this *NodeController) Delete() {
	nodeId := this.GetString("node_id", "")

	if nodeId == "" {
		this.jsonError("节点不存在！")
	}

	node, err := models.NodeModel.GetNodeByNodeId(nodeId)
	if err != nil {
		this.jsonError("节点不存在！")
	}
	if len(node) == 0 {
		this.jsonError("节点不存在！")
	}

	nodeValue := map[string]interface{}{
		"is_delete":   models.NODE_DELETE,
		"update_time": time.Now().Unix(),
	}

	_, err = models.NodeModel.Update(nodeId, nodeValue)

	if err != nil {
		this.RecordLog("删除节点 " + nodeId + " 失败: " + err.Error())
		this.jsonError("删除节点失败！")
	}
	err = models.ModuleNodeModel.DeleteModuleNodeByNodeId(nodeId)
	if err != nil {
		this.RecordLog("删除模块节点关系，节点ID： " + nodeId + " 失败: " + err.Error())
		this.jsonError("删除模块节点关系失败！")
	}
	err = models.NodeNodesModel.DeleteNodeNodesByNodeId(nodeId)
	if err != nil {
		this.RecordLog("删除节点节点组关系，节点ID： " + nodeId + " 失败: " + err.Error())
		this.jsonError("删除节点节点组关系失败！")
	}
	this.RecordLog("删除节点 " + nodeId + " 成功")
	this.jsonSuccess("删除节点成功", nil, "/node/list")
}
