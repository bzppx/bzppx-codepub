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

	nodeGroups, err := models.NodesModel.GetNodeGroups();
	if err != nil {
		this.RecordLog("获取节点组失败："+err.Error())
		this.viewError("获取节点组失败！", "/node/list")
	}

	this.Data["nodeGroups"] = nodeGroups
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
	nodeGroupIds := this.GetStrings("nodes_ids", []string{})

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
	if len(nodeGroupIds) == 0 {
		this.jsonError("请至少选择一个节点组！")
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
	}
	this.RecordLog("保存节点 " + utils.NewConvert().IntToString(nodeId, 10) + " 成功")

	// 绑定节点组和几点关系
	var insertValues []map[string]interface{}
	for _, nodeGroupId := range nodeGroupIds {
		insertValue := map[string]interface{}{
			"node_id": nodeId,
			"nodes_id": nodeGroupId,
			"create_time": time.Now().Unix(),
		}
		insertValues = append(insertValues, insertValue)
	}
	_, err = models.NodeNodesModel.InsertBatch(insertValues)
	if err != nil {
		this.RecordLog("添加节点绑定节点节点组关系失败: " + err.Error())
		this.jsonError("添加节点失败！")
	}

	this.RecordLog("添加节点绑定节点节点组关系成功")
	this.jsonSuccess("添加节点成功", nil, "/node/list")
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

	nodeGroups, err := models.NodesModel.GetNodeGroups();
	if err != nil {
		this.RecordLog("获取节点组失败："+err.Error())
		this.viewError("获取节点组失败！", "/node/list")
	}

	nodeNodes, err := models.NodeNodesModel.GetNodeNodesByNodeId(nodeId)
	if err != nil {
		this.RecordLog("获取节点节点组关系失败："+err.Error())
		this.viewError("获取节点节点组关系失败！", "/node/list")
	}

	var newNodeGroups []map[string]string
	for _, nodeGroup := range nodeGroups {
		newNodeGroup := map[string]string{
			"nodes_id": nodeGroup["nodes_id"],
			"name": nodeGroup["name"],
			"comment": nodeGroup["comment"],
			"is_default": "0",
		}
		for _, nodeNode := range nodeNodes {
			if nodeGroup["nodes_id"] == nodeNode["nodes_id"] {
				newNodeGroup["is_default"] = "1"
			}
		}
		newNodeGroups = append(newNodeGroups, newNodeGroup)
	}

	this.Data["node"] = node
	this.Data["nodeGroups"] = newNodeGroups
	this.viewLayoutTitle("修改节点组", "node/form", "page")
}

// 修改保存
func (this *NodeController) Modify() {

	nodeId := this.GetString("node_id", "")
	ip := strings.Trim(this.GetString("ip", ""), "")
	port, _ := this.GetInt("port", 0)
	comment := strings.Trim(this.GetString("comment", ""), "")
	nodeGroupIds := this.GetStrings("nodes_ids", []string{})

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
	if len(nodeGroupIds) == 0 {
		this.jsonError("请至少选择一个节点组！")
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
	}
	// 重新绑定节点组和几点关系
	// 先删除
	err = models.NodeNodesModel.DeleteNodeNodesByNodeId(nodeId)
	if err != nil {
		this.RecordLog("删除节点 " + nodeId + " 与节点组对应关系失败: " + err.Error())
		this.jsonError("修改节点失败！")
	}
	var insertValues []map[string]interface{}
	for _, nodeGroupId := range nodeGroupIds {
		insertValue := map[string]interface{}{
			"node_id": nodeId,
			"nodes_id": nodeGroupId,
			"create_time": time.Now().Unix(),
		}
		insertValues = append(insertValues, insertValue)
	}
	_, err = models.NodeNodesModel.InsertBatch(insertValues)
	if err != nil {
		this.RecordLog("修改节点绑定节点节点组关系失败: " + err.Error())
		this.jsonError("修改节点失败！")
	}

	this.RecordLog("修改节点 " + nodeId + " 成功")
	this.jsonSuccess("修改节点成功", nil, "/node/list")

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
	err = models.NodeNodesModel.DeleteNodeNodesByNodeNodesId(nodeId)
	if err != nil {
		this.RecordLog("删除节点节点组关系，节点ID： " + nodeId + " 失败: " + err.Error())
		this.jsonError("删除节点节点组关系失败！")
	}
	this.RecordLog("删除节点 " + nodeId + " 成功")
	this.jsonSuccess("删除节点成功", nil, "/node/list")
}
