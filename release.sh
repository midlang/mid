#!/bin/bash

set -e

CMD_GO=go
released_dir=targets
languages=`cat languages.txt`

version_file=VERSION
version=`cat $version_file`

cd ./hack
source ./genmeta.sh $version
cd ..

function mid_release_for() {
	local _version=$1
	local _os=$2
	local _arch=$3

	local _target_dir=mid$_version.$_os-$_arch
	local _target_midroot=$_target_dir/mid
	mkdir -p $_target_dir/bin
	mkdir -p $_target_midroot
	local _suffix=
	if [[ "$_os" == "windows" ]]; then
		_suffix=.exe
	fi

	# Building compiler `midc`
	echo "GOOS=$_os GOARCH=$_arch $CMD_GO build -o $_target_dir/bin/midc$_suffix ./src/cmd/midc/"
	GOOS=$_os GOARCH=$_arch $CMD_GO build -o $_target_dir/bin/midc$_suffix ./src/cmd/midc/

	# Building generators
	local _lang
	for _lang in $languages
	do
		local _bin=mid-gen-$_lang
		echo "GOOS=$_os GOARCH=$_arch $CMD_GO build -o $_target_dir/bin/$_bin$_suffix ./src/cmd/$_bin"
		GOOS=$_os GOARCH=$_arch $CMD_GO build -o $_target_dir/bin/$_bin$_suffix ./src/cmd/$_bin
	done

	# Coping files
	cp ./midconfig $_target_dir/
	cp $version_file $_target_dir/
	cp ./README.md $_target_dir/
	cp -r ./templates $_target_midroot/
	cp -r ./extensions $_target_midroot/
	if [[ "$_os" != "windows" ]]; then
		cp ./install.sh $_target_dir/
	fi

	# Targz or zip( for windows )
	if [[ "$_os" == "windows" ]]; then
		zip -q $_target_dir.zip -r $_target_dir
		mv $_target_dir.zip $released_dir/$_version/
	else
		tar zcf $_target_dir.tar.gz $_target_dir
		mv $_target_dir.tar.gz $released_dir/$_version/
	fi
	rm -rf $_target_dir 2> /dev/null
}

if [[ -d "$released_dir/$_version" ]]; then
	rm -r $released_dir/$_version
fi
mkdir -p $released_dir/$version

mid_release_for $version windows 386
mid_release_for $version windows amd64
mid_release_for $version linux 386
mid_release_for $version linux amd64
mid_release_for $version darwin 386
mid_release_for $version darwin amd64

