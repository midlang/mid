{% unless page.notoc %}

function toccheckboxChanged(e) {
	//console.log("toccheckboxChanged");
	if (e.is(':checked')) {
		$('.toccontent').hide();
	} else {
		$('.toccontent').show();
	}
}

function tocMarkupHelper(ul, deeps) {
	var count = 0;
	ul.children("li").each(function() {
		count++;
		var li = $(this);
		deeps.push(count);
		li.children("a").each(function() {
			var a = $(this);
			var prepended = '<span class="toc-title-markup">' + deeps.join(".") + '</span>';
			a.html(prepended + "&nbsp;&nbsp;" + a.html());
		});
		li.children("ul").each(function() {
			tocMarkupHelper($(this), deeps.slice());
		});
		deeps.pop();
	});
}

$(".toccontent").children("ul").each(function() {
	$(this).css("padding-left", 0);
	tocMarkupHelper($(this), []);
});

/*
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
	onTopVisibleReverse: function(calculations) {
		console.log("onTopVisibleReverse");
	},
	onTopPassedReverse: function(calculations) {
		console.log("onTopPassedReverse");
	},
	onTopPassed: function(calculations) {
		console.log("onTopPassed");
	},
	onBottomVisible: function(calculations) {
		console.log("onBottomVisible");
	},
	onBottomVisibleReverse: function(calculations) {
		console.log("onBottomVisibleReverse");
	},
	onBottomPassedReverse: function(calculations) {
		console.log("onBottomPassedReverse");
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
	onPassingReverse: function(calculations) {
		console.log("onPassingReverse");
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
*/

$('#toctogglecheckbox').change(function() {
	toccheckboxChanged($(this));
});

$('.main-content h2').each(function() {
	var e = $(this);
	e.html('§ ' + e.html());
});

{% endunless %}
