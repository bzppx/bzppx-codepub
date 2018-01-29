/**
 * 首页
 * Copyright (c) 2017 phachon@163.com
 */

var Main = {

	/**
	 * 获取服务器状态
	 * @param url
	 * @constructor
	 */
	GetServerStatus: function (url) {
		$.ajax({
			type : 'post',
			url : url,
			data : {'arr':data},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					console.log(response.message);
					return
				}
				// 获取服务器信息
				
			},
			error : function(response) {
				console.log(response.message)
			}
		});
	}
};