var Layers = {

	/**
	 * 皮肤
	 */
	skin : 'layui-layer-lan',

	/**
	 * success 提示信息框
	 * @param title
	 */
	success : function (title) {
		layer.alert(title+"<br/>", {
			title: "操作成功",
			icon: 1,
			skin: Layers.skin,
			closeBtn: 0
		})
	},

	/**
	 * error 提示信息框
	 * @param title
	 */
	error : function (title) {
		layer.alert(title, {
			title: "操作失败",
			icon: 2,
			skin: Layers.skin,
			closeBtn: 0
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
	 * bind iframe 窗
	 */
	bindIframe: function (element, title, url, width, height) {
		$(element).each(function () {
			$(this).bind('click', function () {
				height = height||"500px";
				width = width||"1000px";
				url = url|| $(this).attr("data-link");
				layer.open({
					type: 2,
					skin: Layers.skin,
					title: title,
					shadeClose: true,
					shade : 0.6,
					maxmin: true,
					area: [width, height],
					content: url,
					padding:"10px"
				});
			})
		})
	}
};