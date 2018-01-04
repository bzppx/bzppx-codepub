package controllers

import (
	"bzppx-codepub/app/models"
	"strings"
)

type PublishController struct {
	BaseController
}

// 模块列表
func (this *PublishController) Module() {

	userId := this.UserID
	modulesId := this.GetString("modules_id", "")
	page, _:= this.GetInt("page", 1)
	keyword := strings.Trim(this.GetString("keyword", ""), "")
	keywords := map[string]string {
		"modules_id": modulesId,
		"keyword": keyword,
	}

	var err error
	var moduleGroups []map[string]string
	if (this.isAdmin() || this.isRoot()) {
		moduleGroups, err = models.ModulesModel.GetModuleGroups()
		if err != nil {
			this.viewError("查找模块出错")
		}
	}else {
		userModules, err := models.UserModuleModel.GetUserModuleByUserId(userId)
		if err != nil {
			this.viewError("查找模块出错")
		}
		moduleIds := []string{}
		for _, userModule := range userModules {
			moduleIds = append(moduleIds, userModule["module_id"])
		}
		modules, err := models.ModuleModel.GetModuleByModuleIds(moduleIds)
		if err != nil {
			this.viewError("查找模块出错")
		}
		modulesIds := []string{}
		for _, module := range modules {
			modulesIds = append(modulesIds, module["modules_id"])
		}
		moduleGroups, err = models.ModulesModel.GetModuleGroupByModulesIds(modulesIds)
		if err != nil {
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
	if (len(keywords) > 0) {
		count, err = models.ModuleModel.CountModulesByKeywords(keywords)
		modules, err = models.ModuleModel.GetModulesByKeywordsAndLimit(keywords, limit, number)
	}else {
		count, err = models.ModuleModel.CountModules()
		modules, err = models.ModuleModel.GetModulesByLimit(limit, number)
	}
	if err != nil {
		this.viewError("查找模块出错")
	}

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
		this.viewError("模块不存在", "/publish/module")
	}
	if len(module) == 0 {
		this.viewError("模块不存在", "/publish/module")
	}
	moduleGroups, err := models.ModulesModel.GetModuleGroups();
	if err != nil {
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
		this.viewError("查找模块信息出错")
	}
	var nodeIds []string
	for _, moduleNode := range moduleNodes {
		nodeIds = append(nodeIds, moduleNode["node_id"])
	}
	nodes, err := models.NodeModel.GetNodeByNodeIds(nodeIds)
	if err != nil {
		this.viewError("查找模块信息出错")
	}

	this.Data["nodes"] = nodes
	this.Data["module"] = module
	this.Data["moduleGroupName"] = moduleGroupName

	this.viewLayoutTitle("模块详细信息", "publish/info", "page")
}

