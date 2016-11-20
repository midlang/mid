#!/bin/bash

set -e

cd ./src/cmd/midc
echo "Installing compiler: midc"
go install
cd ../../..

generators='
go
'

for lang in $generators
do
	_pwd=`pwd`
	cd ./src/cmd/gen$lang
	echo "Installing generator: gen$lang"
	go install
	cd $_pwd
done

echo "Coping config file"
cp ./midconfig $HOME/midconfig

echo "Coping templates"
templates_dir=$HOME/mid_templates
if [[ -d "$templates_dir" ]]; then
	rm -r $templates_dir
fi
cp -r ./templates $templates_dir
