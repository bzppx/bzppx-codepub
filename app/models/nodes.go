package models

import (
	"bzppx-codepub/app/utils"
	"github.com/snail007/go-activerecord/mysql"
)

const (
	NODES_DELETE = 1
	NODES_NORMAL = 0
)

const Table_Nodes_Name = "nodes"

type Nodes struct {
}

var NodesModel = Nodes{}

// 根据 nodes_id 获取节点组
func (p *Nodes) GetNodeGroupByNodesId(nodesId string) (nodes map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"nodes_id": nodesId,
		"is_delete": NODES_NORMAL,
	}))
	if err != nil {
		return
	}
	nodes = rs.Row()
	return
}

// 根据多个 nodes_id 获取节点组
func (p *Nodes) GetNodeGroupsByNodesIds(nodesIds []string) (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"nodes_id": nodesIds,
		"is_delete": NODES_NORMAL,
	}))
	if err != nil {
		return
	}
	nodes = rs.Rows()
	return
}

// 节点组名称是否存在
func (p *Nodes) HasSameNodesName(nodesId, name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"nodes_id <>": nodesId,
		"name":   name,
		"is_delete": NODES_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 节点组名称是否存在
func (p *Nodes) HasNodesName(name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": NODES_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 根据 name 获取节点组
func (p *Nodes) GetNodesByName(name string) (nodes map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": NODES_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	nodes = rs.Row()
	return
}

// 删除节点组
func (p *Nodes) Delete(nodesId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_Nodes_Name, map[string]interface{}{
		"is_delete": NODES_DELETE,
	}, map[string]interface{}{
		"nodes_id": nodesId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入节点组
func (p *Nodes) Insert(nodes map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Nodes_Name, nodes))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改节点组
func (p *Nodes) Update(nodesId string, nodes map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Nodes_Name, nodes, map[string]interface{}{
		"nodes_id": nodesId,
		"is_delete": NODES_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据关键字分页获取节点组
func (nodes *Nodes) GetNodeGroupsByKeywordAndLimit(keyword string, limit int, number int) (nodeGroups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"name LIKE": "%" + keyword + "%",
		"is_delete": NODES_NORMAL,
	}).Limit(limit, number).OrderBy("nodes_id", "DESC"))
	if err != nil {
		return
	}
	nodeGroups = rs.Rows()

	return
}

//分页获取节点组
func (nodes *Nodes) GetNodeGroupsByLimit(limit int, number int) (nodeGroups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Nodes_Name).
			Where(map[string]interface{}{
				"is_delete": NODES_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("nodes_id", "DESC"))
	if err != nil {
		return
	}
	nodeGroups = rs.Rows()

	return
}

// 节点组总数
func (nodes *Nodes) CountNodeGroups() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Nodes_Name).
			Where(map[string]interface{}{
			"is_delete": NODES_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取节点组总数
func (nodes *Nodes) CountNodeGroupsByKeyword(keyword string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From(Table_Nodes_Name).
		Where(map[string]interface{}{
			"name LIKE": "%" + keyword + "%",
			"is_delete": NODES_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 获取所有的节点组
func (p *Nodes) GetNodeGroups() (nodes []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Nodes_Name).Where(map[string]interface{}{
		"is_delete": NODES_NORMAL,
	}))
	if err != nil {
		return
	}
	nodes = rs.Rows()
	return
}