package controllers

import (
	"bzppx-codepub/app/models"
	"strings"
	"time"
	"bzppx-codepub/app/utils"
)

type ProjectController struct {
	BaseController
}

// 添加项目
func (this *ProjectController) Add() {

	groups, err := models.GroupModel.GetGroups();
	if err != nil {
		this.viewError("获取项目组错误", "project/list")
	}

	this.Data["groups"] = groups
	this.viewLayoutTitle("添加项目", "project/form", "page")
}

// 保存项目
func (this *ProjectController) Save() {

	name := strings.Trim(this.GetString("name", ""), "")
	groupId := strings.Trim(this.GetString("group_id", ""), "")
	pullType := strings.Trim(this.GetString("pull_type", "http"), "")
	repositoryUrl := strings.Trim(this.GetString("repository_url", ""), "")
	sshKey := strings.Trim(this.GetString("ssh_key", ""), "")
	sshKeySalt := strings.Trim(this.GetString("ssh_key_salt", ""), "")
	httpsUsername := strings.Trim(this.GetString("https_username", ""), "")
	httpsPassword := strings.Trim(this.GetString("https_password", ""), "")
	branch := strings.Trim(this.GetString("branch", ""), "")
	codePath := strings.Trim(this.GetString("code_path", ""), "")
	codeDirUser := strings.Trim(this.GetString("code_dir_user", ""), "")
	comment := strings.Trim(this.GetString("comment", ""), "")

	if name == "" {
		this.jsonError("项目名不能为空！")
	}
	if groupId == "" {
		this.jsonError("没有选择项目组！")
	}
	if repositoryUrl == "" {
		this.jsonError("代码仓库地址不能为空！")
	}
	if pullType == "http" {
		if repositoryUrl[0:7] != "http://" {
			this.jsonError("代码仓库地址必须是 http:// 开头！")
		}
	}else if(pullType == "https") {
		if repositoryUrl[0:8] != "https://" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else if(pullType == "ssh") {
		if repositoryUrl[0:4] != "ssh@" {
			this.jsonError("代码仓库地址必须是 ssh@ 开头！")
		}
	}else {
		this.jsonError("无效的代码拉取方式！")
	}

	if pullType == "http" || pullType == "https" {
		if httpsUsername == "" {
			this.jsonError("用户名不能为空！")
		}
		if httpsPassword == "" {
			this.jsonError("密码不能为空！")
		}
	}
	if pullType == "ssh" {
		if sshKey == "" {
			this.jsonError("ssh key 不能为空！")
		}
	}
	if branch == "" {
		this.jsonError("代码分支不能为空！")
	}
	if codePath == "" {
		this.jsonError("代码发布路径不能为空！")
	}
	if codeDirUser == "" {
		this.jsonError("目录所属用户不能为空！")
	}

	project, err := models.ProjectModel.GetProjectByName(name)
	if err != nil {
		this.ErrorLog("添加项目查找项目名是否存在失败: "+err.Error())
		this.jsonError("添加项目失败！")
	}
	if len(project) > 0 {
		this.jsonError("该项目名已存在！")
	}

	projectValue := map[string]interface{}{
		"name": name,
		"user_id": this.UserID,
		"group_id": groupId,
		"repository_url": repositoryUrl,
		"branch": branch,
		"ssh_key": sshKey,
		"ssh_key_salt": sshKeySalt,
		"https_username": httpsUsername,
		"https_password": httpsPassword,
		"code_path": codePath,
		"code_dir_user": codeDirUser,
		"comment": comment,
		"pre_command": "",
		"post_command": "",
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}

	projectId, err := models.ProjectModel.Insert(projectValue)
	if err != nil {
		this.ErrorLog("添加项目插入数据失败: "+err.Error())
		this.jsonError("添加项目失败！")
	}

	this.InfoLog("添加项目 "+utils.NewConvert().IntToString(projectId, 10)+" 成功")
	this.jsonSuccess("添加项目成功, 请继续配置节点", nil, "/project/node?flag=insert&project_id="+utils.NewConvert().IntToString(projectId, 10))
}

// 项目列表
func (this *ProjectController) List() {

	page, _:= this.GetInt("page", 1)
	groupId := this.GetString("group_id", "")
	keyword := strings.Trim(this.GetString("keyword", ""), "")
	keywords := map[string]string {
		"group_id": groupId,
		"keyword": keyword,
	}

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var projects []map[string]string
	if (len(keywords) > 0) {
		count, err = models.ProjectModel.CountProjectsByKeywords(keywords)
		projects, err = models.ProjectModel.GetProjectsByKeywordsAndLimit(keywords, limit, number)
	}else {
		count, err = models.ProjectModel.CountProjects()
		projects, err = models.ProjectModel.GetProjectsByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("项目列表查找项目数据失败: "+err.Error())
		this.viewError("查找项目失败", "/project/list")
	}

	groups, err := models.GroupModel.GetGroups();
	if err != nil {
		this.ErrorLog("项目列表查找项目组数据失败: "+err.Error())
		this.viewError("查找项目失败", "project/list")
	}

	this.Data["projects"] = projects
	this.Data["keywords"] = keywords
	this.Data["groups"] = groups
	this.SetPaginator(number, count)

	this.viewLayoutTitle("项目列表", "project/list", "page")
}

// 修改项目
func (this *ProjectController) Edit() {

	projectId := this.GetString("project_id", "")
	if projectId == "" {
		this.viewError("项目不存在", "/project/list")
	}

	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.viewError("项目不存在", "/project/list")
	}
	if len(project) == 0 {
		this.viewError("项目不存在", "/project/list")
	}
	
	repUrl := project["repository_url"]
	pullType := "http"
	if repUrl[0:4] == "ssh@" {
		pullType = "ssh"
	}else if repUrl[0:8] == "https://" {
		pullType = "https"
	}else {
		pullType = "http"
	}

	groups, err := models.GroupModel.GetGroups();
	if err != nil {
		this.ErrorLog("获取项目组数据失败: "+err.Error())
		this.viewError("获取项目组错误", "/project/list")
	}

	this.Data["project"] = project
	this.Data["pullType"] = pullType
	this.Data["groups"] = groups
	this.viewLayoutTitle("修改项目", "project/form", "page")
}

