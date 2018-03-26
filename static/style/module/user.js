/**
 * 项目
 * Copyright (c) 2017 phachon@163.com
 */
var User = {

	defaults : function(defaults) {
		var arr = defaults.split(",");
		console.log(arr);
		$('[data-type="project"][name="project_id"]').each(function() {
			console.log(this.value);
			var checked = $.inArray(this.value, arr) > -1 ? true : false;
			this.checked = checked;
		});
	},

	project : function(element) {
		var userId = $("input[name='user_id']").val();
		var id = $(element).val();
		var checked = $(element).is(':checked');
		var type = $(element).attr("data-type");
		var projectIds = [];

		if(type == 'group') {
			$('[name="project_id"][data-parent='+id+']').each(function() {
				this.checked = checked;
				projectIds.push($(this).val());
			});
		}
		if(type == 'project') {
			projectIds.push(id);
			var parentId = $(element).attr("data-parent");
			var flag = true;
			$('[data-type="project"][data-parent='+parentId+']').each(function() {
				if (this.checked != checked) {
					flag = false;
				}
			});
			if (flag) {
				$('[data-type="group"][value='+parentId+']').each(function() {
					this.checked = checked
				});
			}
		}

		var isChecked = 1;
		if(checked == false) {
			isChecked = 0;
		}
		console.log(projectIds);
		User.save(projectIds.join(","), userId, isChecked);
	},
	
	save: function (projectIds, userId, isChecked) {
		$.ajax({
			type : 'post',
			url : '/user/projectSave',
			data : {'project_ids':projectIds, 'user_id': userId, 'is_check': isChecked},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log("失败: "+response.message)
					Common.errorBox(Form.failedBox, response.message);
				}
			},
			error : function(response) {
				console.log("失败");
				Common.errorBox(Form.failedBox, "添加节点失败");
			}
		});
	}
};