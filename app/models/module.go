package models

import (
	"bzppx-codepub/app/utils"
	"github.com/snail007/go-activerecord/mysql"
)

const (
	MODULE_DELETE = 1
	MODULE_NORMAL = 0
)

const Table_Module_Name = "module"

type Module struct {
}

var ModuleModel = Module{}

// 根据 module_id 获取模块
func (p *Module) GetModuleByModuleId(moduleId string) (module map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Module_Name).Where(map[string]interface{}{
		"module_id": moduleId,
		"is_delete": MODULE_NORMAL,
	}))
	if err != nil {
		return
	}
	module = rs.Row()
	return
}

// 模块名称是否存在
func (p *Module) HasSameModuleName(moduleId, name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Module_Name).Where(map[string]interface{}{
		"module_id <>": moduleId,
		"name":   name,
		"is_delete": MODULE_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 模块名称是否存在
func (p *Module) HasModuleName(name string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Module_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": MODULE_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 根据 name 获取模块
func (p *Module) GetModuleByName(name string) (module map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Module_Name).Where(map[string]interface{}{
		"name": name,
		"is_delete": MODULE_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	module = rs.Row()
	return
}

// 删除模块
func (p *Module) Delete(moduleId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_Module_Name, map[string]interface{}{
		"is_delete": MODULE_DELETE,
	}, map[string]interface{}{
		"module_id": moduleId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入模块
func (p *Module) Insert(module map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Module_Name, module))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改模块
func (p *Module) Update(moduleId string, module map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Module_Name, module, map[string]interface{}{
		"module_id": moduleId,
		"is_delete": MODULE_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

//根据关键字分页获取模块
func (module *Module) GetModulesByKeywordsAndLimit(keywords map[string]string, limit int, number int) (modules []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet

	where := map[string]interface{}{
		"name LIKE": "%" + keywords["keyword"] + "%",
		"is_delete": MODULE_NORMAL,
	}
	modulesId, _ := keywords["modules_id"];
	if modulesId != "" {
		where["modules_id"] = keywords["modules_id"]
	}
	sql := db.AR().From(Table_Module_Name).Where(where).Limit(limit, number).OrderBy("module_id", "DESC")

	rs, err = db.Query(sql)
	if err != nil {
		return
	}
	modules = rs.Rows()

	return
}

//分页获取模块
func (module *Module) GetModulesByLimit(limit int, number int) (modules []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Module_Name).
			Where(map[string]interface{}{
				"is_delete": MODULE_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("module_id", "DESC"))
	if err != nil {
		return
	}
	modules = rs.Rows()

	return
}

// 模块总数
func (module *Module) CountModules() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Module_Name).
			Where(map[string]interface{}{
			"is_delete": MODULE_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取模块总数
func (module *Module) CountModulesByKeywords(keywords map[string]string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet

	where := map[string]interface{}{
		"name LIKE": "%" + keywords["keyword"] + "%",
		"is_delete": MODULE_NORMAL,
	}
	modulesId, _ := keywords["modules_id"];
	if modulesId != "" {
		where["modules_id"] = keywords["modules_id"]
	}
	sql := db.AR().Select("count(*) as total").From(Table_Module_Name).Where(where)

	rs, err = db.Query(sql)
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}
