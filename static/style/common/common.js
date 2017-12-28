var Common = {
	
	ajaxSubmit : function (url, data) {
		$.ajax({
			type : 'post',
			url : url,
			data : {'arr':data},
			dataType: "json",
			success : function(response) {
				if(response.code == 0) {
					Layers.failedMsg(response.message)
				} else {
					Layers.successMsg(response.message)
				}
				Common.redirect(response.redirect.url);
			},
			error : function(response) {
				Layers.failedMsg(response.message)
			}
		});
	},

	redirect: function (redirect) {
		if(redirect) {
			setTimeout(function() {
				location.href = redirect;
			}, 2000);
			setTimeout(function() {
				location.reload();
			}, 2000);
		}
	}
};
