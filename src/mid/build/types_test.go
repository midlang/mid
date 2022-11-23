package build

import "testing"

func TestParseTag(t *testing.T) {
	var tag = Tag(`key:"value"`)
	t.Logf("tag.key=%s", tag.Get("key"))
}
