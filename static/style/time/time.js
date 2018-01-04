$(function(){
	$(".date").datetimepicker({
        format: "yyyy-mm-dd",
		minView: 2,
		language: 'zh-CN',
		autoclose: true,
    });
	
	$(".time").datetimepicker({
        format: "yyyy-mm-dd hh:ii",
		minView: 0,
		language: 'zh-CN',
		autoclose: true,
    });
})