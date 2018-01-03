/**
 * 模块
 * Copyright (c) 2017 phachon@163.com
 */
var User = {

	defaults : function(defaults) {
		var arr = defaults.split(",");
		console.log(arr);
		$('[data-type="module"][name="module_id"]').each(function() {
			console.log(this.value);
			var checked = $.inArray(this.value, arr) > -1 ? true : false;
			this.checked = checked;
		});
	},

	module : function(element) {
		var userId = $("input[name='user_id']").val();
		var id = $(element).val();
		var checked = $(element).is(':checked');
		var type = $(element).attr("data-type");
		var moduleIds = [];

		if(type == 'group') {
			$('[name="module_id"][data-parent='+id+']').each(function() {
				this.checked = checked;
				moduleIds.push($(this).val());
			});
		}
		if(type == 'node') {
			moduleIds.push(id);
			var parentId = $(element).attr("data-parent");
			var flag = true;
			$('[data-type="module"][data-parent='+parentId+']').each(function() {
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
		console.log(moduleIds);
		User.save(moduleIds.join(","), userId, isChecked);
	},
	
	save: function (moduleIds, userId, isChecked) {
		$.ajax({
			type : 'post',
			url : '/user/moduleSave',
			data : {'module_ids':moduleIds, 'user_id': userId, 'is_check': isChecked},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					Form.failedBox(response.message);
					console.log("失败: "+response.message)
				}
			},
			error : function(response) {
				Form.failedBox("添加节点失败");
				console.log("失败");
			}
		});
	}
};