package models

import (
	"bzppx-codepub/app/utils"
	"fmt"

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
	fmt.Println(rs)
	nodes = rs.Rows()
	fmt.Println(nodes)
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
		"ip":        keywords["ip"],
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
		"ip":        keywords["ip"],
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
