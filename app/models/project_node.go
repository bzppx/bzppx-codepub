package models

import "github.com/snail007/go-activerecord/mysql"

const Table_ProjectNode_Name = "project_node"

type ProjectNode struct {
}

var ProjectNodeModel = ProjectNode{}

func (p *ProjectNode) DeleteProjectNodeByNodeId(nodeId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_ProjectNode_Name, map[string]interface{}{
		"node_id": nodeId,
	}))
	return
}

func (p *ProjectNode) DeleteByProjectIdNodeIds(projectId string, nodeIds []string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_ProjectNode_Name, map[string]interface{}{
		"node_id": nodeIds,
		"project_id": projectId,
	}))
	return
}

func (p *ProjectNode) DeleteByProjectId(projectId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_ProjectNode_Name, map[string]interface{}{
		"project_id": projectId,
	}))
	return
}

// 批量插入
func (n *ProjectNode) InsertBatch(insertValues []map[string]interface{}) (id int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().InsertBatch(Table_ProjectNode_Name, insertValues))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *ProjectNode) GetProjectNodeByProjectId(projectId string) (projectNodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ProjectNode_Name).Where(map[string]interface{}{
		"project_id": projectId,
	}))
	if err != nil {
		return
	}
	projectNodes = rs.Rows()
	return
}