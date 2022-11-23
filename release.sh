#!/bin/bash

set -e

CMD_GO=go
RELEASE_DIR=targets
LANGUAGES=`cat languages.txt`
VERSION=`cat VERSION`

cd ./hack
source ./genmeta.sh $VERSION
cd ..

function release_mid_for() {
	local _os=$1
	local _arch=$2

	local _target_dir=mid$VERSION.$_os-$_arch
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
	for _lang in $LANGUAGES
	do
		local _bin=mid-gen-$_lang
		echo "GOOS=$_os GOARCH=$_arch $CMD_GO build -o $_target_dir/bin/$_bin$_suffix ./src/cmd/$_bin"
		GOOS=$_os GOARCH=$_arch $CMD_GO build -o $_target_dir/bin/$_bin$_suffix ./src/cmd/$_bin
	done

	# Coping files
	cp ./midconfig $_target_dir/
	cp VERSION $_target_dir/
	cp ./README.md $_target_dir/
	cp -r ./templates $_target_midroot/
	cp -r ./extensions $_target_midroot/

	cp ./install.sh $_target_dir/install.sh
	chmod +x $_target_dir/install.sh

	# Targz or zip (for windows)
	if [[ "$_os" == "windows" ]]; then
		zip -q $_target_dir.zip -r $_target_dir
		mv $_target_dir.zip $RELEASE_DIR/$VERSION/
	else
		tar zcf $_target_dir.tar.gz $_target_dir
		mv $_target_dir.tar.gz $RELEASE_DIR/$VERSION/
	fi
	rm -rf $_target_dir 2> /dev/null
}

if [[ -d "$RELEASE_DIR/$VERSION" ]]; then
	rm -r $RELEASE_DIR/$VERSION
fi
mkdir -p $RELEASE_DIR/$VERSION

release_mid_for windows 386
release_mid_for windows amd64
release_mid_for windows arm64
release_mid_for linux 386
release_mid_for linux amd64
release_mid_for linux arm64
release_mid_for darwin amd64
release_mid_for darwin arm64
