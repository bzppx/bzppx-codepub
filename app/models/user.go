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
)

type User struct {
}

var UserModel = User{}

func (p *User) GetUserByUserId(userId string) (user map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("user").Where(map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	user = rs.Row()
	return
}

func (p *User) HasSameUsername(userId, username string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("user").Where(map[string]interface{}{
		"user_id <>": userId,
		"username":   username,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *User) HasUsername(username string) (has bool, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("user").Where(map[string]interface{}{
		"username": username,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	if rs.Len() > 0 {
		has = true
	}
	return
}
func (p *User) GetUserByName(username string) (user map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From("user").Where(map[string]interface{}{
		"username": username,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	user = rs.Row()
	return
}

//禁用
func (p *User) Forbidden(userId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update("user", map[string]interface{}{
		"is_forbidden": 1,
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
	_, err = db.Exec(db.AR().Update("user", map[string]interface{}{
		"is_forbidden": 0,
	}, map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	return
}

func (p *User) Insert(user map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert("user", user))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *User) Update(userId string, user map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update("user", user, map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

func (p *User) ChangePassword(userId, newpassword, oldpassword string) (err error) {
	db := G.DB()
	user, err := p.GetUserByUserId(userId)
	if user["password"] != p.EncodePassword(oldpassword) {
		return errors.New("旧密码错误")
	}
	if err != nil {
		return
	}
	_, err = db.Exec(db.AR().Update("user", map[string]interface{}{
		"password": p.EncodePassword(newpassword),
	}, map[string]interface{}{
		"user_id": userId,
	}))
	if err != nil {
		return
	}
	return
}
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
	rs, err = db.Query(db.AR().From("user").Where(map[string]interface{}{
		"username LIKE": "%" + keyword + "%",
	}).Limit(limit, number))
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
	rs, err = db.Query(db.AR().From("user").Limit(limit, number))
	if err != nil {
		return
	}
	users = rs.Rows()

	return
}

func (user *User) CountUsers() (count int, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From("user"))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}

func (user *User) CountUsersByKeyword(keyword string) (count int, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From("user").
		Where(map[string]interface{}{
			"username LIKE": "%" + keyword + "%",
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt(rs.Value("total"))
	return
}
