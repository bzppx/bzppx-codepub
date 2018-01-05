package models

import (
	"bzppx-codepub/app/utils"
	"github.com/snail007/go-activerecord/mysql"
)

const (
	NOTICE_DELETE = 1
	NOTICE_NORMAL = 0
)

const Table_Notice_Name = "notice"

type Notice struct {
}

var NoticeModel = Notice{}

// 根据 notice_id 获取公告
func (p *Notice) GetNoticeByNoticeId(noticeId string) (notice map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Notice_Name).Where(map[string]interface{}{
		"notice_id": noticeId,
		"is_delete": NOTICE_NORMAL,
	}))
	if err != nil {
		return
	}
	notice = rs.Row()
	return
}

// 根据公告标题查找公告
func (p *Notice) GetNoticeByTitle(title string) (notice map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Notice_Name).Where(map[string]interface{}{
		"title": title,
		"is_delete": NOTICE_NORMAL,
	}).Limit(0, 1))
	if err != nil {
		return
	}
	notice = rs.Row()
	return
}

// 删除
func (p *Notice) Delete(noticeId string) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().Update(Table_Notice_Name, map[string]interface{}{
		"is_delete": NOTICE_DELETE,
	}, map[string]interface{}{
		"notice_id": noticeId,
	}))
	if err != nil {
		return
	}
	return
}

// 插入
func (p *Notice) Insert(notice map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Notice_Name, notice))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 修改
func (p *Notice) Update(noticeId string, notice map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_Notice_Name, notice, map[string]interface{}{
		"notice_id": noticeId,
		"is_delete": NOTICE_NORMAL,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}


//根据关键字分页获取公告
func (notice *Notice) GetNoticesByKeywordAndLimit(keyword string, limit int, number int) (notices []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Notice_Name).Where(map[string]interface{}{
		"title LIKE": "%" + keyword + "%",
		"is_delete": NOTICE_NORMAL,
	}).Limit(limit, number).OrderBy("notice_id", "DESC"))
	if err != nil {
		return
	}
	notices = rs.Rows()

	return
}

//分页获取公告
func (notice *Notice) GetNoticesByLimit(limit int, number int) (notices []map[string]string, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			From(Table_Notice_Name).
			Where(map[string]interface{}{
				"is_delete": NOTICE_NORMAL,
			}).
			Limit(limit, number).
			OrderBy("notice_id", "DESC"))
	if err != nil {
		return
	}
	notices = rs.Rows()

	return
}

// 获取公告总数
func (notice *Notice) CountNotices() (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(
		db.AR().
			Select("count(*) as total").
			From(Table_Notice_Name).
			Where(map[string]interface{}{
			"is_delete": NOTICE_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

// 根据关键字获取公告总数
func (notice *Notice) CountNoticesByKeyword(keyword string) (count int64, err error) {

	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().
		Select("count(*) as total").
		From(Table_Notice_Name).
		Where(map[string]interface{}{
			"title LIKE": "%" + keyword + "%",
			"is_delete": NOTICE_NORMAL,
		}))
	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}
