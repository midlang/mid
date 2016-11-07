package build

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"strings"

	"github.com/midlang/mid/src/mid/ast"
	"github.com/mkideal/log"
)

func init() {
	gob.Register(&Builder{})
	gob.Register(map[string]*Package{})
}

type ObjectId string

func (id ObjectId) IsValid() bool { return id != "" }

func (id ObjectId) index(pre bool) int {
	i := strings.Index(string(id), ".")
	if i < 0 {
		if pre {
			i = 0
		} else {
			i = -1
		}
	}
	return i
}

func (id ObjectId) Package() string {
	return string(id[:id.index(true)])
}

func (id ObjectId) Name() string {
	return string(id[id.index(false)+1:])
}

func (id ObjectId) String() string { return string(id) }

// Builder
type Builder struct {
	Packages      map[string]*Package
	encodedString string
}

func NewBuilder() *Builder {
	return &Builder{
		Packages: map[string]*Package{},
	}
}

func Build(pkgs map[string]*ast.Package) (*Builder, error) {
	builder := NewBuilder()
	for _, pkg := range pkgs {
		builder.Packages[pkg.Name] = BuildPackage(pkg)
	}
	return builder, nil
}

func (builder *Builder) Encode() string {
	if builder.encodedString != "" {
		return builder.encodedString
	}
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(builder)
	if err != nil {
		log.Fatal("encode builder error: %v", err)
	}
	builder.encodedString = base64.StdEncoding.EncodeToString(buf.Bytes())
	return builder.encodedString
}

func (builder *Builder) Decode(s string) error {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(builder)
}
