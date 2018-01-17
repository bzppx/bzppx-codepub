package controllers

import (
	"strings"
	"bzppx-codepub/app/models"
	"time"
	"bzppx-codepub/app/utils"
)

type NodesController struct {
	BaseController
}

// 添加节点组
func (this *NodesController) Add() {
	this.viewLayoutTitle("添加节点组", "nodes/form", "page")
}

// 保存节点组
func (this *NodesController) Save() {

	name := strings.Trim(this.GetString("name", ""), "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if name == "" {
		this.jsonError("节点组名称不能为空！")
	}

	nodeGroup, err := models.NodesModel.HasNodesName(name)
	if err != nil {
		this.ErrorLog("查找节点组 "+name+" 失败: "+err.Error())
		this.jsonError("添加节点组失败！")
	}
	if nodeGroup {
		this.jsonError("节点组名称已存在！")
	}

	nodesValue := map[string]interface{}{
		"name": name,
		"comment": comment,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	nodesId, err := models.NodesModel.Insert(nodesValue)
	if err != nil {
		this.ErrorLog("添加节点组失败: "+err.Error())
		this.jsonError("添加节点组失败！")
	}else {
		this.InfoLog("添加节点组 "+utils.NewConvert().IntToString(nodesId, 10)+" 成功")
		this.jsonSuccess("添加节点组成功", nil, "/nodes/list")
	}
}

// 节点组列表
func (this *NodesController) List() {

	page, _:= this.GetInt("page", 1)
	keyword := strings.Trim(this.GetString("keyword", ""), "")

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var nodeGroups []map[string]string
	if (keyword != "") {
		count, err = models.NodesModel.CountNodeGroupsByKeyword(keyword)
		nodeGroups, err = models.NodesModel.GetNodeGroupsByKeywordAndLimit(keyword, limit, number)
	}else {
		count, err = models.NodesModel.CountNodeGroups()
		nodeGroups, err = models.NodesModel.GetNodeGroupsByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("查找节点组列表失败: "+err.Error())
		this.viewError(err.Error(), "/nodes/list")
	}

	this.Data["nodeGroups"] = nodeGroups
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)
	this.viewLayoutTitle("节点组列表", "nodes/list", "page")
}

// 修改
func (this *NodesController) Edit() {

	nodesId := this.GetString("nodes_id", "")
	if nodesId == "" {
		this.viewError("节点组不存在", "/nodes/list")
	}

	nodeGroup, err := models.NodesModel.GetNodeGroupByNodesId(nodesId)
	if err != nil {
		this.ErrorLog("查找节点组 "+nodesId+" 失败: "+err.Error())
		this.viewError("节点组不存在", "/nodes/list")
	}

	this.Data["nodeGroup"] = nodeGroup
	this.viewLayoutTitle("修改节点组", "nodes/form", "page")
}

// 修改保存
func (this *NodesController) Modify() {

	nodesId := this.GetString("nodes_id", "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if nodesId == "" {
		this.jsonError("节点组不存在！")
	}

	nodes, err := models.NodesModel.GetNodeGroupByNodesId(nodesId)
	if err != nil {
		this.ErrorLog("查找节点组 "+nodesId+" 失败: "+err.Error())
		this.jsonError("节点组不存在！")
	}
	if len(nodes) == 0 {
		this.jsonError("节点组不存在！")
	}

	nodesValue := map[string]interface{}{
		"comment": comment,
		"update_time": time.Now().Unix(),
	}

	_, err = models.NodesModel.Update(nodesId, nodesValue)
	if err != nil {
		this.ErrorLog("修改节点组 "+nodesId+" 失败: "+err.Error())
		this.jsonError("修改节点组失败！")
	}else {
		this.InfoLog("修改节点组 "+nodesId+" 成功")
		this.jsonSuccess("修改节点组成功", nil, "/nodes/list")
	}
}

// 删除
func (this *NodesController) Delete() {

	nodesId := this.GetString("nodes_id", "")

	if nodesId == "" {
		this.jsonError("没有选择节点组！")
	}

	nodes, err := models.NodesModel.GetNodeGroupByNodesId(nodesId)
	if err != nil {
		this.ErrorLog("查找节点组 "+nodesId+" 失败: "+err.Error())
		this.jsonError("节点组不存在！")
	}
	if len(nodes) == 0 {
		this.jsonError("节点组不存在！")
	}

	// todo 判断节点组下的项目是否需要一起删除

	nodesValue := map[string]interface{}{
		"is_delete": models.NODES_DELETE,
		"update_time": time.Now().Unix(),
	}

	_, err = models.NodesModel.Update(nodesId, nodesValue)
	if err != nil {
		this.ErrorLog("删除节点组 "+nodesId+" 失败: "+err.Error())
		this.jsonError("删除节点组失败！")
	}

	this.InfoLog("删除节点组 "+nodesId+" 成功")
	this.jsonSuccess("删除节点组成功", nil, "/nodes/list")
}