// 修改保存项目
func (this *ProjectController) Modify() {
	
	projectId := strings.Trim(this.GetString("project_id", ""), "")
	groupId := strings.Trim(this.GetString("group_id", ""), "")
	pullType := strings.Trim(this.GetString("pull_type", "http"), "")
	repositoryUrl := strings.Trim(this.GetString("repository_url", ""), "")
	sshKey := strings.Trim(this.GetString("ssh_key", ""), "")
	sshKeySalt := strings.Trim(this.GetString("ssh_key_salt", ""), "")
	httpsUsername := strings.Trim(this.GetString("https_username", ""), "")
	httpsPassword := strings.Trim(this.GetString("https_password", ""), "")
	branch := strings.Trim(this.GetString("branch", ""), "")
	codePath := strings.Trim(this.GetString("code_path", ""), "")
	codeDirUser := strings.Trim(this.GetString("code_dir_user", ""), "")
	comment := strings.Trim(this.GetString("comment", ""), "")
	
	if projectId == "" {
		this.jsonError("项目id错误！")
	}
	if groupId == "" {
		this.jsonError("没有选择项目组！")
	}
	if repositoryUrl == "" {
		this.jsonError("代码仓库地址不能为空！")
	}
	if pullType == "http" {
		if repositoryUrl[0:7] != "http://" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else if(pullType == "https") {
		if repositoryUrl[0:8] != "https://" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else if(pullType == "ssh") {
		if repositoryUrl[0:4] != "git@" {
			this.jsonError("代码仓库地址必须是 https:// 开头！")
		}
	}else {
		this.jsonError("无效的代码拉取方式！")
	}
	
	if pullType == "http" || pullType == "https" {
		if httpsUsername == "" {
			this.jsonError("用户名不能为空！")
		}
		if httpsPassword == "" {
			this.jsonError("密码不能为空！")
		}
	}
	if pullType == "ssh" {
		if sshKey == "" {
			this.jsonError("ssh key 不能为空！")
		}
	}
	if branch == "" {
		this.jsonError("代码分支不能为空！")
	}
	if codePath == "" {
		this.jsonError("代码发布路径不能为空！")
	}
	if codeDirUser == "" {
		this.jsonError("目录所属用户不能为空！")
	}
	
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 "+projectId+" 失败: "+err.Error())
		this.jsonError("项目不存在！")
	}
	if len(project) == 0 {
		this.jsonError("项目不存在！")
	}
	
	projectValue := map[string]interface{}{
		"group_id": groupId,
		"repository_url": repositoryUrl,
		"branch": branch,
		"ssh_key": sshKey,
		"ssh_key_salt": sshKeySalt,
		"https_username": httpsUsername,
		"https_password": httpsPassword,
		"code_path": codePath,
		"code_dir_user": codeDirUser,
		"comment": comment,
		"pre_command": "",
		"post_command": "",
		"update_time": time.Now().Unix(),
	}

	_, err = models.ProjectModel.Update(projectId, projectValue)
	if err != nil {
		this.ErrorLog("修改项目 "+projectId+" 失败: "+err.Error())
		this.jsonError("修改项目失败！")
	}else {
		this.InfoLog("修改项目 "+projectId+" 成功")
		this.jsonSuccess("修改项目成功", nil, "/project/list")
	}
}

// 项目详细信息
func (this *ProjectController) Info() {
	
	projectId := this.GetString("project_id", "")
	if projectId == "" {
		this.viewError("项目不存在", "/project/list")
	}
	
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 "+projectId+" 失败: "+err.Error())
		this.viewError("项目不存在", "/project/list")
	}
	if len(project) == 0 {
		this.viewError("项目不存在", "/project/list")
	}
	
	groups, err := models.GroupModel.GetGroups();
	if err != nil {
		this.ErrorLog("查找项目组失败: "+err.Error())
		this.viewError("获取项目组错误", "/project/list")
	}
	
	groupName := ""
	for _, group := range groups {
		if group["group_id"] == project["group_id"] {
			groupName = group["name"]
		}
	}
	
	this.Data["project"] = project
	this.Data["groupName"] = groupName
	this.viewLayoutTitle("项目详细信息", "project/info", "page")
}

