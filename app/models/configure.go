package models

import (
	"bzppx-codepub/app/utils"
	"time"

	"github.com/snail007/go-activerecord/mysql"
)

const (
	CONFIGURE_DELETE = 1
	CONFIGURE_NORMAL = 0
)

const Table_Configure_Name = "configure"

type Configure struct {
}

var ConfigureModel = Configure{}

// 获取封版信息
func (config *Configure) GetBlock() (block map[string]string, err error) {
	db := G.DB()
	keys := []string{"block_message", "block_is_enable", "block_start_time", "block_end_time"}
	block = make(map[string]string)
	var rs *mysql.ResultSet

	rs, err = db.Query(db.AR().From(Table_Configure_Name).Where(map[string]interface{}{
		"key":       keys,
		"is_delete": CONFIGURE_NORMAL,
	}))

	if err != nil {
		return
	}
	data := rs.Rows()
	for _, v := range data {
		block[v["key"]] = v["value"]
	}
	return
}

// 修改封版信息
func (config *Configure) InsertBlock(blockValue []map[string]interface{}) (err error) {
	db := G.DB()
	where := []string{"key", "key", "key", "key"}
	_, err = db.Exec(db.AR().Where(map[string]interface{}{
		"is_delete": CONFIGURE_NORMAL,
	}).UpdateBatch(Table_Configure_Name, blockValue, where))

	return
}

func (config *Configure) GetEmail() (email map[string]string, err error) {
	db := G.DB()
	keys := []string{"email_host", "email_port", "email_username", "email_password", "email_is_ssl", "email_from", "email_cc_list"}
	email = make(map[string]string)
	var rs *mysql.ResultSet

	rs, err = db.Query(db.AR().From(Table_Configure_Name).Where(map[string]interface{}{
		"key":       keys,
		"is_delete": CONFIGURE_NORMAL,
	}))
	if err != nil {
		return
	}
	data := rs.Rows()

	for _, v := range data {
		email[v["key"]] = v["value"]
	}
	return
}

func (config *Configure) InsertEmailConfig(emailValue []map[string]interface{}) (err error) {
	db := G.DB()
	where := []string{"key", "key", "key", "key", "key", "key", "key"}
	_, err = db.Exec(db.AR().Where(map[string]interface{}{
		"is_delete": CONFIGURE_NORMAL,
	}).UpdateBatch(Table_Configure_Name, emailValue, where))

	return
}

func (config *Configure) CheckIsBlock() (isBlock bool, block map[string]string, err error) {

	block, err = config.GetBlock()
	if err != nil {
		return
	}
	if block["block_is_enable"] == "1" {
		isBlock = false
		return
	}
	blockStartTime := utils.NewConvert().StringToInt64(block["block_start_time"])
	blockEndTime := utils.NewConvert().StringToInt64(block["block_end_time"])
	now := time.Now().Unix()
	if now > blockEndTime || now < blockStartTime {
		isBlock = false
	} else {
		isBlock = true
	}
	return
}
