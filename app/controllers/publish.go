package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"strings"
	"time"
)

type PublishController struct {
	BaseController
}

// 模块列表
func (this *PublishController) Module() {

	userId := this.UserID
	modulesId := this.GetString("modules_id", "")
	page, _ := this.GetInt("page", 1)
	keyword := strings.Trim(this.GetString("keyword", ""), "")
	keywords := map[string]string{
		"modules_id": modulesId,
		"keyword":    keyword,
	}

	var err error
	var moduleGroups []map[string]string
	if this.isAdmin() || this.isRoot() {
		moduleGroups, err = models.ModulesModel.GetModuleGroups()
		if err != nil {
			this.ErrorLog("查找模块组失败: " + err.Error())
			this.viewError("查找模块出错")
		}
	} else {
		userModules, err := models.UserModuleModel.GetUserModuleByUserId(userId)
		if err != nil {
			this.ErrorLog("查找用户 " + userId + " 模块失败: " + err.Error())
			this.viewError("查找模块出错")
		}
		moduleIds := []string{}
		for _, userModule := range userModules {
			moduleIds = append(moduleIds, userModule["module_id"])
		}
		modules, err := models.ModuleModel.GetModuleByModuleIds(moduleIds)
		if err != nil {
			this.ErrorLog("查找模块失败: " + err.Error())
			this.viewError("查找模块出错")
		}
		modulesIds := []string{}
		for _, module := range modules {
			modulesIds = append(modulesIds, module["modules_id"])
		}
		moduleGroups, err = models.ModulesModel.GetModuleGroupByModulesIds(modulesIds)
		if err != nil {
			this.ErrorLog("查找用户组失败: " + err.Error())
			this.viewError("查找模块出错")
		}
	}
	if keywords["modules_id"] == "" {
		keywords["modules_id"] = moduleGroups[0]["modules_id"]
	}

	number := 12
	limit := (page - 1) * number
	var count int64
	var modules []map[string]string
	if len(keywords) > 0 {
		count, err = models.ModuleModel.CountModulesByKeywords(keywords)
		modules, err = models.ModuleModel.GetModulesByKeywordsAndLimit(keywords, limit, number)
	} else {
		count, err = models.ModuleModel.CountModules()
		modules, err = models.ModuleModel.GetModulesByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("查找用户模块列表失败: " + err.Error())
		this.viewError("查找模块出错")
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
	this.Data["modules"] = modules
	this.Data["keywords"] = keywords
	this.Data["moduleGroups"] = moduleGroups
	this.SetPaginator(number, count)
	this.viewLayoutTitle("模块列表", "publish/module", "page")
}

// 模块信息
func (this *PublishController) Info() {

	moduleId := this.GetString("module_id", "")
	if moduleId == "" {
		this.viewError("模块不存在", "/publish/module")
	}

	module, err := models.ModuleModel.GetModuleByModuleId(moduleId)
	if err != nil {
		this.ErrorLog("查找模块 " + moduleId + " 失败: " + err.Error())
		this.viewError("模块不存在", "/publish/module")
	}
	if len(module) == 0 {
		this.viewError("模块不存在", "/publish/module")
	}
	moduleGroups, err := models.ModulesModel.GetModuleGroups()
	if err != nil {
		this.ErrorLog("查找模块组失败: " + err.Error())
		this.viewError("获取模块组错误", "/publish/module")
	}
	moduleGroupName := ""
	for _, moduleGroup := range moduleGroups {
		if moduleGroup["modules_id"] == module["modules_id"] {
			moduleGroupName = moduleGroup["name"]
		}
	}

	// 查找该模块的节点
	moduleNodes, err := models.ModuleNodeModel.GetModuleNodeByModuleId(moduleId)
	if err != nil {
		this.ErrorLog("查找模块 " + moduleId + " 节点关系失败: " + err.Error())
		this.viewError("查找模块信息出错")
	}
	var nodeIds []string
	for _, moduleNode := range moduleNodes {
		nodeIds = append(nodeIds, moduleNode["node_id"])
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.ErrorLog("查找模块失败: " + err.Error())
		this.viewError("查找模块信息出错")
	}

	this.Data["nodes"] = nodes
	this.Data["module"] = module
	this.Data["moduleGroupName"] = moduleGroupName

	this.viewLayoutTitle("模块详细信息", "publish/info", "page")
}

// 发布页面
func (this *PublishController) Publish() {
	moduleId := this.GetString("module_id", "")
	module, err := models.ModuleModel.GetModuleByModuleId(moduleId)
	if err != nil {
		this.viewError("查找模块信息出错")
	}

	this.Data["module"] = module
	this.viewLayoutTitle("发布代码", "publish/publish", "page")
}

// 回滚页面
func (this *PublishController) Reset() {
	moduleId := this.GetString("module_id", "")
	module, err := models.ModuleModel.GetModuleByModuleId(moduleId)
	if err != nil {
		this.viewError("查找模块信息出错")
	}

	this.Data["module"] = module
	this.viewLayoutTitle("回滚代码", "publish/reset", "page")
}

// 发布操作
func (this *PublishController) DoPublish() {

	taskValue := make(map[string]interface{}, 4)
	moduleId := this.GetString("module_id", "")
	taskValue["module_id"] = moduleId
	taskValue["user_id"] = this.UserID
	taskValue["comment"] = this.GetString("comment", "")
	taskValue["create_time"] = utils.NewConvert().IntToString(time.Now().Unix(), 10)
	if taskValue["comment"] == "" {
		this.jsonError("发版说明不能为空！")
	}

	this.addTaskAndTaskLog(taskValue, moduleId)
}

// 回滚操作
func (this *PublishController) DoReset() {

	taskValue := make(map[string]interface{}, 4)
	moduleId := this.GetString("module_id", "")
	taskValue["sha1_id"] = this.GetString("sha1_id", "")
	taskValue["module_id"] = moduleId
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

	this.addTaskAndTaskLog(taskValue, moduleId)
}

func (this *PublishController) addTaskAndTaskLog(taskValue map[string]interface{}, moduleId string) {
	taskId, err := models.TaskModel.Insert(taskValue)
	if err != nil {
		this.ErrorLog("创建任务失败：" + err.Error())
		this.jsonError("创建任务失败！")
	}

	moduleNodes, err := models.ModuleNodeModel.GetModuleNodeByModuleId(moduleId)
	if len(moduleNodes) <= 0 {
		this.jsonError("该模块下没有节点！")
	}
	if err != nil {
		this.ErrorLog("查询模块节点关系失败：" + err.Error())
		this.jsonError("查询模块节点关系失败！")
	}

	taskLog := make([]map[string]interface{}, len(moduleNodes))
	for index, moduleNode := range moduleNodes {
		taskLog[index] = make(map[string]interface{})
		taskLog[index]["task_id"] = taskId
		taskLog[index]["node_id"] = moduleNode["node_id"]
		taskLog[index]["status"] = "0"
		taskLog[index]["is_success"] = "0"
		taskLog[index]["result"] = ""
		taskLog[index]["commit_id"] = ""
		taskLog[index]["create_time"] = time.Now().Unix()
		taskLog[index]["update_time"] = time.Now().Unix()
	}

	err = models.TaskLogModel.InsertBatch(taskLog)
	if err != nil {
		this.ErrorLog("创建任务日志失败：" + err.Error())
		this.jsonError("创建任务日志失败！")
	}

	this.jsonSuccess("创建任务成功！", nil, "/task/center")
}
