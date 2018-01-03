package models

import "github.com/snail007/go-activerecord/mysql"

const Table_NodeNodes_Name = "node_nodes"

type NodeNodes struct {
}

var NodeNodesModel = NodeNodes{}

// 根据 node_nodes_id 删除关系
func (p *NodeNodes) DeleteNodeNodesByNodeNodesId(nodeNodesId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_NodeNodes_Name, map[string]interface{}{
		"node_nodes_id": nodeNodesId,
	}))
	return
}

// 根据 node_id 删除关系
func (p *NodeNodes) DeleteNodeNodesByNodeId(nodeId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_NodeNodes_Name, map[string]interface{}{
		"node_id": nodeId,
	}))
	return
}

// 根据 node_id 和 nodes_id 删除关系
func (p *NodeNodes) DeleteByNodeIdAndNodesId(nodeId string, nodesId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_NodeNodes_Name, map[string]interface{}{
		"node_id": nodeId,
		"nodes_id": nodesId,
	}))
	return
}

// 批量插入
func (n *NodeNodes) InsertBatch(insertValues []map[string]interface{}) (id int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().InsertBatch(Table_NodeNodes_Name, insertValues))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 批量绑定（先删除后插入）
func (n *NodeNodes) batchBindNodeNodes(insertValues []map[string]interface{}) (id int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().InsertBatch(Table_NodeNodes_Name, insertValues))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 根据 node_id 获取节点节点组关系
func (p *NodeNodes) GetNodeNodesByNodeId(nodeId string) (nodeNodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_NodeNodes_Name).Where(map[string]interface{}{
		"node_id": nodeId,
	}))
	if err != nil {
		return
	}
	nodeNodes = rs.Rows()
	return
}

// 根据 nodes_id 获取节点节点组关系
func (p *NodeNodes) GetNodeNodesByNodesId(nodesId string) (nodeNodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_NodeNodes_Name).Where(map[string]interface{}{
		"nodes_id": nodesId,
	}))
	if err != nil {
		return
	}
	nodeNodes = rs.Rows()
	return
}