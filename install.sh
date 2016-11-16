#!/bin/bash

cat > $HOME/.midconfig <<EOF
{
	"suffix": "mid",

	"plugins": [
		{
			"lang": "go",
			"name": "std",
			"bin": "gengo",
			"supported_exts": ["proto", "redis", "mysql"]
		}
	]
}
EOF

PWD=`pwd`
cd cmd/gengo
echo "install gengo"
go install

cd $PWD
cd cmd/midc
echo "install midc"
go install

cd $PWD
echo "copy templates"
TEMP_ROOTDIR=$HOME/.mid/templates
mkdir -p $TEMP_ROOTDIR

mkdir $TEMP_ROOTDIR/go
cp -r cmd/gengo/templates/* $TEMP_ROOTDIR/go/
