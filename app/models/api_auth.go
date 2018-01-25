package models

import (
	"github.com/snail007/go-activerecord/mysql"
)

const Table_ApiAuth_Name = "api_auth"

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
