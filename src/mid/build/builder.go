package build

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"sort"
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
	Packages       map[string]*Package
	SortedPackages []*Package
	encodedString  string
}

func NewBuilder() *Builder {
	return &Builder{
		Packages: map[string]*Package{},
	}
}

func Build(pkgs map[string]*ast.Package) (*Builder, error) {
	builder := NewBuilder()
	log.Debug("pkgs: %v", pkgs)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, imp := range file.Imports {
				_, pkgId := imp.Package.IsString()
				importedPkg, ok := pkgs[pkgId]
				if !ok {
					return nil, fmt.Errorf("package `%s` not found", pkgId)
				}
				obj := ast.NewObj(ast.Pkg, importedPkg.Name)
				obj.Decl = importedPkg
				pkg.Imports[pkgId] = obj
			}
		}
	}
	for _, pkg := range pkgs {
		builtPkg := BuildPackage(pkg)
		builder.Packages[pkg.Name] = builtPkg
		builder.SortedPackages = append(builder.SortedPackages, builtPkg)
	}
	sort.Slice(builder.SortedPackages, func(i, j int) bool {
		return builder.SortedPackages[i].Name < builder.SortedPackages[j].Name
	})
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
