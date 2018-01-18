package controllers

import (
	"bzppx-codepub/app/container"
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"strconv"
	"strings"
	"time"
)

type PublishController struct {
	BaseController
}

// 项目列表
func (this *PublishController) Project() {

	userId := this.UserID
	groupId := this.GetString("group_id", "")
	page, _ := this.GetInt("page", 1)
	keyword := strings.Trim(this.GetString("keyword", ""), "")
	keywords := map[string]string{
		"group_id": groupId,
		"keyword":  keyword,
	}

	var err error
	var groups []map[string]string
	if this.isAdmin() || this.isRoot() {
		groups, err = models.GroupModel.GetGroups()
		if err != nil {
			this.ErrorLog("查找项目组失败: " + err.Error())
			this.viewError("查找项目出错")
		}
	} else {
		userProjects, err := models.UserProjectModel.GetUserProjectByUserId(userId)
		if err != nil {
			this.ErrorLog("查找用户 " + userId + " 项目失败: " + err.Error())
			this.viewError("查找项目出错")
		}
		projectIds := []string{}
		for _, userProject := range userProjects {
			projectIds = append(projectIds, userProject["project_id"])
		}
		projects, err := models.ProjectModel.GetProjectByProjectIds(projectIds)
		if err != nil {
			this.ErrorLog("查找项目失败: " + err.Error())
			this.viewError("查找项目出错")
		}
		groupIds := []string{}
		for _, project := range projects {
			groupIds = append(groupIds, project["group_id"])
		}
		groups, err = models.GroupModel.GetGroupsByGroupIds(groupIds)
		if err != nil {
			this.ErrorLog("查找项目组失败: " + err.Error())
			this.viewError("查找项目出错")
		}
	}

	number := 12
	limit := (page - 1) * number
	var count int64
	var projects []map[string]string
	if len(keywords) > 0 {
		count, err = models.ProjectModel.CountProjectsByKeywords(keywords)
		projects, err = models.ProjectModel.GetProjectsByKeywordsAndLimit(keywords, limit, number)
	} else {
		count, err = models.ProjectModel.CountProjects()
		projects, err = models.ProjectModel.GetProjectsByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("查找用户项目列表失败: " + err.Error())
		this.viewError("查找项目出错")
	}

	//判断是否封版
	var isBlock bool
	block := make(map[string]string)
	if this.isRoot() || this.isAdmin() {
		isBlock = false
	} else {
		isBlock, block, err = models.ConfigureModel.CheckIsBlock()
		if err != nil {
			this.viewError("获取封版配置出错")
		}
	}

	this.Data["isBlock"] = isBlock
	this.Data["block"] = block
	this.Data["projects"] = projects
	this.Data["keywords"] = keywords
	this.Data["groups"] = groups
	this.SetPaginator(number, count)
	this.viewLayoutTitle("项目列表", "publish/project", "page")
}

// 项目信息
func (this *PublishController) Info() {

	projectId := this.GetString("project_id", "")
	if projectId == "" {
		this.viewError("项目不存在", "/publish/project")
	}

	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 " + projectId + " 失败: " + err.Error())
		this.viewError("项目不存在", "/publish/project")
	}
	if len(project) == 0 {
		this.viewError("项目不存在", "/publish/project")
	}
	groups, err := models.GroupModel.GetGroups()
	if err != nil {
		this.ErrorLog("查找项目组失败: " + err.Error())
		this.viewError("获取项目组错误", "/publish/project")
	}
	groupName := ""
	for _, group := range groups {
		if group["group_id"] == project["group_id"] {
			groupName = group["name"]
		}
	}

	// 查找该项目的节点
	projectNodes, err := models.ProjectNodeModel.GetProjectNodeByProjectId(projectId)
	if err != nil {
		this.ErrorLog("查找项目 " + projectId + " 节点关系失败: " + err.Error())
		this.viewError("查找项目信息出错")
	}
	var nodeIds []string
	for _, projectNode := range projectNodes {
		nodeIds = append(nodeIds, projectNode["node_id"])
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.ErrorLog("查找项目失败: " + err.Error())
		this.viewError("查找项目信息出错")
	}

	this.Data["nodes"] = nodes
	this.Data["project"] = project
	this.Data["groupName"] = groupName

	this.viewLayoutTitle("项目详细信息", "publish/info", "page")
}

// 发布页面
func (this *PublishController) Publish() {
	projectId := this.GetString("project_id", "")
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.viewError("查找项目信息出错")
	}

	//判断是否封版
	var isBlock bool
	if this.isRoot() || this.isAdmin() {
		isBlock = false
	} else {
		isBlock, _, err = models.ConfigureModel.CheckIsBlock()
		if err != nil {
			this.viewError("获取封版配置出错")
		}
	}
	if isBlock {
		this.viewError("已封版")
	}

	this.Data["project"] = project
	this.viewLayoutTitle("发布代码", "publish/publish", "page")
}

