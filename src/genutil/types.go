package genutil

import (
	"bytes"

	"github.com/midlang/mid/src/mid/build"
)

// Package wraps build.Package
type Package struct {
	*build.Package
}

func (pkg Package) GenerateDeclsBySubTemplates() (string, error) {
	buf := new(bytes.Buffer)

	if temp := context.Root.Lookup("T_const"); temp != nil {
		constDecls := make([]*GenDecl, 0)
		for _, f := range pkg.Files {
			for _, c := range f.Decls {
				if len(c.Consts) > 0 {
					constDecls = append(constDecls, NewGenDecl(f, c))
				}
			}
		}
		if len(constDecls) > 0 {
			if err := temp.Execute(buf, constDecls); err != nil {
				return "", err
			}
		}
	}

	for _, f := range pkg.Files {
		for _, b := range f.Beans {
			if temp := context.Root.Lookup("T_" + b.Kind); temp != nil {
				if err := temp.Execute(buf, NewBean(f, b)); err != nil {
					return "", err
				}
			}
		}
	}
	return buf.String(), nil
}

// File wraps build.File
type File struct {
	*build.File
}

func (f File) GenerateDeclsBySubTemplates() (string, error) {
	buf := new(bytes.Buffer)

	if temp := context.Root.Lookup("T_const"); temp != nil {
		for _, c := range f.Decls {
			if len(c.Consts) > 0 {
				if err := temp.Execute(buf, NewGenDecl(f.File, c)); err != nil {
					return "", err
				}
			}
		}
	}

	for _, b := range f.Beans {
		if temp := context.Root.Lookup("T_" + b.Kind); temp != nil {
			if err := temp.Execute(buf, NewBean(f.File, b)); err != nil {
				return "", err
			}
		}
	}
	return buf.String(), nil
}

// GenDecl wraps build.GenDecl
type GenDecl struct {
	*build.GenDecl
	File *build.File
}

func NewGenDecl(file *build.File, c *build.GenDecl) *GenDecl {
	return &GenDecl{
		GenDecl: c,
		File:    file,
	}
}

// Bean wraps build.Bean
type Bean struct {
	*build.Bean
	File *build.File
}

func NewBean(file *build.File, b *build.Bean) *Bean {
	return &Bean{
		Bean: b,
		File: file,
	}
}

// Extends gets extends of bean as a string slice
func (bean *Bean) Extends(ctx *Context) []string {
	extends := bean.Bean.Extends()
	if len(extends) == 0 {
		return nil
	}
	var strs []string
	for _, ex := range extends {
		strs = append(strs, ctx.BuildType(ex))
	}
	return strs
}

// AddTag is a chain function for adding tag
func (bean *Bean) AddTag(key, value string, field *build.Field) *build.Field {
	return bean.addTag(key, value, field, true)
}

func (bean *Bean) AddTagNX(key, value string, field *build.Field) *build.Field {
	return bean.addTag(key, value, field, false)
}

func (bean *Bean) addTag(key, value string, field *build.Field, force bool) *build.Field {
	_, found := field.Tag.Lookup(key)
	if !found || force {
		field.Tag.Set(key, value)
	}
	return field
}

type Field struct {
	*build.Field
	Bean *Bean
	Type string
}

func NewField(bean *Bean, field *build.Field, typ string) *Field {
	return &Field{
		Field: field,
		Bean:  bean,
		Type:  typ,
	}
}
