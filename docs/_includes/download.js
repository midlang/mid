
$('.ui.dropdown').dropdown();

function active_selection(menu, text) {
	$(menu + " .menu").children(".item").each(function() {
		var value = $(this).text();
		if (text === value) {
			$(this).addClass("active selected");
			$(menu + " .text").html($(this).html());
		} else {
			$(this).removeClass("active selected");
		}
	});
}
var current_os = platform.os.family;
if (current_os == "OS X") {
	active_selection("#os-selection", "macOS");
} else if (current_os == "Windows") {
	active_selection("#os-selection", "windows");
}
var current_arch = platform.os.architecture;
if (current_arch == 32) {
	active_selection("#arch-selection", "32 bit");
}

$('#download-button').click(function() {
	var os = "linux";
	var arch = "64 bit";
	var version = "";
	var suffix = ".tar.gz";
	$("#os-selection .menu").children(".selected").each(function() {
		os = $(this).text();
	});
	$("#arch-selection .menu").children(".selected").each(function() {
		arch = $(this).text();
	});
	$("#version-selection .menu").children(".selected").each(function() {
		version = $(this).text();
	});
	if (os === "macOS") {
		os = "darwin";
	}
	if (os === "windows") {
		suffix = ".zip";
	}
	if (arch === "32 bit") {
		arch = "386";
	} else if (arch === "64 bit") {
		arch = "amd64";
	}
	version = version.replace(/^v/, '');
	var addr = "https://github.com/midlang/mid/releases/download/";
	addr += "v" + version +  "/mid" + version +  "." + os +  "-" + arch + suffix;
	console.log("address: " + addr);
	window.location.href = addr;
});