// 回滚页面
func (this *PublishController) Reset() {
	projectId := this.GetString("project_id", "")
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.viewError("查找项目信息出错")
	}

	//判断是否封版
	var isBlock bool
	if this.isRoot() || this.isAdmin() {
		isBlock = false
	} else {
		isBlock, _, err = models.ConfigureModel.CheckIsBlock()
		if err != nil {
			this.viewError("获取封版配置出错！")
		}
	}
	if isBlock {
		this.viewError("已封版！")
	}

	this.Data["project"] = project
	this.viewLayoutTitle("回滚代码", "publish/reset", "page")
}

// 发布操作
func (this *PublishController) DoPublish() {

	taskValue := make(map[string]interface{}, 4)
	projectId := this.GetString("project_id", "")
	taskValue["project_id"] = projectId
	taskValue["user_id"] = this.UserID
	taskValue["comment"] = this.GetString("comment", "")
	taskValue["create_time"] = utils.NewConvert().IntToString(time.Now().Unix(), 10)
	if taskValue["comment"] == "" {
		this.jsonError("发版说明不能为空！")
	}

	this.addTaskAndTaskLog(taskValue, projectId)
}

// 回滚操作
func (this *PublishController) DoReset() {

	taskValue := make(map[string]interface{}, 4)
	projectId := this.GetString("project_id", "")
	taskValue["sha1_id"] = this.GetString("sha1_id", "")
	taskValue["project_id"] = projectId
	taskValue["user_id"] = this.UserID
	taskValue["comment"] = this.GetString("comment", "")
	taskValue["create_time"] = utils.NewConvert().IntToString(time.Now().Unix(), 10)
	taskValue["publish_time"] = "0"
	if taskValue["comment"] == "" {
		this.jsonError("回滚说明不能为空！")
	}
	if taskValue["sha1_id"] == "" {
		this.jsonError("commit_id不能为空！")
	}

	this.addTaskAndTaskLog(taskValue, projectId)
}

func (this *PublishController) addTaskAndTaskLog(taskValue map[string]interface{}, projectId string) {

	//判断是否封版
	var isBlock bool
	var err error
	if this.isRoot() || this.isAdmin() {
		isBlock = false
	} else {
		isBlock, _, err = models.ConfigureModel.CheckIsBlock()
		if err != nil {
			this.jsonError("获取封版配置出错！")
		}
	}
	if isBlock {
		this.jsonError("已封版！")
	}
	taskId, err := models.TaskModel.Insert(taskValue)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务失败！")
	}

	projectNodes, err := models.ProjectNodeModel.GetProjectNodeByProjectId(projectId)
	if len(projectNodes) <= 0 {
		this.jsonError("该项目下没有节点！")
	}
	if err != nil {
		this.ErrorLog("查询项目节点关系失败：" + err.Error())
		this.jsonError("查询项目节点关系失败！")
	}

	taskLog := make([]map[string]interface{}, len(projectNodes))
	nodeIds := []string{}
	for index, projectNode := range projectNodes {
		taskLog[index] = make(map[string]interface{})
		taskLog[index]["task_id"] = taskId
		taskLog[index]["node_id"] = projectNode["node_id"]
		taskLog[index]["status"] = "0"
		taskLog[index]["is_success"] = "0"
		taskLog[index]["result"] = ""
		taskLog[index]["commit_id"] = ""
		taskLog[index]["create_time"] = time.Now().Unix()
		taskLog[index]["update_time"] = time.Now().Unix()

		nodeIds = append(nodeIds, projectNode["node_id"])
	}

	err = models.TaskLogModel.InsertBatch(taskLog)
	if err != nil {
		this.ErrorLog("创建任务日志失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}

	// rpc 调用 agent 开始发布
	taskIdStr := strconv.FormatInt(taskId, 10)
	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskId(taskIdStr)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}

	for _, taskLog := range taskLogs {
		ip := ""
		port := ""
		args := map[string]interface{}{
			"task_log_id":  taskLog["task_log_id"],
			"url":          project["repository_url"],
			"ssh_key":      project["ssh_key"],
			"ssh_key_salt": project["ssh_key_salt"],
			"path":         project["code_path"],
			"branch":       project["branch"],
			"username":     project["https_username"],
			"password":     project["https_password"],
		}
		for _, node := range nodes {
			if node["node_id"] == taskLog["node_id"] {
				ip = node["ip"]
				port = node["port"]
				break
			}
		}
		agentMessage := container.AgentMessage{
			Ip:   ip,
			Port: port,
			Args: args,
		}
		container.Worker.SendPublishChan(agentMessage)
	}

	this.InfoLog("发布任务成功")
	this.jsonSuccess("操作成功！", nil, "/task/taskLog?task_id="+taskIdStr)
}