// 项目配置
func (this *ProjectController) Config() {
	
	projectId := this.GetString("project_id", "")
	if projectId == "" {
		this.viewError("项目不存在", "/project/list")
	}
	
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 "+projectId+" 失败: "+err.Error())
		this.viewError("项目不存在", "/project/list")
	}
	if len(project) == 0 {
		this.viewError("项目不存在", "/project/list")
	}
	
	this.Data["project"] = project
	this.viewLayoutTitle("项目配置", "project/config", "page")
}

// 项目配置保存
func (this *ProjectController) ConfigSave() {
	
	projectId := this.GetString("project_id", "")
	preCommand := strings.Trim(this.GetString("pre_command", ""), "")
	preCommandExecType := strings.Trim(this.GetString("pre_command_exec_type", ""), "")
	preCommandExecTimeout := strings.Trim(this.GetString("pre_command_exec_timeout", ""), "")
	postCommand := strings.Trim(this.GetString("post_command", ""), "")
	postCommandExecType := strings.Trim(this.GetString("post_command_exec_type", ""), "")
	postCommandExecTimeout := strings.Trim(this.GetString("post_command_exec_timeout", ""), "")
	if projectId == "" {
		this.viewError("项目不存在", "/project/list")
	}
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 "+projectId+" 失败: "+err.Error())
		this.jsonError("项目不存在")
	}
	if len(project) == 0 {
		this.jsonError("项目不存在")
	}
	
	configValue := map[string]interface{}{
		"pre_command": preCommand,
		"pre_command_exec_type": preCommandExecType,
		"pre_command_exec_timeout": preCommandExecTimeout,
		"post_command": postCommand,
		"post_command_exec_type": postCommandExecType,
		"post_command_exec_timeout": postCommandExecTimeout,
		"update_time": time.Now().Unix(),
	}
	
	_, err = models.ProjectModel.Update(projectId, configValue)
	if err != nil {
		this.ErrorLog("项目 "+projectId+" 配置失败: "+err.Error())
		this.jsonError("项目配置失败!")
	}
	
	this.InfoLog("项目 " +projectId+" 配置成功!")
	this.jsonSuccess("项目配置成功", nil, "/project/list")
}

