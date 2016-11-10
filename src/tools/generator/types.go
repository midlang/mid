package generator

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/midlang/mid/src/mid/build"
)

type Context struct {
	Pkg       *build.Package
	Root      *template.Template
	buildType BuildTypeFunc
}

func (ctx *Context) BuildType(typ build.Type) string {
	return ctx.buildType(typ)
}

func (ctx *Context) Title(s string) string {
	return strings.Title(s)
}

type Package struct {
	*build.Package
	Context *Context
}

func (pkg Package) GenerateDeclsBySubTemplates() (string, error) {
	buf := new(bytes.Buffer)

	if temp := pkg.Context.Root.Lookup("T_const"); temp != nil {
		for _, f := range pkg.Files {
			for _, c := range f.Decls {
				if len(c.Consts) > 0 {
					if err := temp.Execute(buf, NewGenDecl(pkg.Context, f, c)); err != nil {
						return "", err
					}
				}
			}
		}
	}

	for _, f := range pkg.Files {
		for _, b := range f.Beans {
			if temp := pkg.Context.Root.Lookup("T_" + b.Kind); temp != nil {
				if err := temp.Execute(buf, NewBean(pkg.Context, f, b)); err != nil {
					return "", err
				}
			}
		}
	}
	return buf.String(), nil
}

type File struct {
	*build.File
	Context *Context
}

func (f File) GenerateDeclsBySubTemplates() (string, error) {
	buf := new(bytes.Buffer)

	if temp := f.Context.Root.Lookup("T_const"); temp != nil {
		for _, c := range f.Decls {
			if len(c.Consts) > 0 {
				if err := temp.Execute(buf, NewGenDecl(f.Context, f.File, c)); err != nil {
					return "", err
				}
			}
		}
	}

	for _, b := range f.Beans {
		if temp := f.Context.Root.Lookup("T_" + b.Kind); temp != nil {
			if err := temp.Execute(buf, NewBean(f.Context, f.File, b)); err != nil {
				return "", err
			}
		}
	}
	return buf.String(), nil
}

type GenDecl struct {
	*build.GenDecl
	File    *build.File
	Context *Context
}

func NewGenDecl(ctx *Context, file *build.File, c *build.GenDecl) *GenDecl {
	return &GenDecl{
		GenDecl: c,
		File:    file,
		Context: ctx,
	}
}

type Bean struct {
	*build.Bean
	File    *build.File
	Context *Context
}

func NewBean(ctx *Context, file *build.File, b *build.Bean) *Bean {
	return &Bean{
		Bean:    b,
		File:    file,
		Context: ctx,
	}
}
