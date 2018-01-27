package models

import (
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

// 根据 log_id 获取日志
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
