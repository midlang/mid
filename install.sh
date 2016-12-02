#!/bin/bash

set -e

cd ./src/cmd/midc
echo "Installing compiler: midc"
go install
cd ../../..

languages=`cat languages.txt`

for lang in $languages
do
	_pwd=`pwd`
	cd ./src/cmd/gen$lang
	echo "Installing generator: gen$lang"
	go install
	cd $_pwd
done

echo "Coping config file"
cp ./midconfig $HOME/midconfig

echo "Coping templates and extensions"
midroot=$HOME/.mid
if [[ -d "$midroot" ]]; then
	rm -r $midroot
fi
mkdir -p $midroot
cp -r ./templates $midroot/
cp -r ./extensions $midroot/
