package models

import (
	"bzppx-codepub/app/utils"

	"github.com/snail007/go-activerecord/mysql"
)

const Table_Task_Name = "task"

type Task struct {
}

var TaskModel = Task{}

// 根据 task_id 获取任务
func (t *Task) GetTaskByTaskId(taskId string) (tasks map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Task_Name).Where(map[string]interface{}{
		"task_id": taskId,
	}))
	if err != nil {
		return
	}
	tasks = rs.Row()
	return
}

// 插入一条任务
func (l *Task) Insert(task map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_Task_Name, task))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 通过module_id和task_id查找task
func (l *Task) GetTaskByModuleIdsAndTaskIds(moduleIds, taskIds []string) (task []map[string]string, err error) {
	db := G.DB()
	where := make(map[string]interface{})
	where["task_id"] = taskIds
	var rs *mysql.ResultSet
	if len(moduleIds) > 0 {
		where["module_id"] = moduleIds
	}
	rs, err = db.Query(db.AR().From(Table_Task_Name).Where(where))
	if err != nil {
		return
	}
	task = rs.Rows()
	return
}

func (l *Task) GetTaskByModuleId(moduleId string, limit, number int) (tasks []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Task_Name).Where(map[string]interface{}{
		"module_id": moduleId,
	}).Limit(limit, number).OrderBy("task_id", "DESC"))

	if err != nil {
		return
	}
	tasks = rs.Rows()
	return
}

func (l *Task) GetTaskByModuleIdAndUserId(moduleId, userId string, limit, number int) (tasks []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_Task_Name).Where(map[string]interface{}{
		"module_id": moduleId,
		"user_id":   userId,
	}).Limit(limit, number).OrderBy("task_id", "DESC"))

	if err != nil {
		return
	}
	tasks = rs.Rows()
	return
}

func (l *Task) CountTaskByModuleId(moduleId string) (count int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From(Table_Task_Name).Where(map[string]interface{}{
		"module_id": moduleId,
	}))

	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}

func (l *Task) CountTaskByModuleIdAndUserId(moduleId, userId string) (count int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().Select("count(*) as total").From(Table_Task_Name).Where(map[string]interface{}{
		"module_id": moduleId,
		"user_id":   userId,
	}))

	if err != nil {
		return
	}
	count = utils.NewConvert().StringToInt64(rs.Value("total"))
	return
}