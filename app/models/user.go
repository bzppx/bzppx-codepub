package models

import (
	"bzppx-codepub/app/utils"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/snail007/go-activerecord/mysql"
)

const (
	USER_ROLE_ROOT  = 3
	USER_ROLE_ADMIN = 2
	USER_ROLE_USER  = 1

	USER_DELETE = 1
	USER_NORMAL = 0
)

const Table_User_Name = "user"

type User struct {
}

var UserModel = User{}

// 根据 user_id 获取用户
func (p *User) GetUserByUserId(userId string) (user map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"user_id":   userId,
		"is_delete": USER_NORMAL,
	}))
	if err != nil {
		return
	}
	user = rs.Row()
	return
}

// 用户名是否存在
func (p *User) HasSameUsername(userId, username string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"user_id <>": userId,
		"username":   username,
		"is_delete":  USER_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 用户名是否存在
func (p *User) HasUsername(username string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"username":  username,
		"is_delete": USER_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}

// 根据用户名查找用户
func (p *User) GetUserByName(username string) (user map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"username":  username,
		"is_delete": USER_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	user = rs.Row()
	return
}

// 删除
func (p *User) Delete(userId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_User_Name, map[string]interface{}{
		"is_delete": USER_DELETE,
	}, map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	return
}

//恢复
func (p *User) Review(userId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_User_Name, map[string]interface{}{
		"is_delete": USER_NORMAL,
	}, map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入
func (p *User) Insert(user map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_User_Name, user))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改
func (p *User) Update(userId string, user map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_User_Name, user, map[string]interface{}{
		"user_id":   userId,
		"is_delete": USER_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改密码
func (p *User) ChangePassword(userId, newpassword, oldpassword string) (err error) {
	db := G.DB()
	user, err := p.GetUserByUserId(userId)
	if user["password"] != p.EncodePassword(oldpassword) {
		return errors.New("旧密码错误")
	}
	if err != nil {
		return
	}
	_, err = db.Exec(db.AR().Update(Table_User_Name, map[string]interface{}{
		"password": p.EncodePassword(newpassword),
	}, map[string]interface{}{
		"user_id":   userId,
		"is_delete": USER_NORMAL,
	}))
	if err != nil {
		return
	}
	return
}

// 加密password
func (p *User) EncodePassword(password string) (passwordHash string) {
	hasher := md5.New()
	hasher.Write([]byte(password))
	passwordHash = strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
	return
}

//根据关键字分页获取用户
func (user *User) GetUsersByKeywordAndLimit(keyword string, limit int, number int) (users []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"username LIKE": "%" + keyword + "%",
		"is_delete":     USER_NORMAL,
	}).Limit(limit, number).OrderBy("user_id", "DESC"))
	if err != nil {
		return
	}
	users = rs.Rows()

	return
}

//分页获取用户
func (user *User) GetUsersByLimit(limit int, number int) (users []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_User_Name).
			Where(map[string]interface{}{
				"is_delete": USER_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("user_id", "DESC"))
	if err != nil {
		return
	}
	users = rs.Rows()

	return
}

// 获取用户总数
func (user *User) CountUsers() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_User_Name).
			Where(map[string]interface{}{
				"is_delete": USER_NORMAL,
			}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取用户总数
func (user *User) CountUsersByKeyword(keyword string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From(Table_User_Name).
		Where(map[string]interface{}{
			"username LIKE": "%" + keyword + "%",
			"is_delete":     USER_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据用户名模糊查找用户
func (p *User) GetUserByLikeName(username string) (user []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"username Like": "%" + username + "%",
		"is_delete":     USER_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	user = rs.Rows()
	return
}

// 根据 user_ids 获取用户
func (p *User) GetUserByUserIds(userIds []string) (users []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_User_Name).Where(map[string]interface{}{
		"user_id":   userIds,
		"is_delete": USER_NORMAL,
	}))
	if err != nil {
		return
	}
	users = rs.Rows()
	return
}
