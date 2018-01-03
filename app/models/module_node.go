package models

import "github.com/snail007/go-activerecord/mysql"

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

func (p *ModuleNode) DeleteByModuleIdNodeIds(moduleId string, nodeIds []string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_ModuleNode_Name, map[string]interface{}{
		"node_id": nodeIds,
		"module_id": moduleId,
	}))
	return
}

func (p *ModuleNode) DeleteByModuleId(moduleId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_ModuleNode_Name, map[string]interface{}{
		"module_id": moduleId,
	}))
	return
}

// 批量插入
func (n *ModuleNode) InsertBatch(insertValues []map[string]interface{}) (id int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().InsertBatch(Table_ModuleNode_Name, insertValues))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *ModuleNode) GetModuleNodeByModuleId(moduleId string) (moduleNodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ModuleNode_Name).Where(map[string]interface{}{
		"module_id": moduleId,
	}))
	if err != nil {
		return
	}
	moduleNodes = rs.Rows()
	return
}