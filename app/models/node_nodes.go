package models

const Table_NodeNodes_Name = "node_nodes"

type NodeNodes struct {
}

var NodeNodesModel = NodeNodes{}

func (p *NodeNodes) DeleteNodeNodesByNodeId(NodeNodesId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_NodeNodes_Name, map[string]interface{}{
		"node_nodes_id": NodeNodesId,
	}))
	return
}
