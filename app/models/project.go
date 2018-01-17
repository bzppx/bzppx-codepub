package models

import (
	"bzppx-codepub/app/utils"
	"github.com/snail007/go-activerecord/mysql"
)

const (
	PROJECT_DELETE = 1
	PROJECT_NORMAL = 0
)

const Table_Project_Name = "project"

type Project struct {

}

var ProjectModel = Project{}

// 根据 project_id 获取项目
func (p *Project) GetProjectByProjectId(projectId string) (project map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Project_Name).Where(map[string]interface{}{
		"project_id": projectId,
		"is_delete": PROJECT_NORMAL,
	}))
	if err != nil {
		return
	}
	project = rs.Row()
	return
}

// 根据 project_ids 获取项目
func (p *Project) GetProjectByProjectIds(projectIds []string) (projects []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Project_Name).Where(map[string]interface{}{
		"project_id": projectIds,
		"is_delete": PROJECT_NORMAL,
	}))
	if err != nil {
		return
	}
	projects = rs.Rows()
	return
}

// 项目名称是否存在
func (p *Project) HasSameProjectName(projectId, name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Project_Name).Where(map[string]interface{}{
		"project_id <>": projectId,
		"name":   name,
		"is_delete": PROJECT_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 项目名称是否存在
func (p *Project) HasProjectName(name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Project_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": PROJECT_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 根据 name 获取项目
func (p *Project) GetProjectByName(name string) (project map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Project_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": PROJECT_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	project = rs.Row()
	return
}

// 删除项目
func (p *Project) Delete(projectId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_Project_Name, map[string]interface{}{
		"is_delete": PROJECT_DELETE,
	}, map[string]interface{}{
		"project_id": projectId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入项目
func (p *Project) Insert(project map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Project_Name, project))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改项目
func (p *Project) Update(projectId string, project map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Project_Name, project, map[string]interface{}{
		"project_id": projectId,
		"is_delete": PROJECT_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据关键字分页获取项目
func (p *Project) GetProjectsByKeywordsAndLimit(keywords map[string]string, limit int, number int) (projects []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet

	where := map[string]interface{}{
		"name LIKE": "%" + keywords["keyword"] + "%",
		"is_delete": PROJECT_NORMAL,
	}
	groupId, _ := keywords["group_id"];
	if groupId != "" {
		where["group_id"] = keywords["group_id"]
	}
	sql := db.AR().From(Table_Project_Name).Where(where).Limit(limit, number).OrderBy("project_id", "DESC")

	rs, err = db.Query(sql)
	if err != nil {
		return
	}
	projects = rs.Rows()

	return
}

//分页获取项目
func (project *Project) GetProjectsByLimit(limit int, number int) (projects []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Project_Name).
			Where(map[string]interface{}{
				"is_delete": PROJECT_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("project_id", "DESC"))
	if err != nil {
		return
	}
	projects = rs.Rows()

	return
}

// 项目总数
func (project *Project) CountProjects() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Project_Name).
			Where(map[string]interface{}{
			"is_delete": PROJECT_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取项目总数
func (project *Project) CountProjectsByKeywords(keywords map[string]string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet

	where := map[string]interface{}{
		"name LIKE": "%" + keywords["keyword"] + "%",
		"is_delete": PROJECT_NORMAL,
	}
	groupId, _ := keywords["group_id"];
	if groupId != "" {
		where["group_id"] = keywords["group_id"]
	}
	sql := db.AR().Select("count(*) as total").From(Table_Project_Name).Where(where)

	rs, err = db.Query(sql)
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 获取所有的项目
func (project *Project) GetProjects() (projects []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Project_Name).Where(map[string]interface{}{
		"is_delete": NODE_NORMAL,
	}))
	if err != nil {
		return
	}

	projects = rs.Rows()
	return
}
