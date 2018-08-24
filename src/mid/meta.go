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

var Meta = Map{
	"version":        "0.1.3.head",
	"officialAuthor": "midc",
}