// 项目节点列表
func (this *ProjectController) Node() {

	projectId := this.GetString("project_id", "")
	flag := this.GetString("flag", "update")

	if projectId == "" {
		this.viewError("项目不存在", "/project/list")
	}

	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 "+projectId+" 失败: "+err.Error())
		this.viewError("项目不存在", "/project/list")
	}
	if len(project) == 0 {
		this.viewError("项目不存在", "/project/list")
	}

	// 查找所有的节点组
	nodeGroups, err := models.NodesModel.GetNodeGroups()
	if err != nil {
		this.ErrorLog("查找所有的节点组失败: "+err.Error())
		this.viewError("查找节点出错", "/project/list")
	}
	// 查找所有的节点节点组关系
	nodeNodes, err := models.NodeNodesModel.GetNodeNodes()
	if err != nil {
		this.ErrorLog("查找所有的节点节点组关系失败: "+err.Error())
		this.viewError("查找节点出错", "/project/list")
	}
	//查找所有的节点
	nodes, err := models.NodeModel.GetNodes()
	if err != nil {
		this.ErrorLog("查找所有的节点失败: "+err.Error())
		this.viewError("查找节点出错", "/project/list")
	}

	var projectNodes []map[string]interface{}
	for _, nodeGroup := range nodeGroups {
		projectNode := map[string]interface{}{
			"nodes_id": nodeGroup["nodes_id"],
			"nodes_name": nodeGroup["name"],
			"nodes": []map[string]string{},
		}
		nodeIds := []string{}
		for _, nodeNode := range nodeNodes {
			if nodeGroup["nodes_id"] == nodeNode["nodes_id"] {
				nodeIds = append(nodeIds, nodeNode["node_id"])
			}
		}
		nodeIdsStr := strings.Join(nodeIds, ",")
		nodeGroupNodes := []map[string]string{}
		for _, node := range nodes {
			if strings.Contains(nodeIdsStr, node["node_id"]) {
				nodeValue := map[string]string{
					"node_id": node["node_id"],
					"name": node["name"],
					"ip": node["ip"],
					"port": node["port"],
				}
				nodeGroupNodes = append(nodeGroupNodes, nodeValue)
			}
		}
		projectNode["nodes"] = nodeGroupNodes
		projectNodes = append(projectNodes, projectNode)
	}

	//查找默认的节点
	defaultProjectNodes, _ := models.ProjectNodeModel.GetProjectNodeByProjectId(projectId)
	var defaultNodeIds = []string{}
	for _, defaultProjectNode := range defaultProjectNodes {
		defaultNodeIds = append(defaultNodeIds, defaultProjectNode["node_id"])
	}

	this.Data["flag"] = flag
	this.Data["project"] = project
	this.Data["projectNodes"] = projectNodes
	this.Data["defaultNodeIds"] = strings.Join(defaultNodeIds, ",")
	this.viewLayoutTitle("项目节点", "project/node", "page")
}

// 项目节点保存
func (this *ProjectController) NodeSave() {
	projectId := this.GetString("project_id", "")
	nodeIdsStr := this.GetString("node_ids")
	isCheck := this.GetString("is_check", "")

	if projectId == "" {
		this.jsonError("项目不存在")
	}
	if nodeIdsStr == "" {
		this.jsonError("没有选择节点")
	}

	nodeIds := strings.Split(nodeIdsStr, ",")
	// 先删除
	err := models.ProjectNodeModel.DeleteByProjectIdNodeIds(projectId, nodeIds)
	if err != nil {
		this.ErrorLog("修改项目 "+projectId+" 删除节点"+strings.Join(nodeIds, ",")+" 失败")
		this.jsonError("修改项目节点失败！")
	}
	if isCheck == "1" {
		var insertValues []map[string]interface{}
		for _, nodeId := range nodeIds {
			insertValue := map[string]interface{}{
				"node_id": nodeId,
				"project_id": projectId,
				"create_time": time.Now().Unix(),
			}
			insertValues = append(insertValues, insertValue)
		}
		_, err = models.ProjectNodeModel.InsertBatch(insertValues)
		if err != nil {
			this.ErrorLog("修改项目 "+projectId+" 添加节点"+strings.Join(nodeIds, ",")+" 失败")
			this.jsonError("修改项目节点失败！")
		}
	}

	if isCheck == "1" {
		this.InfoLog("修改项目 "+projectId+" 添加节点"+strings.Join(nodeIds, ",")+" 成功")
	}else {
		this.InfoLog("修改项目 "+projectId+" 删除节点"+strings.Join(nodeIds, ",")+" 成功")
	}

	this.jsonSuccess("修改节点成功", nil)
}

// 删除节点
func (this *ProjectController) Delete() {

	projectId := this.GetString("project_id", "")
	
	if projectId == "" {
		this.jsonError("没有选择项目！")
	}
	
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 "+projectId+" 失败: "+err.Error())
		this.jsonError("项目不存在！")
	}
	if len(project) == 0 {
		this.jsonError("项目不存在！")
	}

	// 删除项目节点关系
	err = models.ProjectNodeModel.DeleteByProjectId(projectId)
	if err != nil {
		this.ErrorLog("删除项目 "+projectId+" 删除项目节点关系失败: "+err.Error())
		this.jsonError("删除项目失败！")
	}

	// 删除项目
	projectValue := map[string]interface{}{
		"is_delete": models.PROJECT_DELETE,
		"update_time": time.Now().Unix(),
	}
	_, err = models.ProjectModel.Update(projectId, projectValue)
	if err != nil {
		this.ErrorLog("删除项目 "+projectId+" 失败: "+err.Error())
		this.jsonError("删除项目失败！")
	}


	this.InfoLog("删除项目 "+projectId+" 成功")
	this.jsonSuccess("删除项目成功", nil, "/project/list")
}