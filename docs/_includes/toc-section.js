function toccheckboxChanged(e) {
	//console.log("toccheckboxChanged");
	if (e.is(':checked')) {
		$('.toccontent').hide();
	} else {
		$('.toccontent').show();
	}
}

$('.overlay').visibility({
	type: 'fixed',
	offset: 20,
	onTopVisible: function(calculations) {
		console.log("onTopVisible");
		// 如果原来是隐藏状态则改为显示
		var e = $('#toctogglecheckbox');
		if (e.is(':checked')) {
			e.prop("checked", false);
			toccheckboxChanged(e);
		}
	},
	onBottomPassed: function(calculations) {
		console.log("onBottomPassed");
		// 如果原来是显示状态则改为隐藏
		var e = $('#toctogglecheckbox');
		if (!e.is(':checked')) {
			e.prop("checked", true);
			toccheckboxChanged(e);
		}
	},
	onPassing: function(calculations) {
		console.log("onPassing");
		// 如果原来是显示状态则改为隐藏
		var e = $('#toctogglecheckbox');
		if (!e.is(':checked')) {
			e.prop("checked", true);
			toccheckboxChanged(e);
		}
	},
});

$('#toctogglecheckbox').change(function() {
	toccheckboxChanged($(this));
});
