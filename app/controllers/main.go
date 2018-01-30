package controllers

import (
	"bzppx-codepub/app/models"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"time"
)

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	this.Data["isAdmin"] = (this.isAdmin() || this.isRoot())
	this.viewLayoutTitle("CodePub POWVEREDBY BZPPX", "main/index", "main")
}

func (this *MainController) Default() {
	var err error

	// 获取最新公告
	notices, err := models.NoticeModel.GetNoticesByLimit(0, 6)
	if err != nil {
		this.ErrorLog("获取最新公告失败：" + err.Error())
		this.viewError("获取最新公告失败")
	}

	this.Data["notices"] = notices
	this.viewLayoutTitle("首页", "main/default", "index")
}

// ajax 获取我的项目总数
func (this *MainController) GetMyProjectCount() {

	var err error
	// 我的项目数
	var projectCount int64
	// 我的项目组数
	var groupCount int64
	// 我的发布失败次数
	var failedPublish int64
	// 我的发布成功次数
	var successPublish int64
	if (this.isAdmin() || this.isRoot()) {
		projectCount, err = models.ProjectModel.CountProjects()
		groupCount, err = models.GroupModel.CountGroups()
	}else {
		userProjects, err := models.UserProjectModel.GetUserProjectByUserId(this.UserID)
		if err != nil {
			this.ErrorLog("获取我的项目组数据失败: " + err.Error())
			this.jsonError("获取数据失败")
		}
		projectIds := []string{}
		for _, userProject := range userProjects {
			projectIds = append(projectIds, userProject["project_id"])
		}
		projectCount = int64(len(projectIds))
		groupCount, err = models.ProjectModel.CountGroupByProjectIds(projectIds)
		if err != nil {
			this.ErrorLog("获取我的项目组数据失败: " + err.Error())
			this.jsonError("获取数据失败")
		}
	}
	if err != nil {
		this.ErrorLog("获取我的项目组数据失败：" + err.Error())
		this.jsonError("获取数据失败")
	}

	tasks, err := models.TaskModel.GetTasksByUserId(this.UserID)
	if err != nil {
		this.ErrorLog("获取我的发布成功次数失败：" + err.Error())
		this.jsonError("获取数据失败")
	}
	taskIds := []string{}
	for _, task := range tasks {
		taskIds = append(taskIds, task["task_id"])
	}
	if len(taskIds) > 0 {
		failedPublish, err = models.TaskLogModel.CountTaskLogByTaskIdsAndIsSuccess(taskIds, models.TASKLOG_FAILED)
		if err != nil {
			this.ErrorLog("获取我的项目数据失败：" + err.Error())
			this.jsonError("获取数据失败")
		}
		successPublish = int64(len(taskIds)) - failedPublish;
	}

	data := map[string]interface{}{
		"project_count": projectCount,
		"group_count": groupCount,
		"success_publish_count": successPublish,
		"failed_publish_count": failedPublish,
	}

	this.jsonSuccess("ok", data)
}

// ajax 获取活跃项目排行榜
func (this *MainController) GetActiveProjectRank() {

	projectPublishCountRank := []map[string]string{}
	projectCountIds, err := models.TaskModel.GetProjectIdsOrderByCountProjectLimit(15)
	if err != nil {
		this.ErrorLog("获取项目总排行数据失败：" + err.Error())
		this.jsonError("获取数据失败")
	}
	if len(projectCountIds) > 0 {
		projectIds := []string{}
		for _, projectCountId := range projectCountIds {
			projectIds = append(projectIds, projectCountId["project_id"])
		}
		projects, err := models.ProjectModel.GetProjectByProjectIds(projectIds)
		if err != nil {
			this.ErrorLog("获取项目总排行数据失败：" + err.Error())
			this.jsonError("获取数据失败")
		}
		for _, projectCountId := range projectCountIds {
			projectPublishCount := map[string]string{
				"project_name": "",
				"total": projectCountId["total"],
			}
			for _, project := range projects {
				if projectCountId["project_id"] == project["project_id"] {
					projectPublishCount["project_name"] = project["name"]
					break
				}
			}
			projectPublishCountRank = append(projectPublishCountRank, projectPublishCount)
		}
	}

	this.jsonSuccess("ok", projectPublishCountRank)
}

// ajax 获取发布统计数据
func (this *MainController) GetPublishData() {

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Unix()
	yesterday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day() - 1, 0, 0, 0, 0, time.UTC).Unix()
	todayTaskCount, todayUserCount, err := models.TaskModel.CountTaskAndUserByCreateTime(today, time.Now().Unix())
	if err != nil {
		this.ErrorLog("获取今日发布任务总数失败："+err.Error())
		this.jsonError("获取数据失败")
	}
	yesterdayTaskCount, yesterdayUserCount, err := models.TaskModel.CountTaskAndUserByCreateTime(yesterday, today)
	if err != nil {
		this.ErrorLog("获取昨日发布任务总数失败："+err.Error())
		this.jsonError("获取数据失败")
	}


	data := map[string]interface{}{
		"task_total": map[string]interface{}{
			"today": todayTaskCount,
			"yesterday": yesterdayTaskCount,
		},
		"user_total": map[string]interface{}{
			"today": todayUserCount,
			"yesterday": yesterdayUserCount,
		},
		"success_node": map[string]interface{}{
			"today": "",
			"yesterday": "",
		},
		"failed_node": map[string]interface{}{
			"today": "",
			"yesterday": "",
		},
	}

	this.jsonSuccess("ok", data)
}

// ajax 获取服务器状态
func (this *MainController) ServerStatus() {
	vm, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	d, _ := disk.Usage("/")

	data := map[string]interface{}{
		"memory_used_percent": int(vm.UsedPercent),
		"cpu_used_percent": int(cpuPercent[0]),
		"disk_used_percent": int(d.UsedPercent),
	}

	this.jsonSuccess("ok", data)
}