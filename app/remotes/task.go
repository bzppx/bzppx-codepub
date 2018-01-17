package remotes

import (
	"encoding/json"
	"strconv"
	"bzppx-codepub/app/models"
	"time"
)

var Task = TaskRemote{}

const (
	Rpc_Task_Service           = "ServiceTask"
	Rpc_Task_Method_Publish    = Rpc_Task_Service+".Publish"
	Rpc_Task_Method_GetStatus    = Rpc_Task_Service+".Status"
	Rpc_Task_Method_Delete    = Rpc_Task_Service+".Delete"
)

type TaskRemote struct {
	BaseRemote
}

// 发布
func (this *TaskRemote) Publish(ip string, port string, args map[string]interface{}) error {
	_, err := this.Call(ip, port, Rpc_Task_Method_Publish, args)
	return err
}

// 获取节点执行结果
func (this *TaskRemote) GetResults(ip string, port string, args map[string]interface{}) (bool, error) {

	replay, err := this.Call(ip, port, Rpc_Task_Method_GetStatus, args)
	if err != nil {
		return false, err
	}

	res := map[string]string{}
	json.Unmarshal([]byte(replay), &res)

	// 任务执行完成
	if res["status"] == strconv.Itoa(models.TASKLOG_STATUS_FINISH) {
		taskLogId := args["task_log_id"].(string)
		taskLogValue := map[string]interface{}{
			"status": res["status"],
			"is_success": res["is_success"],
			"result": res["result"],
			"commit_id": res["commit_id"],
			"update_time": time.Now().Unix(),
		}
		_, err := models.TaskLogModel.Update(taskLogId, taskLogValue)
		if err != nil {
			return false, err
		}
		// 删除 agent
		_, err = this.Call(ip, port, Rpc_Task_Method_Delete, args)
		return true, nil
	} else {
		return false, nil
	}
}