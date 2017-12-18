var Layers = {

	/**
	 * 皮肤
	 */
	skin : 'layui-layer-lan',

	/**
	 * alert 提示信息框
	 * @param title
	 */
	alert : function (title) {
		layer.alert(title, {
			icon: 1,
			skin: Layers.skin
		})
	},

	/**
	 * confirm 提示框
	 * @param title
	 * @param url
	 */
	confirm: function (title) {
		layer.confirm(title, {
			btn: ['是','否']
		}, function() {
			layer.msg('的确很重要', {icon: 1});
		}, function() {
			layer.msg('也可以这样', {
				time: 20000,
				btn: ['明白了', '知道了']
			});
		});
	},

	/**
	 * iframe 窗
	 */
	iframe: function (title, url, width, height) {
		height = height||"500px";
		width = width||"1000px";
		layer.open({
			type: 2,
			skin: Layers.skin,
			title: title,
			shadeClose: true,
			shade : 0.6,
			maxmin: true,
			area: [width, height],
			content: url
		});
	}
};