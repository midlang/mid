#!/bin/bash

set -e

file=../src/mid/meta.go
version=$1

if [[ "$version" == "" ]]; then
	echo "var version missing as first argument"
	exit 1
fi

cat > $file <<EOF
package mid

import "fmt"

type Map map[string]interface{}

func (m Map) String(key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return fmt.Sprintf("%v", v)
	}
}

var Meta = Map {
	"version": "$version",
	"officialAuthor": "midc",
}
EOF

gofmt -w $file
