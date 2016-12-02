package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/midlang/mid/src/mid"
)

var (
	ErrExtensionName = errors.New("invalid extension name")
)

const (
	FileHead       = "file_head"
	BeforeImport   = "before_import"
	InImport       = "in_import"
	AfterImport    = "after_import"
	BeforeConst    = "before_const"
	AfterConst     = "after_const"
	BeforeEnum     = "before_enum"
	AfterEnum      = "after_enum"
	BeforeStruct   = "before_struct"
	AfterStruct    = "after_struct"
	BeforeProtocol = "before_protocol"
	AfterProtocol  = "after_protocol"
	BeforeService  = "before_service"
	AfterService   = "after_service"
)

type EmbeddedPosition string

func (pos EmbeddedPosition) split() (kind string, at string) {
	index := strings.Index(string(pos), ".")
	if index >= 0 {
		return string(pos[:index]), string(pos[index+1:])
	}
	return "", string(pos)
}

func (pos EmbeddedPosition) IsValid() bool {
	kind, at := pos.split()
	switch kind {
	case "", "package", "file", "const", "enum", "struct", "protocol", "service":
	default:
		return false
	}
	switch at {
	case FileHead,
		BeforeImport,
		InImport,
		AfterImport,
		BeforeConst,
		AfterConst,
		BeforeEnum,
		AfterEnum,
		BeforeStruct,
		AfterStruct,
		BeforeProtocol,
		AfterProtocol,
		BeforeService,
		AfterService:
	default:
		return false
	}
	return true
}

func (pos EmbeddedPosition) Match(kind, at string) bool {
	k, a := pos.split()
	return (k == "" || k == kind) && a == at
}

type EmbeddedValue struct {
	Text     string `join:"text"`
	Template string `json:"template"`
}

func (v EmbeddedValue) IsValid() bool {
	return (v.Text != "" && v.Template == "") || (v.Text == "" && v.Template != "")
}

type Extention struct {
	Name    string
	Author  string
	URL     string
	Version string

	// language -> position -> embedded_values
	EmbeddedAt map[string]map[EmbeddedPosition][]EmbeddedValue
}

func (e Extention) Find(lang, kind, at string) []EmbeddedValue {
	if e.EmbeddedAt == nil {
		return nil
	}
	m, ok := e.EmbeddedAt[lang]
	if !ok || m == nil {
		return nil
	}
	values := make([]EmbeddedValue, 0)
	for pos, vals := range m {
		if pos.Match(kind, at) {
			values = append(values, vals...)
		}
	}
	return values
}

type ExtentionKey struct {
	Author string
	Name   string
}

func GetExtentionKey(name string) (ExtentionKey, error) {
	key := ExtentionKey{}
	strs := strings.SplitN(name, "/", 2)
	if len(strs) == 0 {
		return key, ErrExtensionName
	}
	if len(strs) == 1 {
		key.Author = mid.Meta.String("officialAuthor")
		key.Name = name
	} else {
		key.Author = strs[0]
		key.Name = strs[1]
	}
	return key, nil
}

func (key ExtentionKey) Path(rootdir string) string {
	return filepath.Join(rootdir, key.Author, key.Name)
}

func LoadExtensions(rootdir string, names []string) ([]Extention, error) {
	seen := map[ExtentionKey]bool{}
	exts := make([]Extention, 0, len(names))
	for _, name := range names {
		key, err := GetExtentionKey(name)
		if err != nil {
			return nil, err
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		data, err := ioutil.ReadFile(key.Path(rootdir))
		if err != nil {
			return nil, err
		}
		ext := Extention{}
		if err := json.Unmarshal(data, &ext); err != nil {
			return nil, err
		}
		if ext.Author == "" {
			ext.Author = key.Author
		}
		if ext.Name == "" {
			ext.Name = key.Name
		}
		if ext.EmbeddedAt != nil {
			for lang, x := range ext.EmbeddedAt {
				if x != nil {
					for pos, values := range x {
						if !pos.IsValid() {
							return nil, fmt.Errorf("extension %s lang %s: invalid pos %s", name, lang, pos)
						}
						for i, value := range values {
							if !value.IsValid() {
								return nil, fmt.Errorf("extension %s lang %s: %dth: invalid value %s", name, lang, i+1, value)
							}
						}
					}
				}
			}
		}
		exts = append(exts, ext)
	}
	return exts, nil
}
