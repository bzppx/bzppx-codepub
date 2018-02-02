/**
 * 首页
 * Copyright (c) 2017 phachon@163.com
 */

var Main = {

	/**
	 * 获取项目总数
	 * @param url
	 * @constructor
	 */
	GetProjectCount: function (url) {
		$.ajax({
			type : 'post',
			url : url,
			data : {'arr':''},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log(response.message);
					return false
				}
				$("#project_count_text").text(response.data.project_count);
				$("#group_count_text").text(response.data.group_count);
				$("#success_publish_text").text(response.data.success_publish_count);
				$("#failed_publish_text").text(response.data.failed_publish_count);
			},
			error : function(response) {
				console.log(response.message)
			}
		});
	},
	
	/**
	 * 获取活跃项目排行榜
	 * @param element
	 * @param url
	 * @constructor
	 */
	GetActiveProjectRank: function (element, url) {
		$.ajax({
			type : 'post',
			url : url,
			data : {'arr':''},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log(response.message);
					return false
				}
				
				Morris.Bar({
					element: element,
					data: response.data,
					xkey: 'project_name',
					ykeys: ['total'],
					labels: ['发布次数'],
					barRatio: 0.4,
					xLabelAngle: 65,
					hideHover: 'auto',
					resize: true
				});
			},
			error : function(response) {
				console.log(response.message)
			}
		});
	},

	/**
	 * 获取发布数据
	 * @param url
	 * @constructor
	 */
	GetPublishData: function (url) {
		$.ajax({
			type : 'post',
			url : url,
			data : {'arr':''},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log(response.message);
					return false
				}
				$("td[data-name='today_task_total']").text(response.data.task_total.today);
				$("td[data-name='yesterday_task_total']").text(response.data.task_total.yesterday);
				$("td[data-name='today_user_total']").text(response.data.user_total.today);
				$("td[data-name='yesterday_user_total']").text(response.data.user_total.yesterday);
				$("td[data-name='today_success_node']").text(response.data.success_tasklog.today);
				$("td[data-name='yesterday_success_node']").text(response.data.success_tasklog.yesterday);
				$("td[data-name='today_failed_node']").text(response.data.failed_tasklog.today);
				$("td[data-name='yesterday_failed_node']").text(response.data.failed_tasklog.yesterday);
			},
			error : function(response) {
				console.log(response.message)
			}
		});
	},
	
	/**
	 * 获取服务器状态
     * @param url
	 * @constructor
	 */
	GetServerStatus: function (url) {
		$.ajax({
			type : 'post',
			url : url,
			data : {'arr':''},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log(response.message);
					return false
				}
				var cpu = response.data.cpu_used_percent;
				var memory = response.data.memory_used_percent;
				var disk = response.data.disk_used_percent;
				// cpu
				$(".cpu_text").each(function () {
					$(this).text(cpu+"%")
				});
				$("#cpu_progress").attr("aria-valuenow", cpu);
				$("#cpu_progress").attr('style', 'min-width: 2em; width: '+cpu+'%');

				// memory
				$(".memory_text").each(function () {
					$(this).text(memory+"%")
				});
				$("#memory_progress").attr("aria-valuenow", memory);
				$("#memory_progress").attr('style', 'min-width: 2em; width: '+memory+'%');

				// disk
				$(".disk_text").each(function () {
					$(this).text(disk+"%")
				});
				$("#disk_progress").attr("aria-valuenow", disk);
				$("#disk_progress").attr('style', 'min-width: 2em; width: '+disk+'%');

			},
			error : function(response) {
				console.log(response.message)
			}
		});
	},
	
	/**
	 * 获取联系人信息
     * @param url
	 * @constructor
	 */
	GetContact: function (url) {
		$.ajax({
			type : 'post',
			url : url,
			data : {},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log(response.message);
					return false
				}
                       
				var html = ""
				for (var i in response.data) {
					html += '<li class="list-group-item">' + response.data[i]['position'] + '(' + response.data[i]['name'] + ')：'
					if (response.data[i]['mobile'] != "") {
						html += '<span class="glyphicon glyphicon-phone"></span>&nbsp;' + response.data[i]['mobile'] + '&nbsp;&nbsp;'
					}
					if (response.data[i]['telephone'] != "") {
						html += '<span class="glyphicon glyphicon-phone-alt"></span>&nbsp;' + response.data[i]['telephone'] + '&nbsp;&nbsp;'
					}
					if (response.data[i]['email'] != "") {
						html += '<span class="glyphicon glyphicon-envelope"></span>&nbsp;' + response.data[i]['email']
					}
					html += '</li>'
				}
				$('#contact').html(html)
			},
			error : function(response) {
				console.log(response.message)
			}
		});
	}
};