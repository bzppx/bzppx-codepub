package models

import (
	"github.com/snail007/go-activerecord/mysql"
	"strconv"
)

const Table_UserProject_Name = "user_project"

type UserProject struct {
}

var UserProjectModel = UserProject{}

func (p *UserProject) DeleteUserProjectByUserId(userId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_UserProject_Name, map[string]interface{}{
		"user_id": userId,
	}))
	return
}

func (p *UserProject) DeleteByUserIdProjectIds(userId string, projectIds []string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_UserProject_Name, map[string]interface{}{
		"user_id": userId,
		"project_id": projectIds,
	}))
	return
}

func (p *UserProject) DeleteByProjectId(projectId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_UserProject_Name, map[string]interface{}{
		"project_id": projectId,
	}))
	return
}

// 批量插入
func (n *UserProject) InsertBatch(insertValues []map[string]interface{}) (id int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().InsertBatch(Table_UserProject_Name, insertValues))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *UserProject) GetUserProjectByProjectId(projectId string) (userProjects []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_UserProject_Name).Where(map[string]interface{}{
		"project_id": projectId,
	}))
	if err != nil {
		return
	}
	userProjects = rs.Rows()
	return
}

func (p *UserProject) GetUserProjectByUserId(userId string) (userProjects []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_UserProject_Name).Where(map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	userProjects = rs.Rows()
	return
}

func (p *UserProject) CountProjectByUserId(userId string) (total int, err error) {
	db := G.DB()
	var rs *mysql.ResultSet

	sql := db.AR().From(Table_UserProject_Name).Select("count(*) as total").Where(map[string]interface{}{
		"user_id": userId,
	})

	rs, err = db.Query(sql)
	if err != nil {
		return
	}

	if rs.Value("total") != "" {
		total, _ = strconv.Atoi(rs.Value("total"))
	}
	return total, nil
}