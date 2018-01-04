package models

import "github.com/snail007/go-activerecord/mysql"

const Table_UserModule_Name = "user_module"

type UserModule struct {
}

var UserModuleModel = UserModule{}

func (p *UserModule) DeleteUserModuleByUserId(userId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_UserModule_Name, map[string]interface{}{
		"user_id": userId,
	}))
	return
}

func (p *UserModule) DeleteByUserIdModuleIds(userId string, moduleIds []string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_UserModule_Name, map[string]interface{}{
		"user_id": userId,
		"module_id": moduleIds,
	}))
	return
}

func (p *UserModule) DeleteByModuleId(moduleId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Delete(Table_UserModule_Name, map[string]interface{}{
		"module_id": moduleId,
	}))
	return
}

// 批量插入
func (n *UserModule) InsertBatch(insertValues []map[string]interface{}) (id int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().InsertBatch(Table_UserModule_Name, insertValues))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *UserModule) GetUserModuleByModuleId(moduleId string) (userModules []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_UserModule_Name).Where(map[string]interface{}{
		"module_id": moduleId,
	}))
	if err != nil {
		return
	}
	userModules = rs.Rows()
	return
}

func (p *UserModule) GetUserModuleByUserId(userId string) (userModules []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_UserModule_Name).Where(map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	userModules = rs.Rows()
	return
}