// 节点列表
func (this *NodesController) Node() {

	nodesId := this.GetString("nodes_id", "")
	if nodesId == "" {
		this.viewError("节点组不存在", "/nodes/list")
	}

	nodeGroup, err := models.NodesModel.GetNodeGroupByNodesId(nodesId)
	if err != nil {
		this.ErrorLog("查找节点组 "+nodesId+" 失败: "+err.Error())
		this.viewError("节点组不存在", "/nodes/list")
	}

	// 查找该节点组下的节点
	nodeNodes, err := models.NodeNodesModel.GetNodeNodesByNodesId(nodesId)
	if err != nil {
		this.ErrorLog("查找节点组 "+nodesId+" 下节点失败: "+err.Error())
		this.viewError("查找节点错误", "/nodes/list")
	}
	if len(nodeNodes) == 0 {
		this.viewError("该节点组无节点", "/nodes/list")
	}
	var nodeIds []string
	for _, nodeNode := range nodeNodes {
		nodeIds = append(nodeIds, nodeNode["node_id"])
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.ErrorLog("查找节点失败: "+err.Error())
		this.viewError("查找节点错误", "/nodes/list")
	}

	// 查找除 nodeIds 外的节点
	otherNodes, err := models.NodeModel.GetNodeByNotNodeIds(nodeIds)
	if err != nil {
		this.ErrorLog("查找节点失败: "+err.Error())
		this.viewError("查找节点错误", "/nodes/list")
	}

	this.Data["nodes"] = nodes
	this.Data["otherNodes"] = otherNodes
	this.Data["nodeGroup"] = nodeGroup
	this.viewLayoutTitle("节点组节点", "nodes/node", "page")
}

// 导入节点
func (this *NodesController) ImportNode() {

	nodesId := this.GetString("nodes_id", "")
	nodeIds := this.GetStrings("node_id", []string{})

	if nodesId == "" {
		this.jsonError("节点组不存在！")
	}
	if len(nodeIds) == 0 {
		this.jsonError("请选择节点！")
	}
	nodeGroup, err := models.NodesModel.GetNodeGroupByNodesId(nodesId)
	if err != nil {
		this.ErrorLog("查找节点组 "+nodesId+" 失败: "+err.Error())
		this.jsonError("节点组不存在！")
	}
	if len(nodeGroup) == 0 {
		this.jsonError("节点组不存在！")
	}

	var insertValues []map[string]interface{}
	for _, nodeId := range nodeIds {
		insertValue := map[string]interface{}{
			"node_id": nodeId,
			"nodes_id": nodesId,
			"create_time": time.Now().Unix(),
		}
		insertValues = append(insertValues, insertValue)
	}
	_, err = models.NodeNodesModel.InsertBatch(insertValues)
	if err != nil {
		this.ErrorLog("节点组导入节点失败: " + err.Error())
		this.jsonError("导入节点失败！")
	}

	this.InfoLog("节点组 "+nodesId+" 导入节点 "+strings.Join(nodeIds, ",")+" 成功" )
	this.jsonSuccess("导入节点成功！", nil, "/nodes/node?nodes_id="+nodesId)
}

// 导入节点
func (this *NodesController) Remove() {

	nodesId := this.GetString("nodes_id", "")
	nodeId := this.GetString("node_id", "")

	if nodesId == "" {
		this.jsonError("节点组不存在！")
	}
	if nodeId == "" {
		this.jsonError("节点不存在！")
	}

	err := models.NodeNodesModel.DeleteByNodeIdAndNodesId(nodeId, nodesId)
	if err != nil {
		this.ErrorLog("移除节点组 "+nodesId+" 下节点 "+nodeId+" 失败：" + err.Error())
		this.jsonError("移除节点失败！")
	}

	this.InfoLog("移除节点组 "+nodesId+" 下节点 "+nodeId+" 成功" )
	this.jsonSuccess("移除节点成功！", nil, "/nodes/node?nodes_id="+nodesId)
}