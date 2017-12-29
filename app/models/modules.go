package models

import (
	"bzppx-codepub/app/utils"
	"github.com/snail007/go-activerecord/mysql"
)

const (
	MODULES_DELETE = 1
	MODULES_NORMAL = 0
)

const Table_Modules_Name = "modules"

type Modules struct {
}

var ModulesModel = Modules{}

// 根据 modules_id 获取模块组
func (p *Modules) GetModuleGroupByModulesId(modulesId string) (modules map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Modules_Name).Where(map[string]interface{}{
		"modules_id": modulesId,
		"is_delete": MODULES_NORMAL,
	}))
	if err != nil {
		return
	}
	modules = rs.Row()
	return
}

// 模块组名称是否存在
func (p *Modules) HasSameModulesName(modulesId, name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Modules_Name).Where(map[string]interface{}{
		"modules_id <>": modulesId,
		"name":   name,
		"is_delete": MODULES_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 模块组名称是否存在
func (p *Modules) HasModulesName(name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Modules_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": MODULES_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 根据 name 获取模块组
func (p *Modules) GetModulesByName(name string) (modules map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Modules_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": MODULES_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	modules = rs.Row()
	return
}

// 删除模块组
func (p *Modules) Delete(modulesId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_Modules_Name, map[string]interface{}{
		"is_delete": MODULES_DELETE,
	}, map[string]interface{}{
		"modules_id": modulesId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入模块组
func (p *Modules) Insert(modules map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Modules_Name, modules))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改模块组
func (p *Modules) Update(modulesId string, modules map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Modules_Name, modules, map[string]interface{}{
		"modules_id": modulesId,
		"is_delete": MODULES_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据关键字分页获取模块组
func (modules *Modules) GetModuleGroupsByKeywordAndLimit(keyword string, limit int, number int) (moduleGroups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Modules_Name).Where(map[string]interface{}{
		"name LIKE": "%" + keyword + "%",
		"is_delete": MODULES_NORMAL,
	}).Limit(limit, number).OrderBy("modules_id", "DESC"))
	if err != nil {
		return
	}
	moduleGroups = rs.Rows()

	return
}

//分页获取模块组
func (modules *Modules) GetModuleGroupsByLimit(limit int, number int) (moduleGroups []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Modules_Name).
			Where(map[string]interface{}{
				"is_delete": MODULES_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("modules_id", "DESC"))
	if err != nil {
		return
	}
	moduleGroups = rs.Rows()

	return
}

// 模块组总数
func (modules *Modules) CountModuleGroups() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Modules_Name).
			Where(map[string]interface{}{
			"is_delete": MODULES_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取模块组总数
func (modules *Modules) CountModuleGroupsByKeyword(keyword string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From(Table_Modules_Name).
		Where(map[string]interface{}{
			"name LIKE": "%" + keyword + "%",
			"is_delete": MODULES_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}
