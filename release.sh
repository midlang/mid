#!/bin/bash

set -e

function verbose() {
	echo "$@"
}

Go=go
released_dir=targets
generators='
	go
'
version_file=VERSION
version=`cat $version_file`

cat > ./src/mid/meta.go <<EOF
package mid

var Meta = map[string]interface{} {
	"version": "$version",
}
EOF

function mid_release_with_os_cpu() {
	local _version=$1
	local _os=$2
	local _cpu=$3

	local _target_dir=mid$_version.$_os-$_cpu
	mkdir -p $_target_dir/bin

	# Building compiler `midc`
	GOOS=$_os GOARCH=$_cpu $Go build -o $_target_dir/bin/midc ./src/cmd/midc/

	# Building generators
	for lang in $generators
	do
		echo "$Go build -o $_target_dir/bin/gen$lang ./src/cmd/gen$lang"
		$Go build -o $_target_dir/bin/gen$lang ./src/cmd/gen$lang
	done

	# Coping files
	cp ./midconfig $_target_dir/
	cp $version_file $_target_dir/
	cp ./README.md $_target_dir/
	cp -r ./templates $_target_dir/mid_templates

	# Targz
	mkdir -p $released_dir/$_version
	tar zcf $_target_dir.tar.gz $_target_dir
	mv $_target_dir.tar.gz $released_dir/$_version/
	rm -rf $_target_dir 2> /dev/null
}

mid_release_with_os_cpu $version windows 386
mid_release_with_os_cpu $version windows amd64
mid_release_with_os_cpu $version linux 386
mid_release_with_os_cpu $version linux amd64
mid_release_with_os_cpu $version darwin 386
mid_release_with_os_cpu $version darwin amd64

