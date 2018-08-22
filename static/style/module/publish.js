/**
 * 发布
 * Copyright (c) 2017 phachon@163.com
 */
var Publish = {

	node : function(element) {
		var id = $(element).val();
		var checked = $(element).is(':checked');
		var type = $(element).attr("data-type");
		var nodeIds = [];

		if(type == 'group') {
			$('[name="node_id"][data-parent='+id+']').each(function() {
				this.checked = checked;
				nodeIds.push($(this).val());
			});
		}
		if(type == 'node') {
			nodeIds.push(id);
			var parentId = $(element).attr("data-parent");
			var flag = true;
			$('[data-type="node"][data-parent='+parentId+']').each(function() {
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
		console.log(nodeIds);
	}
};