package build

import (
	"strings"
	"testing"
)

func TestReadBeanIds(t *testing.T) {
	const content = `
	x=1
	y= 2
	z = 3
	`
	reader := strings.NewReader(content)
	ids, err := ReadBeanIds(reader)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(ids)
	}
}
