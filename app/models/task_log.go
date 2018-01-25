package models

import "github.com/snail007/go-activerecord/mysql"

const (
	TASKLOG_STATUS_CREATE = 0 // 任务状态，创建
	TASKLOG_STATUS_SATART = 1 // 任务状态，开始执行
	TASKLOG_STATUS_FINISH = 2 // 任务状态，执行完成

	TASKLOG_FAILED  = 0 // 执行结果状态，失败
	TASKLOG_SUCCESS = 1 // 执行结果状态，成功
)

const Table_TaskLog_Name = "task_log"

type TaskLog struct {
}

var TaskLogModel = TaskLog{}

// 根据 task_log_id 获取任务日志
func (t *TaskLog) GetTaskLogByTaskLogId(taskLogId string) (tasLogs map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"task_log_id": taskLogId,
	}))
	if err != nil {
		return
	}
	tasLogs = rs.Row()
	return
}

// 根据 task_id 获取任务日志
func (t *TaskLog) GetTaskLogByTaskId(taskId string) (tasLogs []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"task_id": taskId,
	}))
	if err != nil {
		return
	}
	tasLogs = rs.Rows()
	return
}

// 根据状态获取任务日志
func (t *TaskLog) GetTaskLogByStatus(status int) (tasLogs []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"status": status,
	}))
	if err != nil {
		return
	}
	tasLogs = rs.Rows()
	return
}

// 根据 is_success 获取任务日志
func (t *TaskLog) GetTaskLogBySuccess(isSuccess int) (tasLogs []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"is_success": isSuccess,
	}))
	if err != nil {
		return
	}
	tasLogs = rs.Rows()
	return
}

// 插入一条任务日志
func (l *TaskLog) Insert(taskLog map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Insert(Table_TaskLog_Name, taskLog))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 插入多条任务日志
func (l *TaskLog) InsertBatch(taskLogs []map[string]interface{}) (err error) {
	db := G.DB()
	_, err = db.Exec(db.AR().InsertBatch(Table_TaskLog_Name, taskLogs))
	return
}

// 查找未执行完的日志
func (l *TaskLog) GetExcutingTaskIdByTaskLog() (taskIds []string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"status": [2]int{TASKLOG_STATUS_CREATE, TASKLOG_STATUS_SATART},
	}).GroupBy("task_id"))
	if err != nil {
		return
	}

	taskLogs := rs.Rows()
	taskIds = make([]string, len(taskLogs))
	for index, taskLog := range taskLogs {
		taskIds[index] = taskLog["task_id"]
	}

	return
}

// 根据多个 task_id 获取任务日志
func (t *TaskLog) GetTaskLogByTaskIds(taskIds []string) (tasLogs []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"task_id": taskIds,
	}))
	if err != nil {
		return
	}
	tasLogs = rs.Rows()
	return
}

// 修改task_log
func (t *TaskLog) Update(taskLogId string, taskLogValue map[string]interface{}) (id int64, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Exec(db.AR().Update(Table_TaskLog_Name, taskLogValue, map[string]interface{}{
		"task_log_id": taskLogId,
	}))
	if err != nil {
		return
	}
	id = rs.LastInsertId
	return
}

// 获取失败的 task_log
func (t *TaskLog) GetFailedTaskLogByTaskIds(taskIds []string) (tasLogs []map[string]string, err error) {
	db := G.DB()
	var rs *mysql.ResultSet
	rs, err = db.Query(db.AR().From(Table_TaskLog_Name).Where(map[string]interface{}{
		"task_id": taskIds,
		"is_success": TASKLOG_FAILED,
	}))
	if err != nil {
		return
	}
	tasLogs = rs.Rows()
	return
}