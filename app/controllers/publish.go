package controllers

import (
	"bzppx-codepub/app/container"
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"strconv"
	"strings"
	"time"
	"bzppx-codepub/app/remotes"
	"sync"
	"fmt"
)

type PublishController struct {
	BaseController
}

// 项目列表
func (this *PublishController) Project() {

	userId := this.UserID
	groupId := this.GetString("group_id", "")
	keyword := strings.Trim(this.GetString("keyword", ""), "")
	keywords := map[string]string{
		"group_id": groupId,
		"keyword":  keyword,
	}

	var err error
	var groups []map[string]string
	var projects []map[string]string

	if this.isAdmin() || this.isRoot() {
		groups, err = models.GroupModel.GetGroups()
		if err != nil {
			this.ErrorLog("查找项目组失败: " + err.Error())
			this.viewError("查找项目出错")
		}
		if (keywords["group_id"] == "") && (len(groups) > 0) {
			keywords["group_id"] = groups[0]["group_id"]
		}
		//if len(keywords) > 0 {
		projects, err = models.ProjectModel.GetProjectsByKeywords(keywords)
		//} else {
		//	projects, err = models.ProjectModel.GetProjects()
		//}
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
		groupIds, err := models.ProjectModel.GetGroupIdsByProjectIds(projectIds)
		if err != nil {
			this.ErrorLog("查找项目失败: " + err.Error())
			this.viewError("查找项目出错")
		}
		groups, err = models.GroupModel.GetGroupsByGroupIds(groupIds)
		if err != nil {
			this.ErrorLog("查找项目组失败: " + err.Error())
			this.viewError("查找项目出错")
		}
		if (keywords["group_id"] == "") && (len(groups) > 0) {
			keywords["group_id"] = groups[0]["group_id"]
		}
		//if len(keywords) > 0 {
		projects, err = models.ProjectModel.GetProjectByProjectIdsAndKeywords(projectIds, keywords)
		if err != nil {
			this.ErrorLog("查找项目失败: " + err.Error())
			this.viewError("查找项目出错")
		}
		//}else {
		//	projects, err = models.ProjectModel.GetProjectByProjectIds(projectIds)
		//	if err != nil {
		//		this.ErrorLog("查找项目失败: " + err.Error())
		//		this.viewError("查找项目出错")
		//	}
		//}
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

	if projectId == "" {
		this.viewError("项目不存在")
	}
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
		this.viewError("封版期间，禁止发布")
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

// 发布历史页面
func (this *PublishController) History() {

	page, _ := this.GetInt("page", 1)
	projectId := this.GetString("project_id", "")

	if projectId == "" {
		this.viewError("项目ID参数出错！")
	}
	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil || len(project) == 0 {
		this.viewError("获取项目信息失败！")
	}

	var count int64
	var tasks []map[string]string
	number := 20
	limit := (page - 1) * number

	count, err = models.TaskModel.CountTaskByProjectId(projectId)
	tasks, err = models.TaskModel.GetTaskByProjectId(projectId, limit, number)
	if err != nil {
		this.viewError("获取任务信息出错！")
	}

	array := utils.NewArray()
	userIds := array.ArrayColumn(tasks, "user_id")
	users, err := models.UserModel.GetUserByUserIdsAndNoLimit(userIds)
	if err != nil {
		this.viewError("获取用户信息出错！")
	}
	usersMap := array.ChangeKey(users, "user_id")

	taskIds := []string{}
	for _, task := range tasks {
		taskIds = append(taskIds, task["task_id"])
		task["status"] = "1" // 状态默认为成功
	}

	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskIds(taskIds)
	if err != nil {
		this.viewError("获取任务信息出错！")
	}
	for index, task := range tasks {
		tasks[index]["username"] = usersMap[task["user_id"]].(map[string]string)["username"]
		for _, taskLog := range taskLogs {
			if (taskLog["task_id"] == task["task_id"]) && (taskLog["status"] != "2")  {
				task["status"] = "2" // 执行中
				break
			}
			if (taskLog["task_id"] == task["task_id"]) && (taskLog["is_success"] == "0")  {
				task["status"] = "0" // 执行失败
				continue
			}
		}
	}

	this.Data["tasks"] = tasks
	this.Data["username"] = this.User["username"]
	this.Data["project"] = project
	this.SetPaginator(number, count)
	this.viewLayoutTitle("历史发布信息", "publish/history", "page")
}

// 节点信息
func (this *PublishController) TaskLog() {
	taskId := this.GetString("task_id", "")
	isReturn := this.GetString("is_return", "1")

	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskId(taskId)
	if err != nil {
		this.viewError("获取节点发布信息失败")
	}
	if len(taskLogs) == 0 {
		this.viewError("节点发布信息不存在")
	}

	nodeIds := []string{}
	for _, taskLog := range taskLogs {
		nodeIds = append(nodeIds, taskLog["node_id"])
		taskLog["node_address"] = ""
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.viewError("获取节点信息失败")
	}
	for _, taskLog := range taskLogs {
		for _, node := range nodes {
			if taskLog["node_id"] == node["node_id"] {
				taskLog["node_address"] = node["ip"]+":"+node["port"]
			}
		}
	}

	this.Data["taskLogs"] = taskLogs
	this.Data["isReturn"] = isReturn
	this.viewLayoutTitle("节点信息", "publish/task-log", "page")
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

	// 判断是否封版
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

	// 查找该项目下的节点
	projectNodes, err := models.ProjectNodeModel.GetProjectNodeByProjectId(projectId)
	if len(projectNodes) <= 0 {
		this.jsonError("该项目下还没有节点！请先添加节点后再发布")
	}
	if err != nil {
		this.ErrorLog("创建任务查询项目 "+projectId+" 节点关系失败：" + err.Error())
		this.jsonError("创建任务失败！")
	}

	// 查找该项目所有的节点信息
	nodeIds := []string{}
	for _, projectNode := range projectNodes {
		nodeIds = append(nodeIds, projectNode["node_id"])
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.ErrorLog("创建任务查找项目 "+projectId+" 节点信息失败：" + err.Error())
		this.jsonError("创建任务失败！")
	}

	// 检查所有的节点是否通畅
	type BadNodes struct {
		Nodes []map[string]string
		Lock sync.Mutex
	}
	badNodes := new(BadNodes)
	var wait sync.WaitGroup
	for _, node := range nodes {
		wait.Add(1)
		go func(node map[string]string) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("%v", err)
					badNodes.Lock.Lock()
					badNodes.Nodes = append(badNodes.Nodes, node)
					badNodes.Lock.Unlock()
				}
				wait.Done()
			}()
			_, err := remotes.System.Ping(node["ip"], node["port"], node["token"], nil)
			if err != nil {
				badNodes.Lock.Lock()
				badNodes.Nodes = append(badNodes.Nodes, node)
				badNodes.Lock.Unlock()
				this.ErrorLog("项目 "+projectId+" 创建任务时节点 "+node["node_id"]+" 检测失败：" + err.Error())
			}
		}(node)
	}
	wait.Wait()
	if len(badNodes.Nodes) > 0 {
		this.jsonError("创建任务失败！有部分节点连接失败")
	}

	// 创建发布任务
	taskId, err := models.TaskModel.Insert(taskValue)
	if err != nil {
		this.ErrorLog("项目 "+projectId+" 创建发布任务失败：" + err.Error())
		this.jsonError("创建任务失败！")
	}

	// 修改项目最后一次发布时间
	projectValue := map[string]interface{}{
		"last_publish_time": taskValue["create_time"],
		"update_time": time.Now().Unix(),
	}
	_, err = models.ProjectModel.Update(projectId, projectValue)
	if err != nil {
		this.ErrorLog("项目 "+projectId+" 创建任务修改最后发布时间失败：" + err.Error())
		this.jsonError("创建任务失败！")
	}

	// 创建节点任务
	taskLog := make([]map[string]interface{}, len(projectNodes))
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
	}
	err = models.TaskLogModel.InsertBatch(taskLog)
	if err != nil {
		this.ErrorLog("项目 "+projectId+" 创建节点任务日志失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}

	// rpc 调用 agent 开始发布
	taskIdStr := strconv.FormatInt(taskId, 10)
	taskLogs, err := models.TaskLogModel.GetTaskLogByTaskId(taskIdStr)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}

	project, err := models.ProjectModel.GetProjectByProjectId(projectId)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}

	// 判断是发布还是回滚
	branch := project["branch"]
	sha1Id, ok := taskValue["sha1_id"].(string)
	if ok {
		branch = sha1Id
	}
	for _, taskLog := range taskLogs {
		ip := ""
		port := ""
		token := ""
		args := map[string]interface{}{
			"task_log_id":  taskLog["task_log_id"],
			"url":          project["repository_url"],
			"ssh_key":      project["ssh_key"],
			"ssh_key_salt": project["ssh_key_salt"],
			"path":         project["code_path"],
			"branch":       branch,
			"username":     project["https_username"],
			"password":     project["https_password"],
			"dir_user":     project["code_dir_user"],
			"pre_command":                  project["pre_command"],
			"pre_command_exec_type":        project["pre_command_exec_type"],
			"pre_command_exec_timeout":     project["pre_command_exec_timeout"],
			"post_command":                 project["post_command"],
			"post_command_exec_type":       project["post_command_exec_type"],
			"post_command_exec_timeout":    project["post_command_exec_timeout"],
		}
		for _, node := range nodes {
			if node["node_id"] == taskLog["node_id"] {
				ip = node["ip"]
				port = node["port"]
				token = node["token"]
				break
			}
		}
		agentMessage := container.AgentMessage{
			Ip:   ip,
			Port: port,
			Token: token,
			Args: args,
		}
		container.Worker.SendPublishChan(agentMessage)
	}

	this.InfoLog("添加发布任务 " + taskIdStr + " 成功")
	this.jsonSuccess("创建发布任务成功！", nil, "/publish/taskLog?task_id="+taskIdStr, 1000)
}
