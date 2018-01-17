package models

import (
	"bzppx-codepub/app/utils"
	"github.com/snail007/go-activerecord/mysql"
)

const (
	GROUP_DELETE = 1
	GROUP_NORMAL = 0
)

const Table_Group_Name = "group"

type Group struct {
}

var GroupModel = Group{}

// 根据 group_id 获取项目组
func (g *Group) GetGroupByGroupId(groupId string) (groups map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Group_Name).Where(map[string]interface{}{
		"group_id": groupId,
		"is_delete": GROUP_NORMAL,
	}))
	if err != nil {
		return
	}
	groups = rs.Row()
	return
}

// 根据 group_ids 获取项目组
func (g *Group) GetGroupsByGroupIds(groupIds []string) (groups []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Group_Name).Where(map[string]interface{}{
		"group_id": groupIds,
		"is_delete": GROUP_NORMAL,
	}))
	if err != nil {
		return
	}
	groups = rs.Rows()
	return
}

// 项目组名称是否存在
func (g *Group) HasSameGroupName(groupId, name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Group_Name).Where(map[string]interface{}{
		"group_id <>": groupId,
		"name":   name,
		"is_delete": GROUP_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 项目组名称是否存在
func (g *Group) HasGroupName(name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Group_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": GROUP_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 根据 name 获取项目组
func (g *Group) GetGroupByName(name string) (group map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Group_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": GROUP_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	group = rs.Row()
	return
}

// 删除项目组
func (g *Group) Delete(groupId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_Group_Name, map[string]interface{}{
		"is_delete": GROUP_DELETE,
	}, map[string]interface{}{
		"group_id": groupId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入项目组
func (g *Group) Insert(group map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Group_Name, group))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改项目组
func (g *Group) Update(groupId string, group map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Group_Name, group, map[string]interface{}{
		"group_id": groupId,
		"is_delete": GROUP_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据关键字分页获取项目组
func (g *Group) GetGroupsByKeywordAndLimit(keyword string, limit int, number int) (groups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Group_Name).Where(map[string]interface{}{
		"name LIKE": "%" + keyword + "%",
		"is_delete": GROUP_NORMAL,
	}).Limit(limit, number).OrderBy("group_id", "DESC"))
	if err != nil {
		return
	}
	groups = rs.Rows()

	return
}

//分页获取项目组
func (g *Group) GetGroupsByLimit(limit int, number int) (groups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Group_Name).
			Where(map[string]interface{}{
				"is_delete": GROUP_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("group_id", "DESC"))
	if err != nil {
		return
	}
	groups = rs.Rows()

	return
}

// 获取所有的项目组
func (g *Group) GetGroups() (groups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Group_Name).
			Where(map[string]interface{}{
			"is_delete": GROUP_NORMAL,
		}))
	if err != nil {
		return
	}
	groups = rs.Rows()

	return
}

// 项目组总数
func (g *Group) CountGroups() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Group_Name).
			Where(map[string]interface{}{
			"is_delete": GROUP_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取项目组总数
func (g *Group) CountGroupsByKeyword(keyword string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From(Table_Group_Name).
		Where(map[string]interface{}{
			"name LIKE": "%" + keyword + "%",
			"is_delete": GROUP_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}
