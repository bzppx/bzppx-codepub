/**
 * 任务
 * Copyright (c) 2017 phachon@163.com
 */
var Task = {

	taskLogCreateHtml: '<span class="label label-info">创建</span>',
	taskLogStartHtml: '<span class="label label-primary">开始</span>',
	taskLogEndHtml: '<span class="label label-success">完成</span>',
	taskLogSuccessHtml: '<span class="label label-success">成功</span>',
	taskLogFailedHtml: '<span class="label label-danger end_failed">失败</span>',
	taskLogLoadingHtml: '<img src="/static/style/layui/layer/theme/default/loading-2.gif" style="max-width:17px;">',

	/**
	 * 加载状态(30s)
	 */
	loadTaskLogStatus: function () {

		var taskLogIds = [];

		$("td[data-name='status'][data-status!='2']").each(function () {
			var taskLogId = $(this).parent().attr("data-row");
			taskLogIds.push(taskLogId);
		});
		
		console.log(taskLogIds);
		if(taskLogIds.length == 0) {
			return 0;
		}

		taskLogIds = JSON.stringify(taskLogIds);

		$.ajax({
			type: 'POST',
			url: '/taskLog/getTaskLogs',
			data: 'taskLog_ids=' + taskLogIds,
			dataType: 'json',
			success: function(response) {
				if (!response) {
					return;
				}
				if(response.code == "0") {
					console.log(response.message);
					return;
				}

				var taskLogs = response.data;
				for (var i = 0; i < taskLogs.length; i++) {
					var taskLog = taskLogs[i];
					var taskLogId = taskLog.task_log_id;
					var status = taskLog.status;
					var isSuccess = taskLog.is_success;
					var commitId = taskLog.commit_id;
					var exitResult = taskLog.result;
					var updateTime = taskLog.update_time;
					
					var statusSelect = $('[data-row=' + taskLogId + '] [data-name="status"]');
					var isSuccessSelect = $('[data-row=' + taskLogId + '] [data-name="is_success"]');
					var commitIdSelect = $('[data-row=' + taskLogId + '] [data-name="commit_id"] strong');
					var resultSelect = $('#result_'+taskLogId+" label");
					var updateTimeSelect = $('[data-row=' + taskLogId + '] [data-name="update_time"]');

					statusSelect.attr('data-status', status);
					isSuccessSelect.attr('data-success', isSuccess);
					commitIdSelect.html(commitId);
					resultSelect.html(exitResult);
					updateTimeSelect.html($.myTime.UnixToDate(updateTime, true, 8));

					// create status
					if(status == 0) {
						statusSelect.html(Task.taskLogCreateHtml);
						isSuccessSelect.html(Task.taskLogLoadingHtml);
					}
					// start status
					if(status == 1) {
						statusSelect.html(Task.taskLogStartHtml);
						isSuccessSelect.html(Task.taskLogLoadingHtml);
					}
					// end status
					if(status == 2) {
						statusSelect.html(Task.taskLogEndHtml);
						if (isSuccess == 1) {
							isSuccessSelect.html(Task.taskLogSuccessHtml);
							Task.updateProgress("success");
						}else {
							isSuccessSelect.html(Task.taskLogFailedHtml);
							Task.updateProgress("danger");
						}
					}
				}
			}
		});
	},

	// 刷新进度条
	updateProgress: function (status) {
		if (status == undefined && $(".end_failed").length > 0) {
			status = "danger"
		}
		var totalRow = $("tr[data-row]").length;
		var end = $("td[data-name='status'][data-status='2']").length;
		var p = (end/totalRow).toFixed(2) * 100;
		console.log("progress: "+p+"%");
		var taskProgressSelect = $("#task_progress");
		if (status == "danger") {
			taskProgressSelect.removeClass("progress-bar-success");
			taskProgressSelect.addClass("progress-bar-danger");
		}
		taskProgressSelect.attr('aria-valuenow', p);
		taskProgressSelect.attr('style', 'min-width: 2em; width: '+p+'%');
		$("#task_progress_content").html(p+"%");
	}
};