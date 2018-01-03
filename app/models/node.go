package models

import (
	"bzppx-codepub/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

const (
	NODE_DELETE = 1
	NODE_NORMAL = 0
)

const Table_Node_Name = "node"

type Node struct {
}

var NodeModel = Node{}

//分页获取节点
func (node *Node) GetNodesByLimit(limit int, number int) (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Node_Name).
			Where(map[string]interface{}{
				"is_delete": NODE_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("node_id", "DESC"))
	if err != nil {
		return
	}
	nodes = rs.Rows()
	return
}

// 节点总数
func (node *Node) CountNodes() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Node_Name).
			Where(map[string]interface{}{
				"is_delete": NODE_NORMAL,
			}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

//根据关键字分页获取节点
func (node *Node) GetNodesByKeywordsAndLimit(keywords map[string]string, limit int, number int) (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet

	where := map[string]interface{}{
		"ip LIKE": "%" + keywords["ip"] + "%",
		"is_delete": NODE_NORMAL,
	}

	sql := db.AR().From(Table_Node_Name).Where(where).Limit(limit, number).OrderBy("node_id", "DESC")

	rs, err = db.Query(sql)
	if err != nil {
		return
	}
	nodes = rs.Rows()

	return
}

// 根据关键字获取节点总数
func (node *Node) CountNodesByKeywords(keywords map[string]string) (count int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet

	where := map[string]interface{}{
		"ip LIKE": "%" + keywords["ip"] + "%",
		"is_delete": NODE_NORMAL,
	}

	sql := db.AR().Select("count(*) as total").From(Table_Node_Name).Where(where)

	rs, err = db.Query(sql)
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 插入节点
func (node *Node) Insert(nodeValue map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Node_Name, nodeValue))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改节点
func (node *Node) Update(nodeId string, nodeValue map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Node_Name, nodeValue, map[string]interface{}{
		"node_id":   nodeId,
		"is_delete": NODE_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 是否存在节点
func (node *Node) HasNodeByIpAndPort(nodeId, ip string, port int) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Node_Name).Where(map[string]interface{}{
		"ip":        ip,
		"port":      port,
		"is_delete": NODES_NORMAL,
	}))
	if err != nil {
		return
	}

	has = false
	nodes := rs.Rows()
	for _, node := range nodes {
		if node["node_id"] != nodeId {
			has = true
			return
		}
	}

	return
}

// 通过nodeId获取节点数据
func (node *Node) GetNodeByNodeId(nodeId string) (nodes map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Node_Name).Where(map[string]interface{}{
		"node_id":   nodeId,
		"is_delete": NODE_NORMAL,
	}))
	if err != nil {
		return
	}

	nodes = rs.Row()
	return
}

// 通过多个 node_id 获取节点数据
func (node *Node) GetNodeByNodeIds(nodeIds []string) (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Node_Name).Where(map[string]interface{}{
		"node_id":   nodeIds,
		"is_delete": NODE_NORMAL,
	}))
	if err != nil {
		return
	}

	nodes = rs.Rows()
	return
}

// 除 node_ids 的节点数据
func (node *Node) GetNodeByNotNodeIds(nodeIds []string) (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Node_Name).Where(map[string]interface{}{
		"node_id NOT":   nodeIds,
		"is_delete": NODE_NORMAL,
	}))
	if err != nil {
		return
	}

	nodes = rs.Rows()
	return
}

// 获取所有的节点
func (node *Node) GetNodes() (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Node_Name).Where(map[string]interface{}{
		"is_delete": NODE_NORMAL,
	}))
	if err != nil {
		return
	}

	nodes = rs.Rows()
	return
}