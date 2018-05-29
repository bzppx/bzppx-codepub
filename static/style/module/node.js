/**
 * 节点
 * Copyright (c) 2017 phachon@163.com
 */
var Node = {


	/**
	 * 获取节点状态
	 */
	GetStatus: function () {

		var nodeIds = [];
		$("td[data-name='node_id']").each(function () {
			var nodeId = $(this).parent().attr("data-row");
			nodeIds.push(nodeId);
		});
		console.log(nodeIds);
		if(nodeIds.length == 0) {
			return;
		}
		var nodeIdsJson = JSON.stringify(nodeIds);

		$.ajax({
			type: 'POST',
			url: '/node/status',
			data: 'node_ids=' + nodeIdsJson,
			async: true,
			dataType: 'json',
			success: function(response) {
				if (!response) {
					return;
				}
				if(response.code == "0") {
					console.log(response.message);
					return;
				}
				var nodesStatus = response.data;
				for (var i = 0; i < nodesStatus.length; i++) {
					var nodeStatus = nodesStatus[i];
					var nodeId = nodeStatus.node_id;
					var status = nodeStatus.status;
					var version = nodeStatus.version;
					var tdStatus = $("td[data-row="+nodeId+"][data-name='status']");
					// failed
					if(status == 0) {
						// tdStatus.html('<i class="glyphicon glyphicon-remove-circle text-danger"></i>');
						tdStatus.html('<span class="label label-danger">null</span>');
					}
					// success
					if(status == 1) {
						// tdStatus.html('<i class="glyphicon glyphicon-ok-circle text-success"></i>');
						tdStatus.html('<span class="label label-success">'+version+'</span>');
					}
				}
			}
		});
	}
};