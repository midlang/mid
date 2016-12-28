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
	// positions
	FileHead       = "file_head"
	BeforeImport   = "before_import"
	InImport       = "in_import"
	AfterImport    = "after_import"
	BeforeConst    = "before_const"
	AfterConst     = "after_const"
	ConstFront     = "const_front"
	ConstBack      = "const_back"
	BeforeEnum     = "before_enum"
	EnumFront      = "enum_front"
	EnumBack       = "enum_back"
	AfterEnum      = "after_enum"
	BeforeStruct   = "before_struct"
	StructFront    = "struct_front"
	StructBack     = "struct_back"
	AfterStruct    = "after_struct"
	BeforeProtocol = "before_protocol"
	ProtocolFront  = "protocol_front"
	ProtocolBack   = "protocol_back"
	AfterProtocol  = "after_protocol"
	BeforeService  = "before_service"
	ServiceFront   = "service_front"
	ServiceBack    = "service_back"
	AfterService   = "after_service"

	// Extension config filename
	ExtConfigFilename = "ext.json"
)

func IsValidKind(kind string) bool {
	switch kind {
	case "package", "file", "const", "enum", "struct", "protocol", "service":
		return true
	default:
		return false
	}
}

type EmbeddedPosition string

func (pos EmbeddedPosition) IsValid() bool {
	switch pos {
	case FileHead,
		BeforeImport,
		InImport,
		AfterImport,
		BeforeConst,
		ConstFront,
		ConstBack,
		AfterConst,
		BeforeEnum,
		EnumFront,
		EnumBack,
		AfterEnum,
		BeforeStruct,
		StructFront,
		StructBack,
		AfterStruct,
		BeforeProtocol,
		ProtocolFront,
		ProtocolBack,
		AfterProtocol,
		BeforeService,
		ServiceFront,
		ServiceBack,
		AfterService:
	default:
		return false
	}
	return true
}

func (pos EmbeddedPosition) Match(at string) bool {
	return string(pos) == at
}

type EmbeddedValue struct {
	Text     string   `join:"text"`
	Template string   `json:"template"`
	Suffix   string   `json:"suffix"`
	Kinds    []string `json:"kinds"`
}

func (v EmbeddedValue) IsValid() bool {
	valid := (v.Text != "" && v.Template == "") || (v.Text == "" && v.Template != "")
	if !valid {
		return false
	}
	for _, kind := range v.Kinds {
		if !IsValidKind(kind) {
			return false
		}
	}
	return true
}

func (v EmbeddedValue) MatchKind(kind string) bool {
	if len(v.Kinds) == 0 {
		return true
	}
	for _, k := range v.Kinds {
		if k == kind {
			return true
		}
	}
	return false
}

type Extension struct {
	Name    string   `json:"name"`
	Author  string   `json:"author"`
	URL     string   `json:"url"`
	Version string   `json:"version"`
	Deps    []string `json:"deps"`
	// language -> position -> embedded_values
	EmbeddedAt map[string]map[EmbeddedPosition][]EmbeddedValue `json:"at"`

	Path string `json:"path"`
}

func (e Extension) Find(lang, kind, at string) []EmbeddedValue {
	if e.EmbeddedAt == nil {
		return nil
	}
	m, ok := e.EmbeddedAt[lang]
	if !ok || m == nil {
		return nil
	}
	values := make([]EmbeddedValue, 0)
	for pos, vals := range m {
		if pos.Match(at) {
			for _, val := range vals {
				if val.MatchKind(kind) {
					values = append(values, val)
				}
			}
		}
	}
	return values
}

type ExtensionKey struct {
	Author string
	Name   string
}

func GetExtensionKey(name string) (ExtensionKey, error) {
	key := ExtensionKey{}
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

func (key ExtensionKey) Path(rootdir string) string {
	return filepath.Join(rootdir, key.Subdir())
}

func (key ExtensionKey) Subdir() string {
	return filepath.Join(key.Author, key.Name)
}

func LoadExtensions(rootdir string, names []string) ([]Extension, error) {
	seen := map[ExtensionKey]bool{}
	exts := make([]Extension, 0, len(names))
	for _, name := range names {
		key, err := GetExtensionKey(name)
		if err != nil {
			return nil, err
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		data, err := ioutil.ReadFile(filepath.Join(key.Path(rootdir), ExtConfigFilename))
		if err != nil {
			return nil, err
		}
		ext := Extension{}
		if err := json.Unmarshal(data, &ext); err != nil {
			return nil, err
		}
		if ext.Author == "" {
			ext.Author = key.Author
		}
		if ext.Name == "" {
			ext.Name = key.Name
		}
		ext.Path = key.Subdir()
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
