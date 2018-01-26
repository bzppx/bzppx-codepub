/**
 * 图表
 * Copyright (c) 2017 phachon@163.com
 */

var Chart = {

	/**
	 * 活跃项目排行榜
	 * @param element
	 * @param data
	 */
	activeProjectRank: function (element, data) {
		Morris.Bar({
			element: element,
			data: data,
			xkey: 'project',
			ykeys: ['number'],
			labels: ['数量'],
			barRatio: 0.4,
			xLabelAngle: 35,
			hideHover: 'auto',
			resize: true
		});
	}
};