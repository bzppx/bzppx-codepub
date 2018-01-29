package models

import (
	"bzppx-codepub/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

const (
	Table_ApiAuth_Name = "api_auth"
	API_AUTH_DELETE    = "1"
	API_AUTH_NORMAL    = "0"
)

type ApiAuth struct {
}

var ApiAuthModel = ApiAuth{}

// 获取所有的api_auth
func (l *ApiAuth) GetAllShowApiAuth() (data []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ApiAuth_Name).
		OrderBy("sort", "asc").
		Where(map[string]interface{}{
			"is_show":   1,
			"is_delete": 0,
		}))
	if err != nil {
		return
	}
	data = rs.Rows()
	return
}

// 根据 api_auth_id 获取 api_auth
func (l *ApiAuth) GetApiAuthByApiAuthId(apiAuthId string) (data map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ApiAuth_Name).Where(map[string]interface{}{
		"api_auth_id": apiAuthId,
		"is_delete":   API_AUTH_NORMAL,
	}))
	if err != nil {
		return
	}
	data = rs.Row()
	return
}

// 分页获取api_auth
func (l *ApiAuth) GetApiAuthByLimit(number, limit int) (authApis []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ApiAuth_Name).Where(map[string]interface{}{
		"is_delete": API_AUTH_NORMAL,
	}).Limit(limit, number))
	if err != nil {
		return
	}

	authApis = rs.Rows()
	return
}

// 获取api_auth数量
func (l *ApiAuth) CountApiAuth() (count int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ApiAuth_Name).Select("count(*) as total"))
	if err != nil {
		return
	}

	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 删除 api_auth
func (l *ApiAuth) DeleteByAuthApiId(apiAuthId string) (affect int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_ApiAuth_Name, map[string]interface{}{
		"is_delete": API_AUTH_DELETE,
	}, map[string]interface{}{
		"api_auth_id": apiAuthId,
	}))

	if err != nil {
		return
	}
	affect = rs.RowsAffected
	return
}

// 判断 key 是否存在
func (l *ApiAuth) CheckKeyExist(key, apiAuthId string) (exist bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_ApiAuth_Name).Where(map[string]interface{}{
		"is_delete": API_AUTH_NORMAL,
		"key":       key,
	}).Limit(0, 1))
	if err != nil {
		return
	}

	if len(rs.Row()) > 0 {
		if rs.Value("api_auth_id") != apiAuthId {
			exist = true
		}
	}
	return
}

// 添加 auth_api
func (l *ApiAuth) Insert(apiAuth map[string]interface{}) (apiAuthId int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_ApiAuth_Name, apiAuth))
	if err != nil {
		return
	}

	apiAuthId = rs.LastInsertId
	return
}

// 修改 auth_api
func (l *ApiAuth) UpdateByAuthApiId(apiAuth map[string]interface{}, apiAuthId string) (affect int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_ApiAuth_Name, apiAuth, map[string]interface{}{
		"api_auth_id": apiAuthId,
	}))
	if err != nil {
		return
	}

	affect = rs.RowsAffected
	return
}
