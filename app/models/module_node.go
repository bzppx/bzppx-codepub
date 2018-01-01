package models

const Table_ModuleNode_Name = "module_node"

type ModuleNode struct {
}

var ModuleNodeModel = ModuleNode{}

func (p *ModuleNode) DeleteModuleNodeByNodeId(NodeId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_ModuleNode_Name, map[string]interface{}{
		"node_id": NodeId,
	}))
	return
}
